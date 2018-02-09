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

var folderWatchers map[string]chan struct{} = make(map[string]chan struct{})

//All messages will be broadcast to everyone - so this is not suitable for multiple users.
func broadcast(message string) {
	for _, socket := range server.Sockets() {
		socket.Write(message)
	}
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

	for {
		select {
		case <-folderWatchers[path]:
			f.Close()
			delete(folderWatchers, path)
			return
		case <-update:

			if f.Path() == "" && f.Contents() == nil {
				fmt.Println("folder object gone")
				message := JSONUpdate{Path: path, Deleted: true}
				jsonMsg, _ := json.Marshal(message)
				broadcast(string(jsonMsg))
				delete(folderWatchers, path)
				break
			} else {
				message := JSONUpdate{Path: path, Deleted: false, Contents: fm.GetSortedFileList(f)}
				if path != f.Path() {
					message.NewPath = f.Path()
					folderWatchers[f.Path()] = folderWatchers[path]
					delete(folderWatchers, path)
					path = f.Path()
				}
				jsonMsg, _ := json.Marshal(message)
				broadcast("update: " + string(jsonMsg))
			}
		}
	}
}

type JSONUpdate struct {
	Path     string   `json:"path"`
	NewPath  string   `json:"newPath"`
	Deleted  bool     `json:"deleted"`
	Contents []string `json:"contents"`
}

func closeDir(w http.ResponseWriter, r *http.Request) {
	fmt.Println("close dir")
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

	close(folderWatchers[path])
	fmt.Fprintln(w, "ok")
}

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
	folderWatchers[path] = make(chan struct{})
	go watchFolder(path)

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

func refreshFolder(path string) {
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

	//close folders
	router.HandleFunc(`/api/close/{path:.*}`, closeDir).Methods("GET")

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
	go refreshFolder("./millertoy_html")

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
