#!/bin/sh
echo "Compiling.."
go build

echo "Copying red to /usr/bin"
sudo cp red /usr/bin/

echo "Installing red manual.."
gzip -c red.1 > red.1.gz
sudo cp red.1.gz /usr/share/man/man1/

echo "Done"

