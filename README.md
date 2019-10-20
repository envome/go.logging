# GoLang Colorful Logging

## Example
```go
package main

import (
	"envo.me/logging"
)

func main() {
	logger := logging.GetLogger("main", &logging.Options{
		Color:logging.BGreen,
		OutputToTerminal: true,
		Colorful:true,
	})

	logger.Println("Start...")
}

```

Output
```bash
[        main  ] app 2019/10/20 10:17:17 Start...
```