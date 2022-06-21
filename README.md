# log2zap

The Go linter `log2zap` helps to migrate code from log.Print* calls to uber zap.

For example:

```go
package main

import "log"

func main() {
	log.Printf("set some var=%d err=%v", 215, nil)
	// => zap.L().Info("set some", zap.Int("var", 215), zap.Any("err", nil))
}
```

*You can treat it as an example project and modify it for your case.*