function Folder(path) {
    this.path = path;
    this.name = function() {
        var n = path.lastIndexOf('/');
        name = path;
        if (n != -1 && path != "/") {
            name = name.substring(n + 1);
        }
        return name;
    };
    console.log(this.path);


    //would want an uid of some sorts too
    this.Contents = function() {
        var request = new XMLHttpRequest();
        request.open('GET', 'api/open/' + this.path, false);  // `false` makes the request synchronous
        request.send(null);

        if (request.status === 200) {
            console.log(request.responseText);
            f = JSON.parse(request.responseText);
            if (f.error != null) {
                return []; //emptry
            } else {
                arr = [];
                for (var i = 0; i < f.contents.length; i++) {
                    arr[i] = f.contents[i];
                }
                return arr;
            }
        }
    }
}

function getClickedFolder(column, name) {
    return getPath(column.folder.path, name);
}

function getPath(parentPath, childName) {
    p = parentPath + "/" + childName;              
    p = p.replace("//", "/");
    return p;
}

function Column(folder) {
    this.folder = folder
    this.active = null
    this.HTMLElement = null
}

var columns = [];


function makeHTMLColumn(column, files) {
    var col = document.createElement("div");
    column.HTMLElement = col;
    var title = document.createElement("b");
    title.appendChild(document.createTextNode(column.folder.name()));
    col.appendChild(title);
    col.id = column.folder.name();
    col.className = "miller-column";
    if (files == undefined) {
        console.log("no folder given");
        files = column.folder.Contents();
    } 
    if (files.length == 0) {
        var item = document.createElement("p");
        item.appendChild(document.createTextNode("Empty"));
        col.appendChild(item);
    } else {
        var list = col.appendChild(document.createElement("ul"));
        for (var i = 0; i < files.length; i++) {
            var folder = files[i];
            if (files[i][0] != ".") {
                var item = document.createElement("li");
                item.appendChild(document.createTextNode(folder));
                item.onclick = (function() {
                    var folder = files[i]; 
                    return function() {clickColumn(column, folder);}
                })();
                list.appendChild(item);
            }
        }
    }
    return col
}
function appendColumn(colNode) {
    document.getElementById("miller-container").appendChild(makeHTMLColumn(colNode));
}

function markActive(column) {
    var list = column.HTMLElement.childNodes[1].childNodes;
    for (var i = 0; i < list.length; i++) {
        if (list[i].innerText != column.active) {
            list[i].classList.remove("item-active");
        } else {
            list[i].classList.add("item-active");
        }
    }
}

function clickColumn(column, name) {
    column.active = name;
    markActive(column);
    closeColumnsAfter(column);
    columns[index + 1] = new Column(getFolder(getClickedFolder(column, name))); //so this is the javascript power...
    appendColumn(columns[index + 1]);
}

function closeColumnsAfter(column) {
    index = columns.indexOf(column);
    for (var i = index + 1; i < columns.length; i++) {
        document.getElementById("miller-container").removeChild(columns[i].HTMLElement);
        var request = new XMLHttpRequest();
        request.open('GET', 'api/close/' + columns[i].folder.path);  // `false` makes the request synchronous
        request.send(null);
    }
    columns.length = index + 1;
}

function getFolder(path) {
    console.log(path)
    var request = new XMLHttpRequest();
    request.open('GET', 'api/open/' + path, false);  // `false` makes the request synchronous
    request.send(null);

    if (request.status === 200) {
        console.log(request.responseText);
        f = JSON.parse(request.responseText)
        r = new Folder(path)
        return r
    }
}



columns[0] = new Column(getFolder("/"));
appendColumn(columns[0]);

function refreshFolder(data) {
    console.log("update column: " + data);
    f = JSON.parse(data);
    if (f.error != null) {
        return; //emptry
    } 
    
    var column = null;
    for (var i = 0; i < columns.length; i++) {
        if (columns[i].folder.path == f.path) { 
            column = columns[i];
            break;
        }
    }
    if (columns == null) {
        console.log("don't need this update");
        return;
    }

    if (f.newPath != "") {
        column.folder.path = f.newPath;
        closeColumnsAfter(column);
    }
    
    oldHTMLElement = column.HTMLElement;
    newHTMLElement = makeHTMLColumn(column, f.contents);
    document.getElementById("miller-container").replaceChild(newHTMLElement, oldHTMLElement);
}

var host = "http://localhost:8080";

var opts = {
    // The base URL is appended to the host string. This value has to match with the server value.
    baseURL: "/channel/",

    // Force a socket type.
    // Values: false, "WebSocket", "AjaxSocket"
    forceSocketType: false,

    // Kill the connect attempt after the timeout.
    connectTimeout:  10000,

    // If the connection is idle, ping the server to check if the connection is stil alive.
    pingInterval:           35000,
    // Reconnect if the server did not response with a pong within the timeout.
    pingReconnectTimeout:   5000,

    // Whenever to automatically reconnect if the connection was lost.
    reconnect:          true,
    reconnectDelay:     1000,
    reconnectDelayMax:  5000,
    // To disable set to 0 (endless).
    reconnectAttempts:  10,

    // Reset the send buffer after the timeout.
    resetSendBufferTimeout: 10000
};

// Create and connect to the server.
// Optional pass a host string and options.
var socket = glue(host, opts);

socket.onMessage(function(data) {
    console.log("onMessage: " + data);

    console.log(data.substr(0, 7))
    if (data.substr(0, 8) == "update: ") {
        refreshFolder(data.substring(8))
    }

    if (data == "Folder update") {
        // https://zeit.co/blog/async-and-await
        function sleep (time) {
          return new Promise((resolve) => setTimeout(resolve, time));
        }

        // Usage!
        sleep(500).then(() => {
            location.reload()
            // Do something after the sleep!
        })
    }

    // Echo the message back to the server.
    //socket.send("echo: " + data);
});


socket.on("connected", function() {
    console.log("connected");
});

socket.on("connecting", function() {
    console.log("connecting");
});

socket.on("disconnected", function() {
    console.log("disconnected");
});

socket.on("reconnecting", function() {
    console.log("reconnecting");
});

socket.on("error", function(e, msg) {
    console.log("error: " + msg);
});

socket.on("connect_timeout", function() {
    console.log("connect_timeout");
});

socket.on("timeout", function() {
    console.log("timeout");
});

socket.on("discard_send_buffer", function() {
    console.log("some data could not be send and was discarded.");
});
