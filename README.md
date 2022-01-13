# Manaita

Simple Markdown-Driven Code Generator written by Go

![manaita](./docs/manaita.png "manaita")

# Installation

```shell
go install github.com/igtm/manaita
```

# Getting Started

1. place sample [CODEGEN.md](./docs/CODEGEN.md) on your directory you like
2. Run `manaita`


# Usage

1. Place `CODEGEN.md` on your directory you like
2. Write scaffold code on it
3. Run `manaita`

Available options:

```
  -c                  specify markdown template file path. default name is 'CODEGEN.md'
  -p                  specify parameters for code gen. these must be defined on markdown  e.g. '-p foo=bar,fizz=buzz'
```

Available template params:

```
  {{.ENV}}            can access environment variables
  {{.Params}}         can access given parameters by '-p' option, which must be defined on 'Params' field of markdown header.
```

# references

- [example code](./example)
