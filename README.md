[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/) [![Go Report Card](https://goreportcard.com/badge/github.com/NicoNex/red)](https://goreportcard.com/report/github.com/NicoNex/red) [![License](http://img.shields.io/badge/license-GPL3.0-orange.svg?style=flat)](https://github.com/NicoNex/re/blob/master/LICENSE)

# red
Recursive Editor

Red is an intuitive and fast find & replace cli.
Red replaces all matches of a specified regex with the provided replacement text.
It can be run over single files or entire directories specified in the command line argument.
When red encounters a directory it recursively finds and replaces text in all the files in the directory tree.

Install it with `go install github.com/NicoNex/red@latest`.
Run `red -h` for more options.

## Usage
```bash
red [options] "pattern" "replacement" [input file/dir]
```
