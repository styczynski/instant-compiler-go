# Instant compiler

## External code usage

This project uses [go-project-blueprint](https://github.com/MartinHeinz/go-project-blueprint) as a starter project so
the original Dockerfile/go.mod and config comes from this repository.

I also utilize customized [chroma](https://github.com/alecthomas/chroma) code to provide syntax colouring.

The project itself implements Hindley-Milder typesystem, which were based on various reference papers and implemeneted by me so there's a little change the core code will be similar to what is seen in this reference documents.

The code itself is mine although orignally it was written for Latte compilation.

## Project strcture

- compiler - core compiler features
- errors - error handling
- events_collector - compiler async event colelction system
- events_utils - minor utilities for event collection
- generic_ast - AST implementation (the specification allows parsing and inference algotether)
- parser - parser implemetation
- input_reader - input pipeline abstraction
- printer - pretty-printing tilities
- runner - compiled code runner and tester
- type_checker - extended-HM inference implementation

## Building

To build the project simply call:

```bash
    $ make
```

## Running

You can test bulk of files using glob:

```bash
    $ ./insc_jvm --backend jvm './test/*.ins'
```

Or compile one at the time:

```bash
    $ ./insc_jvm --backend llvm ./tests/a.ins
```

The compiler atomatically detects `*.output` files and runs compiled program agains them to check for output correctnes.
That means thatyou can easily test the executable. You just have to copy all `*.output` files to the specific directory and use './dir/\*.some_extension` glob for compiler input.
