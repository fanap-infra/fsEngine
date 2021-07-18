# File Storage Engine

`File Storage Engine` (FSEngine) is a write-optimized object storage library designed mainly for storing large, continuous streams on commodity hardware and Disk-Drives. This project includes APIs and primitives for storing raw data and fetching data from the storage. Use-cases of this library include storing live video streams, logs, metadata, or any write-optimized workload.

This library might be used with `ArchiverMedia` wrapper library for easily transforming media objects into consumable binary packets for `FSEngine`.

This library uses a pseudo log storage model for storing multiple data streams into a single file. It also supports a circular storage model for a fixed-sized store.

Currently, FSEngine has been offered as a [Golang](https://www.golang.com) only library, and all APIs are usable within any Golang 1.12+ project.

## Features
- [x] Write-optimized storage
- [x] Local storage
- [x] Circular storage
- [ ] Multi-instance storage
- [ ] Spare storage with ECC

## Installing
1. To install FSEngine library into and existing project, having Go installed, use the bellow Go command to install FSEngine.
```shell
$ go get github.com/fanap-infra/fsEngine
```
2. Import it to your library
```go
import "github.com/fanap-infra/fsEngine"
```
## How to use
First you need to establish storage location on local system and initialize storage engine by maximum size of storage used by the engine:

```go
package main

import (
	"github.com/fanap-infra/fsEngine"
	"github.com/fanap-infra/log"
)

func main() {
	// define storage location and size
	const storagePath = "/var/fsEngine/volume1"
	const storageSize = 1 << 32 // 4GB volume

	fileSystem, err := fsEngine.CreateFileSystem(storagePath, storageSize,
		fsEngine.BLOCKSIZE, log.GetScope("Example"))
	if err != nil {
		log.Fatal(err)
	}
	
	return
}


```

This is the barebone definition for using FSEngine, however, to store objects you have to create a virtual identity for the object which here is called a `virtualFile`.
```go
	fileSystem, err := fsEngine.CreateFileSystem(storagePath, storageSize,
		fsEngine.BLOCKSIZE, &EventsTrigger{}, log.GetScope("Example"))
	if err != nil {
		log.Fatal(err)
	}
    const vfID = 1
    const vfName = "Hello"
    virtualFile, err := fileSystem.NewVirtualFile(vfID, vfName)
    if err != nil {
        log.Fatal(err)
    }
    _, err = virtualFile.Write([]byte("HelloStorage"))
    if err != nil {
        log.Fatal(err)
    }
    err = virtualFile.Close()
	if err != nil {
        log.Fatal(err)
    }
```
## Reading Object

To read object from storage you first need to open the archiver which is internally a `Singleton` object and once initialized, returns the main object. For retrieving an object, provide `object-ID` and use the code bellow.
```go
package main

import (
	"github.com/fanap-infra/fsEngine"
	"github.com/fanap-infra/log"
)

func main() {
	// define storage location and size
	const storagePath = "/var/fsEngine/volume1"
	const fileID = 68 // Some random fileID

	fileSystem, err := fsEngine.ParseFileSystem(storagePath, log.GetScope("Example"))
	if err != nil {
		log.Fatal(err)
	}
	
	virtualFile, err := fileSystem.OpenVirtualFile(fileID)
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, virtualFile.Size())
	virtualFile.Read(data, fileID)
	return
}
```
PS: 

[comment]: <> (TODO: Complete readFile section)
[comment]: <> (## Reading an object from storage)

[comment]: <> (To read an object from storage, Open )