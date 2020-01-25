package main

import (
    "os"
    "fmt"
    "sync"
    "flag"
    "regexp"
    "strings"
    "io/ioutil"

    "github.com/logrusorgru/aurora"
)

var regex string
var replacement string
var prnt bool
var wg sync.WaitGroup

func printErr(a interface{}) {
    fmt.Println(aurora.Red(a).Bold())
}

func die(a interface{}) {
    printErr(a)
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

func writeFile(fpath string, content string) {
    file, err := os.OpenFile(fpath, os.O_WRONLY, 0644)
    if err != nil {
        printErr(err)
        return
    }
    defer file.Close()
    if _, err := file.WriteString(content); err != nil {
        printErr(err)
    }
}

func edit(fpath string) {
    defer wg.Done()
    b, err := ioutil.ReadFile(fpath)
    if err != nil {
        printErr(err)
        return
    }

    // TODO: calculate hash and compare to avoid useless I/O.
    content := string(b)
    matches := findMatches(content)
    fmt.Println(matches)
    if prnt {
        fmt.Print(replace(content, matches...))
    } else {
        writeFile(fpath, replace(content, matches...))
    }
}

// Recursively walks in a directory tree.
func walkDir(root string) {
	files, err := readDir(root)
	if err != nil {
		printErr(err)
		return
	}

	for _, finfo := range files {
        fpath := fmt.Sprintf("%s%s", root, finfo.Name())

        if finfo.IsDir() {
            walkDir(fpath+"/")
        } else {
            wg.Add(1)
            go edit(fpath)
		}
	}
}

// Removes the './' from the beginning of directory names and
// adds a '/' at the end if missing.
func sanitise(name string) string {
    if name != "." && name[:2] == "./" {
        name = name[2:]
    }
    if name[len(name)-1] != '/' {
        name += "/"
    }
    return name
}

func main() {
    var inFile string

    flag.BoolVar(&prnt, "p", false, "Print to stdout")
    flag.Parse()

    if flag.NArg() >= 3 {
        regex = flag.Arg(0)
        replacement = flag.Arg(1)
        inFile = flag.Arg(2)
    } else {
        die("not enough arguments specified")
    }

    finfo, err := os.Stat(inFile)
    if err != nil {
        die(err)
    }

    if finfo.IsDir() {
        walkDir(sanitise(inFile))
    } else {
        wg.Add(1)
        go edit(inFile)
    }
    wg.Wait()
}
