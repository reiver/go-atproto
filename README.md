# go-atproto

Package **atproto** provides an implementation of **BlueSky**'s **AT-Protocol**, for the Go programming language.

## Documention

Online documentation, which includes examples, can be found at: http://godoc.org/github.com/reiver/go-atproto

[![GoDoc](https://godoc.org/github.com/reiver/go-atproto?status.svg)](https://godoc.org/github.com/reiver/go-atproto)

## Examples

Here is an example of how to create a new post on Bluesky —

```golang

import
	"time"

	"github.com/reiver/go-atproto/com/atproto/repo"
	"github.com/reiver/go-atproto/com/atproto/server"
)

// ...

var handle   string = "joeblow.bsky.social" // <--- REPLACE THIS WITH THE 'HANDLE' OF THE USER.
var password string = "password123"         // <--- REPLACE THIS WITH THE 'PASSWORD' OF THE USER.


// Login

var bearerToken string
{
	var dst server.CreateSessionResponse

	err :=  server.CreateSession(&dst, identifier, password)

	if nil != err {
		return return
	}

	bearerToken = dst.AccessJWT
}

// Post

var post map[string]any
{
	when := time.Now().Format("2006-01-02T15:04:05.999Z")

	post = map[string]any{
		"$type":"app.bsky.feed.post",
		"text": "Test post at "+when,
		"createdAt": when,
	}
}

var dst repo.CreateRecordResponse
{
	var repoName   string = handle
	var collection string = "app.bsky.feed.post"

	err := repo.CreateRecord(&dst, bearerToken, repoName, collection, post)

	if nil != err {
		return err
	}
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
