// STARTIMPORT OMIT

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	root = flag.String("root", "", "top-level directory to inspect (default=$PWD)")
)

// ENDIMPORT OMIT

func main() {
	var err error
	flag.Parse()

	if *root == "" {
		*root = "."
	}

	*root, err = filepath.Abs(os.ExpandEnv(*root))
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	err = filepath.Walk(*root, func(path string, fi os.FileInfo, err error) error {
		fmt.Printf(">>> %s\n", path)
		if fi.IsDir() {
			return nil
		}
		return err
	})

	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

// STARTSTAT OMIT

// Stat holds metadata about dirs and files
type Stat struct {
	Name string // full path to object
	Type Type   // type of object (Dir|File)
	Size int64  // size of object in bytes
}

// Type describes the type of a filesystem object (dir|file)
type Type int

// possible types of filesystem objects
const (
	Invalid Type = 0
	Dir     Type = 1
	File    Type = 2
)

// ENDSTAT OMIT
