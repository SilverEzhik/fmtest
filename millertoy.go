package main

import (
	"./fm"
	"./fs"
	"encoding/json"
	"fmt"
	"github.com/desertbit/glue"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func getDir(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get dir")
	vars := mux.Vars(r)
	pathId := vars["path"]
	var path string

	if len(pathId) == 0 {
		path = "/"
	} else if string(pathId[0]) == "~" {
		path = pathId
	} else {
		path = "/" + pathId
	}

	fmt.Println("path:\"" + path + "\"")

	a := JSONFolder{Path: path, Contents: make([]string, 0)}

	var f fm.Folder
	f, err := fs.GetFolder(path)
	if err != nil {
		fmt.Println("error:", err)
		a.Error = err
	} else {
		a.Contents = fm.GetSortedFileList(f)
	}

	jsonFiles, err := json.Marshal(a)
	fmt.Fprintln(w, string(jsonFiles))
	fmt.Println(string(jsonFiles))
	fmt.Println(path, len(a.Contents))
}

type JSONFolder struct {
	Path     string   `json:"path"`
	Contents []string `json:"contents"`
	Error    error    `json:"error"`
}

func watchFolder(path string) {
	fmt.Println(path)

	var f fm.Folder
	f, err := fs.GetFolder(path)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	update := f.Watch()

Loop:
	for {
		select {
		case <-update:
			if f.Path() == "" && f.Contents() == nil {
				fmt.Println("folder object gone")
				break Loop
			}

			for _, socket := range server.Sockets() {
				socket.Write("Folder update")
			}

			//weird mechanism for testing this
			for file := range f.Contents() {
				if file == "close" {
					f.Close()
					break Loop
				}
			}
		}
	}
}

var server *glue.Server

// main function to boot up everything
func main() {
	router := mux.NewRouter()

	//read files
	router.HandleFunc(`/api/open/{path:.*}`, getDir).Methods("GET")

	//notify api
	// Create a new glue server.
	server = glue.NewServer(glue.Options{
		HTTPListenAddress: ":8080",
		HTTPSocketType:    glue.HTTPSocketTypeNone,
		HTTPHandleURL:     "/channel/",
	})

	// Release the glue server on defer.
	// This will block new incoming connections
	// and close all current active sockets.
	defer server.Release()

	// Set the glue event function to handle new incoming socket connections.
	server.OnNewSocket(onNewSocket)

	// Run the glue server.
	go server.Run()

	router.PathPrefix("/channel/").Handler(server)

	//serve static
	staticFileDirectory := http.Dir("./millertoy_html/")
	staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDirectory))
	router.PathPrefix("/").Handler(staticFileHandler).Methods("GET")
	//check static updates
	go watchFolder("./millertoy_html")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func onNewSocket(s *glue.Socket) {
	// Set a function which is triggered as soon as the socket is closed.
	s.OnClose(func() {
		log.Printf("socket closed with remote address: %s", s.RemoteAddr())
	})

	// Set a function which is triggered during each received message.
	s.OnRead(func(data string) {
		// Echo the received data back to the client.
		s.Write(data)
	})

	// Send a welcome string to the client.
	s.Write("Hello Client")
	log.Println("Channel pinged")
}
