// Command parisgo-gen-tree generates a hierarchy of directories and files
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

var (
	depth  = flag.Int("depth", 5, "depth of directories to generate")
	ndirs  = flag.Int("ndirs", 5, "number of directories per directory")
	nfiles = flag.Int("nfiles", 10, "maximal number of files to generate per directory")
	root   = flag.String("root", "", "path to top-level directory under which to generate the structure ($PWD by default)")
)

func main() {
	var err error

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`%[1]s generates a hierarchy of directories and files.

Usage of %[1]s:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *root == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
		*root = wd
	}

	*root, err = filepath.Abs(*root)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "generating %d-depth dirs with %d-max files each...\n",
		*ndirs, *nfiles,
	)

	state := State{
		depth:    0,
		depthmax: *depth,
		ndirs:    *ndirs,
		nfiles:   *nfiles,
	}

	for idir := 0; idir < *ndirs; idir++ {
		dirname := filepath.Join(*root, fmt.Sprintf("dir-%04d", idir))
		err := gen(state, dirname)
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
	}

	fmt.Fprintf(os.Stderr, "generating %d-depth dirs with %d-max files each... [done]\n",
		*ndirs, *nfiles,
	)

}

type State struct {
	depth    int
	depthmax int
	ndirs    int
	nfiles   int
}

func gen(state State, root string) error {
	// fmt.Printf(">>> root=%q...\n", root)
	ndirs := rand.Intn(state.ndirs)
	nfiles := rand.Intn(state.nfiles)

	state.depth++

	err := os.MkdirAll(root, 0755)
	if err != nil {
		return err
	}

	for i := 0; i < ndirs; i++ {
		dirname := filepath.Join(root, fmt.Sprintf("dir-%04d", i))
		err := os.MkdirAll(dirname, 0755)
		if err != nil {
			return err
		}

		if state.depth < state.depthmax {
			err := gen(state, dirname)
			if err != nil {
				return err
			}
		}
	}

	for i := 0; i < nfiles; i++ {
		fname := filepath.Join(root, fmt.Sprintf("file-%04d", i))
		f, err := os.Create(fname)
		if err != nil {
			return err
		}
		defer f.Close()

		err = genfile(f)
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return fmt.Errorf("error closing %q: %v", f.Name(), err)
		}
	}

	return nil
}

func genfile(f *os.File) error {
	size := rand.Intn(8*1024) + 1
	data := make([]byte, size)
	for i := range data {
		// generate some ascii
		data[i] = byte(rand.Intn(126-33) + 33)
	}

	_, err := f.Write(data)
	return err
}
