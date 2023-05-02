package xrpcactor

import (
	"github.com/reiver/go-fck"
	"github.com/reiver/go-opt"

	"encoding/json"
)

const (
	errActorNotFound = fck.Error("actor not found")
)

// BlueSky's AT-Protocol XRPC app.bsky.actor.getProfile
type GetProfile struct {
	Actor opt.Optional[string]
}

var _ json.Marshaler = GetProfile{}

func (receiver GetProfile) MarshalJSON() ([]byte, error) {
	if opt.Nothing[string]() == receiver.Actor  {
		return nil, errActorNotFound
	}

        var data map[string]string = map[string]string{}
        receiver.Actor.WhenSomething(func(value string){
                data["actor"] = value
        })

        return json.Marshal(data)
}
