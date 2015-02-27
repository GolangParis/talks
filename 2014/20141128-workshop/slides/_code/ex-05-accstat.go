// STARTIMPORT OMIT

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	root = flag.String("root", "", "top-level directory to inspect (default=$PWD)")
	out  = flag.String("o", "index.json", "path to index file")
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

	// STARTWALK OMIT
	stats := make([]Stat, 0, 10) // HLxxx
	err = filepath.Walk(*root, func(path string, fi os.FileInfo, err error) error {
		fmt.Printf(">>> %s\n", path)
		stat := Stat{
			Name: path,
			Size: fi.Size(),
		}
		switch {
		case fi.IsDir():
			stat.Type = Dir
		case fi.Mode().IsRegular():
			stat.Type = File
		}

		stats = append(stats, stat) // HLxxx

		if fi.IsDir() {
			return nil
		}
		return err
	})
	// ENDWALK OMIT

	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// STARTENCODE OMIT
	f, err := os.Create(*out)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	defer f.Close() // HLxxx

	err = json.NewEncoder(f).Encode(&stats) // HLxxx
	if err != nil {
		log.Fatalf("error encoding stats: %v\n", err)
	}

	err = f.Close() // HLxxx
	if err != nil {
		log.Fatalf("error closing file: %v\n", err)
	}
	// ENDENCODE OMIT
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
