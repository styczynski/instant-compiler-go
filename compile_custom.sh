#!/bin/bash

DIR="/usr/src/comp"

#rm -rf ${DIR}/custom > /dev/null 2> /dev/null
#rm -rf ${DIR}/custom.o > /dev/null 2> /dev/null
#rm -rf ${DIR}/runtime.o > /dev/null 2> /dev/null

#gcc -c ${DIR}/lib/runtime.s -o ${DIR}/runtime.o
#gcc -c ${DIR}/custom.s -o ${DIR}/custom.o

#gcc ${DIR}/custom.o ${DIR}/runtime.o -o ${DIR}/custom
#gcc ${DIR}/custom.o -o ${DIR}/custom
echo "=== Run program ==="
${DIR}/custom
echo $?