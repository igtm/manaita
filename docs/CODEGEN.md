---
Title: example code
Description: this is example code
Params:
- name
---

# foo.go

```golang
package foo

import "fmt"

func Foo() {
	// you can use Params defined above. pass options like '-p name=George'
	fmt.Println("name", "{{.Params.name}}")
	// you can also use Environment Variable. 
	fmt.Println("BAR", "{{.Env.BAR}}")
}
```
