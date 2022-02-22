# Manaita

Simple Markdown-Driven Scaffolding tool written by Go

Write your scaffolding code on `SCAFFOLD.md` and generate files using the scaffold.

Scaffold file is Markdown format. so you can see it on Github and easily understand what will be generated.

![manaita](./docs/manaita.png "manaita")

# Installation

### Brew

```shell
brew install igtm/tap/manaita
```

### Go Install

```shell
go install github.com/igtm/manaita@latest
```

### Install with curl

```shell
curl -sfL https://raw.githubusercontent.com/igtm/manaita/master/install.sh | sudo sh -s -- -b /usr/local/bin
```

# Getting Started

1. put `SCAFFOLD.md` on your directory

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
manaita -c ./path/to/SCAFFOLD.md -p key=value
```

Available options:

```
  -c                  specify markdown scaffold file path. default name is 'SCAFFOLD.md'
                        supported types:
                          - local path         e.g. 'path/to/SCAFFOLD.md'
                          - go get style path: e.g. 'github.com/owner/repo/path/to/SCAFFOLD.md'
                          - http url:          e.g. 'https://example.com/SCAFFOLD.md'
  -p                  specify parameters for scaffold template. these must be defined on markdown  e.g. '-p foo=bar,fizz=buzz'
```

Available template params:

```
  {{.Env}}            can access environment variables
  {{.Params}}         can access given parameters by '-p' option, which must be defined on 'Params' field of markdown header.
```

Available template functions:

`AnyKind of_string`

| Function           | Result               |
|--------------------|----------------------|
| `ToUpper`          | `ANY KIND OF_STRING` |
| `ToLower`          | `anykind of_string`  |
| `ToSnake`          | `any_kind_of_string` |
| `ToScreamingSnake` | `ANY_KIND_OF_STRING` |
| `ToKebab`          | `any-kind-of-string` |
| `ToScreamingKebab` | `ANY-KIND-OF-STRING` |
| `ToCamel`          | `AnyKindOfString`    |
| `ToLowerCamel`     | `anyKindOfString`    |

This library uses [Go Template](https://pkg.go.dev/text/template).

so you can use any Go Template syntax like '{{if foo}} .. {{end}}' and like that.

# References

- [Example Code](./example)

# Trouble Shooting

- cannot get private repo.

it is cloned through `ssh`. so you need setup `~/.ssh/config`.
(currently only @branch is supported on private repo)

- when you get `ssh: handshake failed: knownhosts: key mismatch`.

```shell
ssh-keyscan -H github.com >> ~/.ssh/known_hosts
```
