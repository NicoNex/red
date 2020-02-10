#!/bin/sh
echo "Installing red..."
cp red /usr/bin/

echo "Installing red manual..."
gzip -c red.1 > red.1.gz
cp red.1.gz /usr/share/man/man1/

echo "Done"
