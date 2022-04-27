// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is main program driver for the loopback filesystem from
// github.com/hanwen/go-fuse/fs/, a filesystem that shunts operations
// to an underlying file system.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
)

const (
	MountPoint = "/share/Public/mnt/hykuan-test"
	SourceDir = "/share/CACHEDEV1_DATA/.qpkg/hbs3-apigateway/tmp/hykuan"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	// Scans the arg list and sets up flags
	debug := flag.Bool("debug", false, "print debugging messages.")
	//other := flag.Bool("allow-other", false, "mount with -o allowother.")

	flag.Parse()
	//if flag.NArg() < 2 {
	//	fmt.Printf("usage: %s MOUNTPOINT ORIGINAL\n", path.Base(os.Args[0]))
	//	fmt.Printf("\noptions:\n")
	//	flag.PrintDefaults()
	//	os.Exit(2)
	//}

	loopbackRoot, err := fs.NewLoopbackRoot(SourceDir)
	if err != nil {
		log.Fatalf("NewLoopbackRoot(%s): %v\n", SourceDir, err)
	}

	sec := time.Second
	opts := &fs.Options{
		// These options are to be compatible with libfuse defaults,
		// making benchmarking easier.
		AttrTimeout:  &sec,
		EntryTimeout: &sec,
	}
	opts.Debug = *debug
	//opts.AllowOther = *other
	//if opts.AllowOther {
	//	// Make the kernel check file permissions for us
	//	opts.MountOptions.Options = append(opts.MountOptions.Options, "default_permissions")
	//}

	// First column in "df -T": original dir
	opts.MountOptions.Options = append(opts.MountOptions.Options, "fsname="+SourceDir)
	// Second column in "df -T" will be shown as "fuse." + Name
	opts.MountOptions.Name = "loopback"
	// Leave file permissions on "000" files as-is
	opts.NullPermissions = true

	server, err := Mount(MountPoint, loopbackRoot, opts)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}

	fmt.Println("Mounted!")
	server.Wait()
}