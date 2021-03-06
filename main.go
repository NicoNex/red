package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
)

var prnt bool
var repl []byte
var glob string
var verbose bool
var maxdepth int
var allFiles bool
var wg sync.WaitGroup
var re *regexp.Regexp

func die(a interface{}) {
	fmt.Println(a)
	os.Exit(1)
}

func readDir(filename string) ([]os.FileInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return []os.FileInfo{}, err
	}
	defer file.Close()
	return file.Readdir(-1)
}

func matchGlob(fname string) bool {
	if runtime.GOOS == "windows" {
		re := regexp.MustCompile(`[^\\]+$`)
		fname = re.FindString(fname)
	}

	ok, err := filepath.Match(glob, fname)
	if err != nil {
		fmt.Println(err)
	}

	return ok
}

func edit(fpath string) {
	var match = make(chan bool, 1)
	defer wg.Done()

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !prnt {
		go func(cont []byte, ch chan bool) {
			ch <- re.Match(cont)
		}(b, match)
	}

	tmp := re.ReplaceAll(b, repl)
	if prnt {
		fmt.Print(string(tmp))
	} else if <-match {
		if verbose {
			fmt.Printf("Writing %s\n", fpath)
		}
		ioutil.WriteFile(fpath, tmp, 0644)
	}
	close(match)
}

func editStdin() {
	var reader = bufio.NewReader(os.Stdin)

	b, err := reader.ReadBytes(0)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}

	o := re.ReplaceAll(b, repl)
	fmt.Print(string(o))
}

// Recursively walks in a directory tree.
func walkDir(root string, depth int) {
	if depth != maxdepth {
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
					walkDir(fpath, depth+1)
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
Red allows you to replace all the substrings matched by a specified regex in one or more files.
If it is given a directory as input, it will recursively replace all the matches in the files of the directory tree.

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

func main() {
	var pattern string
	var files []string

	flag.BoolVar(&prnt, "p", false, "Print to stdout.")
	flag.BoolVar(&verbose, "v", false, "Verbose, explain what is being done.")
	flag.StringVar(&glob, "g", "", "Add a pattern the file names must match to be edited.")
	flag.BoolVar(&allFiles, "a", false, "Includes hidden files (starting with a dot).")
	flag.IntVar(&maxdepth, "l", -1, "Max depth.")
	flag.Usage = usage
	flag.Parse()

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
		die(err)
	}
	re = regex

	for _, f := range files {
		if f == "-" {
			editStdin()
			continue
		}

		finfo, err := os.Stat(f)
		if err != nil {
			die(err)
		}

		if finfo.IsDir() {
			walkDir(f, 0)
		} else {
			wg.Add(1)
			go edit(f)
		}
		wg.Wait()
	}
}
