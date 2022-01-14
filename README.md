# Manaita

Simple Markdown-Driven Code Generator written by Go

Write your scaffolding code on `CODEGEN.md` and generate files using the scaffold.

Template file is Markdown format. so you can see it on Github and easily understand what will be generated.

![manaita](./docs/manaita.png "manaita")

# Installation

```shell
go install github.com/igtm/manaita@latest
```

# Getting Started

1. put `CODEGEN.md` on your directory

````
---
Params:
- name
---

# foo.go

```golang
package foo

var foo = "foo"
```


# {{.Params.name}}/bar.py

```python
print("bar.py")
```
````


2. Run `manaita -p name=dog`
3. `foo.go` and `dog/bar.py` files are generated

# Usage

```shell
manaita -c ./path/to/CODEGEN.md -p key=value
```

Available options:

```
  -c                  specify markdown template file path. default name is 'CODEGEN.md'
  -p                  specify parameters for code gen. these must be defined on markdown  e.g. '-p foo=bar,fizz=buzz'
```

Available template params:

```
  {{.Env}}            can access environment variables
  {{.Params}}         can access given parameters by '-p' option, which must be defined on 'Params' field of markdown header.
```

# references

- [example code](./example)
