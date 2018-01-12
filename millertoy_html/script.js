function Folder(name) {
    this.name = name;
    this.content = []; 
    //would want an uid
}

//fake root to test
var root = new Folder("root");
for (var i = 0; i < 20; i++) {
    root.content[i] = new Folder("child " + i);
    for (var j = 0; j < 10; j++) {
        root.content[i].content[j] = new Folder("subchild " + j);
        for (var k = 0; k < 5; k++) {
            root.content[i].content[j].content[k] = new Folder("subsubchild " + k);
            root.content[i].content[j].content[k].content[0] = new Folder("subsubsubchild 1");
        }
    }
}

function Column(folder) {
    this.folder = folder
    this.selected = null
    this.HTMLElement = null
}

var columns = [];
columns[0] = new Column(root);


function makeHTMLColumn(column) {
    var col = document.createElement("div");
    column.HTMLElement = col;
    var title = document.createElement("b");
    title.appendChild(document.createTextNode(column.folder.name));
    col.appendChild(title);
    col.id = column.folder.name;
    col.className = "miller-column";
    if (column.folder.content.length == 0) {
        var item = document.createElement("p");
        item.appendChild(document.createTextNode("Empty"));
        col.appendChild(item);
    } else {
        var list = col.appendChild(document.createElement("ul"));
        for (var i = 0; i < column.folder.content.length; i++) {
            var folder = column.folder.content[i];
            var item = document.createElement("li");
            item.appendChild(document.createTextNode(folder.name));
            item.onclick = (function() {
                var folder = column.folder.content[i]; 
                return function() {clickColumn(column, folder);}
            })();
            list.appendChild(item);
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


appendColumn(columns[0]);

