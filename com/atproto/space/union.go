package space

import (
	"bytes"
	"encoding/json"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
)

// marshalWithType returns the JSON of 'value' (which must marshal to a JSON object), with a "$type" field of 'typ' added first.
func marshalWithType(typ string, value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if nil != err {
		return nil, err
	}
	if len(data) < 2 || '{' != data[0] {
		return nil, erorr.Wrap(ErrBadUnion, "union member is not a JSON object", field.String("type", typ))
	}

	typeJSON, err := json.Marshal(typ)
	if nil != err {
		return nil, err
	}

	var buffer bytes.Buffer
	buffer.WriteString(`{"$type":`)
	buffer.Write(typeJSON)
	if "{}" != string(data) {
		buffer.WriteByte(',')
	}
	buffer.Write(data[1:])
	return buffer.Bytes(), nil
}

// unionType returns the "$type" of the JSON object 'data'.
func unionType(data []byte) (string, error) {
	var wire struct {
		Type *string `json:"$type"`
	}
	if err := json.Unmarshal(data, &wire); nil != err {
		return "", erorr.Wrap(ErrBadUnion, "union member is not a JSON object", field.Error("error", err))
	}
	if nil == wire.Type {
		return "", ErrMissingType
	}
	return *wire.Type, nil
}

// countSet returns how many of 'set' are true.
func countSet(set ...bool) int {
	var count int
	for _, b := range set {
		if b {
			count++
		}
	}
	return count
}

// errNotOneMember is returned when marshaling a union that has no member set, or more than one.
func errNotOneMember(union string, count int) error {
	return erorr.Wrap(ErrBadUnion, "exactly one member of the union must be set", field.String("union", union), field.Int("set", count))
}

// copyRaw returns a copy of 'data' (json.Unmarshaler implementations must not keep the slice they are given).
func copyRaw(data []byte) json.RawMessage {
	return append(json.RawMessage(nil), data...)
}
