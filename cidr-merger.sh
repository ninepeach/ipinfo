#!/usr/bin/env zsh

rm -f ./data/*.lite
./data/ | while read file; do /usr/local/bin/cidr-merger  ./data/$file > ./data/$file.lite ;done
