package main

import (
	"os"
	"fmt"
	"sync"
	"flag"
	"regexp"
	"strings"
	"io/ioutil"
)

var prnt bool
var regex string
var maxdepth int
var editHidden bool
var wg sync.WaitGroup
var replacement string

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

func removeDuplicates(a []string) []string {
	var ret []string

	m := make(map[string]bool)
	for _, v := range a {
		if _, ok := m[v]; !ok {
			m[v] = true
			ret = append(ret, v)
		}
	}
	return ret
}

func findMatches(cnt string) []string {
	re := regexp.MustCompile(regex)
	matches := re.FindAllString(cnt, -1)
	return removeDuplicates(matches)
}

func replace(in string, v ...string) string {
	var elems []string

	for _, s := range v {
		elems = append(elems, s, replacement)
	}

	r := strings.NewReplacer(elems...)
	return r.Replace(in)
}

func edit(fpath string) {
	defer wg.Done()
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Println(err)
		return
	}
	content := string(b)
	matches := findMatches(content)

	if prnt {
		fmt.Print(replace(content, matches...))
	} else if len(matches) > 0 {
		tmp := replace(content, matches...)
		ioutil.WriteFile(fpath, []byte(tmp), 0644)
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

// TODO: add hidden flag
func main() {
	var files []string

	flag.BoolVar(&prnt, "p", false, "Print to stdout")
	flag.BoolVar(&editHidden, "d", false, "Includes hidden files (starting with a dot).")
	flag.IntVar(&maxdepth, "l", -1, "Max depth")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() >= 3 {
		regex = flag.Arg(0)
		replacement = flag.Arg(1)
		files = flag.Args()[2:]
	} else if flag.NArg() == 2 {
		regex = flag.Arg(0)
		replacement = flag.Arg(1)
		files = []string{"./"}
	} else {
		die("not enough arguments specified")
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
