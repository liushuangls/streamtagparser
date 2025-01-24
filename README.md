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
	"fmt"

	"github.com/liushuangls/streamtagparser"
)

func main() {
	// TagParser Non-concurrency safe, a TagParser can only be used for one stream
	parser := streamtagparser.NewTagParser("Artifact")

	mockStream := []string{
		"hello", "world <A", "rtifact", ` "id"=1>`,
		"local a", "=1", "</Ar", "tifact>end",
	}

	for _, data := range mockStream {
		tagStreams := parser.Parse(data)
		for _, tagStream := range tagStreams {
			fmt.Printf("%#v\n", tagStream)
		}
	}

	// Output:
	/*
		&streamtagparser.TagStreamData{Type:"text", Text:"hello", TagName:"", Attrs:[]streamtagparser.TagAttr(nil), Content:""}
		&streamtagparser.TagStreamData{Type:"text", Text:"world ", TagName:"", Attrs:[]streamtagparser.TagAttr(nil), Content:""}
		&streamtagparser.TagStreamData{Type:"start", Text:"", TagName:"Artifact", Attrs:[]streamtagparser.TagAttr{streamtagparser.TagAttr{Name:"id", Value:"1"}}, Content:""}
		&streamtagparser.TagStreamData{Type:"content", Text:"", TagName:"Artifact", Attrs:[]streamtagparser.TagAttr(nil), Content:"local a"}
		&streamtagparser.TagStreamData{Type:"content", Text:"", TagName:"Artifact", Attrs:[]streamtagparser.TagAttr(nil), Content:"=1"}
		&streamtagparser.TagStreamData{Type:"end", Text:"", TagName:"Artifact", Attrs:[]streamtagparser.TagAttr{streamtagparser.TagAttr{Name:"id", Value:"1"}}, Content:"local a=1"}
		&streamtagparser.TagStreamData{Type:"text", Text:"end", TagName:"", Attrs:[]streamtagparser.TagAttr(nil), Content:""}
	*/

	// Sometimes the streaming data returned by the LLM may be incomplete, so finishing touches are needed
	tagStreams := parser.ParseDone()
	for _, tagStream := range tagStreams {
		fmt.Printf("%#v\n", tagStream)
	}
}
```
