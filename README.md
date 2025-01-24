# streamtagparser
Parse the tag data in the stream data returned by the LLM, such as the Artifact

## Installation
```bash
go get github.com/liushuangls/streamtagparser
```

## Usage
```go
package main

import (
	"github.com/liushuangls/streamtagparser"
)

// TagParser Non-concurrency safe, a TagParser can only be used for one stream
parser := streamtagparser.NewTagParser("Artifact")

mockStream := []string{
	"hello", "world <A", "rtifact", ` "id"=1>`,
	"local a", "=1", "</Ar", "tifact>end"
}

for _, data := range mockStream {
	tagStrams := parser.Parse(data)
	for _, tagStram := range tagStrams {
		fmt.Println(tagStram)
	}
}

// Sometimes the streaming data returned by the LLM may be incomplete, so finishing touches are needed
tagStreams := parser.ParseDone()
for _, tagStram := range tagStrams {
	fmt.Println(tagStram)
}
```
