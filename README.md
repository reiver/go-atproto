# go-atproto

Package **atproto** provides an implementation of **BlueSky**'s **AT-Protocol**, for the Go programming language.

## Documention

Online documentation, which includes examples, can be found at: http://godoc.org/github.com/reiver/go-atproto

[![GoDoc](https://godoc.org/github.com/reiver/go-atproto?status.svg)](https://godoc.org/github.com/reiver/go-atproto)

## Example

```go

import "github.com/reiver/go-atproto/atdomain"
import "github.com/reiver/go-atproto/didplc"

import "net/http"


// ...

var domainname string = "example.com"

dids, err := atdomain.LookupDID(domainname)
if nil != err {
	return err
}

// ...

did := dids[0]

uri := didplc.DocURI(did)

// ...

resp, err := http.Get(uri)
if nil != err {
	return err
}

```

## Import

To import package **atproto** use `import` code like the follownig:
```
import "github.com/reiver/go-atproto"
```

## Installation

To install package **atproto** do the following:
```
GOPROXY=direct go get https://github.com/reiver/go-atproto
```

## Author

Package **atproto** was written by [Charles Iliya Krempeaux](http://reiver.link)
