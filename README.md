

# Instant compiler

[![Build test release](https://github.com/styczynski/instant-compiler-go/actions/workflows/release.yml/badge.svg)](https://github.com/styczynski/instant-compiler-go/actions/workflows/release.yml)

This project uses [go-project-blueprint](https://github.com/MartinHeinz/go-project-blueprint) as a starter project so
the original Dockerfile/go.mod and config comes from this repository.

I also utilize customized [chroma](https://github.com/alecthomas/chroma) code to provide syntax colouring.

The project itself implements Hindley-Milder typesystem, which were based on various reference papers and implemeneted by me so there's a little change the core code will be similar to what is seen in this reference documents.

The code itself is mine although orignally it was written for Latte compilation.

## Project strcture

- compiler - core compiler features (this is only frontend release, so the compiler implementation is just a no-op dummy)
- errors - error handling
- compiler_pipeline - Compiler pipeline code (this folder contains all blocks like runner, compiler, type checker connected)
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

## Interactive CLI

You can run interactive CLI by typing:
```bash
   $ ./latc_jvm shell
```

The interactive shell allows user to provide input program and realtime compiler output as well as running program output.

## Running

You can test bulk of files using glob:

```bash
    $ ./latc_jvm --backend jvm './test/*.lat'
```

Or compile one at the time:

```bash
    $ ./latc_jvm --backend llvm './tests/good/*.lat'
```

The compiler atomatically detects `*.valid_output` files and runs compiled program agains them to check for output correctnes.
That means thatyou can easily test the executable. You just have to copy all `*.valid_output` files to the specific directory and use './dir/\*.some_extension` glob for compiler input.
The extension is configrable by using `--runner-test-extension` parameter. Please type `./latc_jvm --help` or `./latc_x86_64 --help` to get more info.

For example if you wish to test all files in `tests/` diretory you can execute the following command:
```bash
    $ ./latc_x86_64 --runner-test-extension output "./tests/core*.lat"
    $ ./latc_jvm --runner-test-extension output "./tests/good/core*.lat"
```

