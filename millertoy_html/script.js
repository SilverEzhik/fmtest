function Folder(path) {
    this.path = path;
    var n = path.lastIndexOf('/');
    name = path;
    if (n > 0) {
        name = name.substring(n + 1);
    }
    this.name = name;
    console.log(this.path)


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
                    arr[i] = new Folder(this.path + "/" + f.contents[i])               
                }
                return arr;
            }
        }
    }
}

/*
//fake root to test
var root = new Folder("root");

for (var i = 0; i < 20; i++) {
    root.contents[i] = new Folder("child " + i);
    for (var j = 0; j < 10; j++) {
        root.contents[i].contents[j] = new Folder("subchild " + j);
        for (var k = 0; k < 5; k++) {
            root.contents[i].contents[j].contents[k] = new Folder("subsubchild " + k);
            root.contents[i].contents[j].contents[k].contents[0] = new Folder("subsubsubchild 1");
        }
    }
}
*/

function Column(folder) {
    this.folder = folder
    this.selected = null
    this.HTMLElement = null
}

var columns = [];


function makeHTMLColumn(column) {
    var col = document.createElement("div");
    column.HTMLElement = col;
    var title = document.createElement("b");
    title.appendChild(document.createTextNode(column.folder.name));
    col.appendChild(title);
    col.id = column.folder.name;
    col.className = "miller-column";
    files = column.folder.Contents();
    if (files.length == 0) {
        var item = document.createElement("p");
        item.appendChild(document.createTextNode("Empty"));
        col.appendChild(item);
    } else {
        var list = col.appendChild(document.createElement("ul"));
        for (var i = 0; i < files.length; i++) {
            var folder = files[i];
            if (files[i].name[0] != ".") {
                var item = document.createElement("li");
                item.appendChild(document.createTextNode(folder.name));
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

function clickColumn(column, folder) {
    var list = column.HTMLElement.childNodes[1].childNodes;
    for (var i = 0; i < list.length; i++) {
        if (list[i].innerText != folder.name) {
            list[i].classList.remove("item-active");
        } else {
            list[i].classList.add("item-active");
        }
    }
    index = columns.indexOf(column);
    for (var i = index + 1; i < columns.length; i++) {
        document.getElementById("miller-container").removeChild(columns[i].HTMLElement);
    }
        
    columns.length = index + 1;
    columns[index + 1] = new Column(folder);
    appendColumn(columns[index + 1]);
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
