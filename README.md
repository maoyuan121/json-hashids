# json-hashids

[json iterator](https://github.com/json-iterator/go) 的扩展，可以把整型 marshal 为唯一的、不可预测的 id
## 用法

100% 兼容标准库

替换
```go
import "encoding/json"

json.Marshal(&data)
json.Unmarshal(input, &data)
```

为
```go
import "github.com/liamylian/jsonhashids"

var json = NewConfigWithHashIDs("abcdefg", 10)

json.Marshal(&data)
json.Unmarshal(input, &data)
```


## 例子

```go
package main

import(
	"fmt"
	"github.com/liamylian/json-hashids"
	"time"
)

var json = jsonhashids.NewConfigWithHashIDs("abcdefg", 10)

type Book struct {
	Id    int    `json:"id" hashids:"true"`
	Name  string `json:"name"`
}

func main() {
	book := Book {
		Id:          1,
		Name:        "Jane Eyre",
	}
	
	bytes, _ := json.Marshal(book)
	
	// output: {"id":"gYEL5rKBnd","name":"Jane Eyre"}
	fmt.Printf("%s", bytes)
}

```
