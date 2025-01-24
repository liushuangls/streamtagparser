# streamtagparser
Parse the tag data in the stream data returned by the LLM, such as the Artifact

[![Go Reference](https://pkg.go.dev/badge/github.com/liushuangls/streamtagparser/v2.svg)](https://pkg.go.dev/github.com/liushuangls/streamtagparser)
[![Go Report Card](https://goreportcard.com/badge/github.com/liushuangls/streamtagparser)](https://goreportcard.com/report/github.com/liushuangls/streamtagparser)
[![codecov](https://codecov.io/gh/liushuangls/streamtagparser/graph/badge.svg?token=U7TovXlix3)](https://codecov.io/gh/liushuangls/streamtagparser)
[![Sanity check](https://github.com/liushuangls/streamtagparser/actions/workflows/pr.yml/badge.svg)](https://github.com/liushuangls/streamtagparser/actions/workflows/pr.yml)

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
