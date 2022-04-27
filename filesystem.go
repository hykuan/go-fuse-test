package main

import (
	"fmt"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"os"
	"path/filepath"
	"time"
)

func Mount(dir string, root fs.InodeEmbedder, options *fs.Options) (*fuse.Server, error) {
	if options == nil {
		oneSecond := time.Second
		options = &fs.Options{
			EntryTimeout: &oneSecond,
			AttrTimeout:  &oneSecond,
		}
	}

	rawFS := NewRawFileSystem(dir, root, options)
	server, err := fuse.NewServer(rawFS, dir, &options.MountOptions)
	if err != nil {
		return nil, err
	}

	go server.Serve()
	if err := server.WaitMount(); err != nil {
		// we don't shutdown the serve loop. If the mount does
		// not succeed, the loop won't work and exit.
		return nil, err
	}

	return server, nil
}

type RawFileSystem fuse.RawFileSystem

type MyRawFileSystem struct {
	RawFileSystem
	Dir string
	JobChan chan int
}

func (sys *MyRawFileSystem) OpenDir(cancel <-chan struct{}, input *fuse.OpenIn, out *fuse.OpenOut) (status fuse.Status) {
	fmt.Println("OpenDir!!!!!!!!!")
	if len(sys.JobChan) == cap(sys.JobChan) {
		return sys.RawFileSystem.OpenDir(cancel, input, out)
	}

	sys.JobChan <- 1
	now := time.Now().Format(time.RFC3339)
	path := filepath.Join(sys.Dir, now)
	os.Create(path)

	time.Sleep(3*time.Second)
	<-sys.JobChan


	return sys.RawFileSystem.OpenDir(cancel, input, out)
}

func (sys *MyRawFileSystem) ReleaseDir(input *fuse.ReleaseIn) {
	fmt.Println("ReleaseDir!!!!!!!!!")
	sys.RawFileSystem.ReleaseDir(input)


}

//func (sys *MyRawFileSystem) StatFs(cancel <-chan struct{}, input *fuse.InHeader, out *fuse.StatfsOut) (code fuse.Status) {
//	fmt.Println("StatFs!!!!!!!!!")
//	now := time.Now().Format(time.RFC3339)
//	path := filepath.Join(sys.Dir, now)
//	os.Create(path)
//
//	time.Sleep(3*time.Second)
//
//
//	return sys.RawFileSystem.StatFs(cancel, input, out)
//}


func NewRawFileSystem(dir string, root fs.InodeEmbedder, opts *fs.Options) fuse.RawFileSystem {

	return &MyRawFileSystem{fs.NewNodeFS(root, opts), dir ,make(chan int, 1)}
}
