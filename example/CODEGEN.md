---
Title: example code
Description: this is example code
Params:
- hoge
- fuga
---

# hoge.go

```golang
package hoge

import "fmt"

func Hoge() {
	fmt.Println("hoge", "{{.Env.HOGE}}")
	fmt.Println("hoge", "{{.Params.hoge}}")
}
```

# fuga.go

```golang
package hoge

import "fmt"

func Hoge() {
	fmt.Println("fuga")
}
```


# {{.Params.fuga}}/bar.go

```golang
package foo

import "fmt"

func Bar() {
	fmt.Println("{{ .Params.fuga | ToUpper }}")
}
```
