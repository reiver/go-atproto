package space

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
)

// Bytes is the lexicon type 'bytes'.
//
// In JSON it is an object with only a "$bytes" field, whose value is the bytes in base64 (standard alphabet). For example:
//
//	{"$bytes":"SGVsbG8sIHdvcmxkIQ"}
//
// It is written without padding, and read with or without padding.
type Bytes []byte

var _ json.Marshaler = Bytes(nil)
var _ json.Unmarshaler = new(Bytes)

func (receiver Bytes) MarshalJSON() ([]byte, error) {
	var wire struct {
		Bytes string `json:"$bytes"`
	}
	wire.Bytes = base64.RawStdEncoding.EncodeToString(receiver)

	return json.Marshal(wire)
}

func (receiver *Bytes) UnmarshalJSON(data []byte) error {
	if nil == receiver {
		return erorr.Wrap(ErrBadBytes, "nil receiver")
	}

	var wire map[string]json.RawMessage
	if err := json.Unmarshal(data, &wire); nil != err {
		return erorr.Wrap(ErrBadBytes, "not a JSON object", field.Error("error", err))
	}
	raw, found := wire["$bytes"]
	if !found {
		return erorr.Wrap(ErrBadBytes, "missing \"$bytes\"")
	}
	if 1 != len(wire) {
		return erorr.Wrap(ErrBadBytes, "fields other than \"$bytes\"", field.Int("fields", len(wire)))
	}

	var encoded string
	if err := json.Unmarshal(raw, &encoded); nil != err {
		return erorr.Wrap(ErrBadBytes, "\"$bytes\" is not a string", field.Error("error", err))
	}

	decoded, err := base64.RawStdEncoding.DecodeString(strings.TrimRight(encoded, "="))
	if nil != err {
		return erorr.Wrap(ErrBadBytes, "\"$bytes\" is not base64", field.Error("error", err))
	}

	*receiver = decoded
	return nil
}
