[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/) [![Go Report Card](https://goreportcard.com/badge/github.com/NicoNex/red)](https://goreportcard.com/report/github.com/NicoNex/red) [![License](http://img.shields.io/badge/license-GPL3.0-orange.svg?style=flat)](https://github.com/NicoNex/re/blob/master/LICENSE)

## Warning
This repository was **moved** to **[Jet](https://github.com/NicoNex/jet)** on **28-07-2023** and thus is **unmantained**.  
The name was changed so that it no longer collides with many other Unix tools.  
For the mantained version please refer to: https://github.com/NicoNex/jet

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
