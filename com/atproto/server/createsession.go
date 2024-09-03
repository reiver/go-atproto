package server

import (
	"github.com/reiver/go-xrpc"

	"github.com/reiver/go-atproto/internal/config"
)

// CreateSession creates an authentication session.
//
// I.e., it lets you login to Bluesky / AT-protocol.
//
// I.e., NSID = com.atproto.server.createSession
//
// Example usage:
//
//	import "github.com/reiver/go-atproto/com/atproto/server"
//	
//	// ...
//	
//	var identifier string = "joeblow.bsky.social" // <-- replace this value
//	var password   string = "password123"         // <-- replace this value
//	
//	var response CreateSessionResponse
//	
//	err := server.CreateSession(&response, identifier, password)
//	if nil != err {
//		return err
//	}
//	
func CreateSession(dst any, identifier string, password string) error {

	const nsid string = "com.atproto.server.createSession"
	const host string = config.CreateSessionHost

	var url string = xrpc.ConstructURL(host, nsid, "").String()

	var src = map[string]string{
		"identifier":identifier,
		"password":password,
	}

	return xrpc.Execute(dst, url, src)
}

type CreateSessionResponse struct {
	DID             string         `json:"did"`
	Handle          string         `json:"handle"`

	Active          bool           `json:"active"`
	Status          string         `json:"status"`

	EMail           string         `json:"email"`
	EMailAuthFactor bool           `json:"emailAuthFactor"`
	EMailConfirmed  bool           `json:"emailConfirmed"`

	AccessJWT       string         `json:"accessJwt"`
	RefreshJWT      string         `json:"refreshJwt"`

	DIDDoc          map[string]any `json:"didDoc"`
}
