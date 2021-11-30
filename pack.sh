#!/bin/bash

DIR=$(pwd)
rm -rfd *.tar.gz
cd ..
cp -R $DIR $DIR/ps386038
cd $DIR
tar -czvf ps386038.tar.gz ./ps386038
rm -rfd ./ps386038
mv ../ps386038.tar.gz .