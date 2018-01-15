# fmtest

I like Miller Columns, but most Linux file managers do not (and the ones that do are not very good). So here I am writing my own.

## Structure

I want to have three layers:

`File System`  
`File Manager`  
`User Interface`  

FM layer defines the interfaces for FS implementations to provide, and will provide various services, such as handling paths, sorting by name, and etc.

FS layer will deal with file system interaction - and will contain actual implementations of interface functions defined by the FM layer - watching folders, getting their contents, copying files, etc. While I am focusing on the local FS layer, it is, in theory, possible to use anything as the file system, from a cloud service to an SQL database - implement the interface and go.

UI layer will be what the user actually sees. This one will be platform-dependent - I don't want to make an Electron app. First priority will be a Linux app in GTK or Qt.

## Current contents

Right now, I have implemented a piece of the FS model for watching folders - channels for getting updates and etc. There are two test Go programs in there, `fmtest`, which accepts folders as arguments and prints their contents when a change is detected. The second is `millertoy`, which lets you browse your file system in the browser - it's absolutely minimal and will not be used as a base for the UI, as I want it to be native.
