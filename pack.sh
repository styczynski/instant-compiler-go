#!/bin/bash

DIR=$(PWD)
git add . && git reset --hard
git clean -ffdx
rm -rfd *.tar.gz
cd ..
tar -czvf ps386038.tar.gz *ps386038
cd $DIR
mv ../ps386038.tar.gz .