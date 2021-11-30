# Latte compiler

[![Build test release](https://github.com/styczynski/instant-compiler-go/actions/workflows/release.yml/badge.svg)](https://github.com/styczynski/instant-compiler-go/actions/workflows/release.yml)

## External code usage 

This project uses [go-project-blueprint](https://github.com/MartinHeinz/go-project-blueprint) as a starter project so
the original Dockerfile/go.mod and config comes from this repository.

I also utilize customized [chroma](https://github.com/alecthomas/chroma) code to provide syntax colouring.

The project itself implements Hindley-Milder typesystem and generic CFG code, which were based on various reference papers and implemeneted by me so there's a little change the core code will be similar to what is seen in this reference documents.

## Building

To build the project simply call:
```bash
    $ make
```

## Running 

You can test bulk of files using glob:
```bash
    $ ./latc --short './examples/bad/*.lat'
```

Or compile one at the time:
```bash
    $ ./latc ./examples/good/core001.lat
```
