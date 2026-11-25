package space

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// methodURL returns the URL of the XRPC method 'nsid' at 'endpoint' (ex: "https://pds.example.com"), with the (already encoded) query 'query'.
func methodURL(endpoint string, nsid string, query string) string {
	var result string = strings.TrimSuffix(endpoint, "/") + "/xrpc/" + nsid
	if "" != query {
		result += "?" + query
	}
	return result
}

// encodeQuery returns the query string for 'params' — a struct whose fields have `query:"name"` tags.
// Fields with their zero value are left out.
func encodeQuery(params any) string {
	var values url.Values = url.Values{}

	var value reflect.Value = reflect.ValueOf(params)
	var typ reflect.Type = value.Type()
	for i := 0; i < typ.NumField(); i++ {
		name, found := typ.Field(i).Tag.Lookup("query")
		if !found {
			continue
		}

		var fieldValue reflect.Value = value.Field(i)
		if fieldValue.IsZero() {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.String:
			values.Set(name, fieldValue.String())
		case reflect.Int, reflect.Int64:
			values.Set(name, strconv.FormatInt(fieldValue.Int(), 10))
		case reflect.Bool:
			values.Set(name, strconv.FormatBool(fieldValue.Bool()))
		default:
			panic("space: unsupported query field type: " + fieldValue.Kind().String())
		}
	}

	return values.Encode()
}
