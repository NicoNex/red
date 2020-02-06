package main

import (
	"os"
	"fmt"
	"sync"
	"flag"
	"regexp"
	"io/ioutil"
)

var prnt bool
var repl string
var maxdepth int
var pattern string
var editHidden bool
var wg sync.WaitGroup

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

func replace(src string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		die(err)
	}
	return re.ReplaceAllString(src, repl)
}

func edit(fpath string) {
	defer wg.Done()
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Println(err)
		return
	}
	content := string(b)
	if prnt {
		fmt.Print(replace(content))
	} else if ok, _ := regexp.Match(pattern, b); ok {
		ioutil.WriteFile(fpath, []byte(replace(content)), 0644)
	}
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
			fpath := fmt.Sprintf("%s%s", root, fname)

			if fname[0] != '.' || editHidden {
				if finfo.IsDir() {
					walkDir(fpath+"/", depth+1)
				} else {
					wg.Add(1)
					go edit(fpath)
				}
			}
		}
	}
}

func usage() {
	var msg = `red - Recursive Editor
Red allows you to replace all the substrings matched by a specified regex in one or more files.
If it is given a directory as input, it will recursively replace the substrings in all the files of the directory.

Usage:
    %s [options] "regex" "replacement" input-files

Options:
    -p    Print to stdout instead of writing each file.
    -d    Includes hidden files (starting with a dot).
    -l int
          Max depth in a directory tree.
`
	fmt.Printf(msg, os.Args[0])
}

func main() {
	var files []string

	flag.BoolVar(&prnt, "p", false, "Print to stdout")
	flag.BoolVar(&editHidden, "d", false, "Includes hidden files (starting with a dot).")
	flag.IntVar(&maxdepth, "l", -1, "Max depth")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() >= 3 {
		pattern = flag.Arg(0)
		repl = flag.Arg(1)
		files = flag.Args()[2:]
	} else {
		usage()
		return
	}

	for _, f := range files {
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
