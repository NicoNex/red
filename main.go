package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

var (
	prnt     bool
	repl     []byte
	glob     string
	verbose  bool
	maxdepth int
	allFiles bool
	wg       sync.WaitGroup
	re       *regexp.Regexp
)

func readDir(filename string) ([]os.FileInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return []os.FileInfo{}, err
	}
	defer file.Close()
	return file.Readdir(-1)
}

func matchGlob(fname string) bool {
	ok, err := filepath.Match(glob, filepath.Base(fname))
	if err != nil {
		fmt.Println(err)
	}
	return ok
}

func edit(fpath string) {
	defer wg.Done()

	b, err := os.ReadFile(fpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	tmp := re.ReplaceAll(b, repl)
	if prnt {
		fmt.Print(string(tmp))
	} else if re.Match(b) {
		if verbose {
			fmt.Printf("Writing %s\n", fpath)
		}
		if err := os.WriteFile(fpath, tmp, 0644); err != nil {
			fmt.Println(err)
		}
	}
}

func editStdin() {
	var reader = bufio.NewReader(os.Stdin)

	b, err := reader.ReadBytes(0)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}
	fmt.Print(string(re.ReplaceAll(b, repl)))
}

// Recursively walks in a directory tree.
func walkDir(root string, depth int) {
	if depth != 0 {
		files, err := readDir(root)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, finfo := range files {
			fname := finfo.Name()
			fpath := filepath.Join(root, fname)

			if fname[0] != '.' || allFiles {
				if finfo.IsDir() {
					walkDir(fpath, depth-1)
				} else {
					if glob == "" || matchGlob(fpath) {
						wg.Add(1)
						go edit(fpath)
					}
				}
			}
		}
	}
}

func usage() {
	var msg = `red - Recursive Editor
Red allows you to replace all the substrings matched by a specified regex in
one or more files.
If it is given a directory as input, it will recursively replace all the
matches in the files of the directory tree.

Usage:
    %s [options] "pattern" "replacement" input-files

Options:
    -p    Print to stdout instead of writing each file.
    -v    Verbose, explain what is being done.
    -g string
          Add a glob the file names must match to be edited.
    -a    Includes hidden files (starting with a dot).
    -l int
          Max depth in a directory tree.
    -h    Prints this help message.
`
	fmt.Printf(msg, os.Args[0])
}

func parseFlags() {
	flag.BoolVar(&prnt, "p", false, "Print to stdout.")
	flag.BoolVar(&verbose, "v", false, "Verbose, explain what is being done.")
	flag.StringVar(&glob, "g", "", "Add a pattern the file names must match to be edited.")
	flag.BoolVar(&allFiles, "a", false, "Includes hidden files (starting with a dot).")
	flag.IntVar(&maxdepth, "l", -1, "Max depth.")
	flag.Usage = usage
	flag.Parse()
}

func main() {
	var (
		pattern string
		files   []string
	)

	parseFlags()
	if flag.NArg() >= 3 {
		pattern = flag.Arg(0)
		repl = []byte(flag.Arg(1))
		files = flag.Args()[2:]
	} else {
		usage()
		return
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println(err)
		return
	}
	re = regex

	for _, f := range files {
		if f == "-" {
			editStdin()
			continue
		}

		finfo, err := os.Stat(f)
		if err != nil {
			fmt.Println(err)
			return
		}

		if finfo.IsDir() {
			walkDir(f, maxdepth)
		} else {
			wg.Add(1)
			go edit(f)
		}
		wg.Wait()
	}
}
