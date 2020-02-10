# red
Recursive Editor

Red is an intuitive and fast find & replace cli.
Red replaces all matches of a specified regex with the provided replacement text.
It can be run over single files or entire directories specified in the command line argument.
When red encounters a directory it recursively finds and replaces text in all the files in the directory tree.

Run `red -h` for more options.

## Usage
```bash
red [options] "pattern" "replacement" [input file/dir]
```
