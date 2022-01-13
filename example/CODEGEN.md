---
Title: example code
Description: this is example code
Params:
- hoge
- fuga
---

# ./hoge.go

```golang
package hoge

import "fmt"

func Hoge() {
	fmt.Println("hoge", "{{.Env.HOGE}}")
	fmt.Println("hoge", "{{.Params.hoge}}")
}
```

# ./fuga.go

```golang
package hoge

import "fmt"

func Hoge() {
	fmt.Println("fuga")
}
```


# ./foo/bar.go

```golang
package foo

import "fmt"

func Bar() {
	fmt.Println("bar")
}
```
