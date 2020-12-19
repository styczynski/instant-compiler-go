#!/bin/bash

git add . && git reset --hard
git clean -ffdx
rm -rfd *.tar.gz
tar -czvf ps386038.tar.gz *ps386038
