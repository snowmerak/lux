---
description: fsnotify is a Go library that provides cross-platform filesystem notifications.
tags: [go, library, filesystem, notifications, events, fsevents, inotify, kqueue]
---

# fsnotify

`fsnotify` provides cross-platform filesystem notifications for Go. It supports Windows, Linux, macOS, BSD, and illumos.

## Installation

```bash
go get github.com/fsnotify/fsnotify
```

## Basic Usage

To use `fsnotify`, create a new watcher and add the paths you want to monitor. You should listen for events and errors in a separate goroutine.

```go
package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add("/tmp")
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
```

## Event Types

The `Op` type describes a set of operations.

- `Create`: A new file or directory was created.
- `Write`: A file was written to.
- `Remove`: A file or directory was removed.
- `Rename`: A file or directory was renamed.
- `Chmod`: Attributes were changed.

## Key Considerations

- **Recursive Watching**: `fsnotify` does **not** support recursive watching out of the box. You must manually add every subdirectory you want to watch.
- **Watching Files vs Directories**: It is generally recommended to watch the **parent directory** rather than individual files. Many editors (like Vim or VS Code) save files by creating a temporary file and renaming it, which can break a direct watch on a file.
- **Buffers**: If events are generated faster than they are consumed, the buffer might overflow.

## Platform Support

- **Linux**: uses `inotify`.
- **macOS/BSD**: uses `kqueue`.
- **Windows**: uses `ReadDirectoryChangesW`.
- **illumos**: uses `FEN`.

## FAQ

- **Do I need a goroutine?**: Yes, you must consume the `Events` and `Errors` channels, typically in a `select` block within a goroutine.
- **Network Filesystems**: `fsnotify` generally does not work with NFS, SMB, or FUSE as they don't provide native OS notifications to the client.
