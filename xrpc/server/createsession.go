package xrpcserver

import (
	"github.com/reiver/go-atproto/xrpc/srv"

	"github.com/reiver/go-fck"
	"github.com/reiver/go-opt"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	errNilHttpResponse = fck.Error("nil *http.Response")
	errNilReceiver     = fck.Error("nil receiver")
)

const (
	errIdentifierNotFound = fck.Error("identifier not found")
	errPasswordNotFound   = fck.Error("password not found")
)

// BlueSky's AT-Protocol XRPC com.atproto.server.createSession
//
// Example usage:
//
//	var req xrpcserver.CreateSession
//	{
//		req.Identifier = opt.Something("username.bsky.social")
//		req.Password = opt.Something("password123")
//	}
//
//	const serviceEndpoint string = "https://bsky.social" // <---- you would probably get this value from a DID document.
//	
//	var httpreq *http.Request
//	{
//		var err error
//		
//		httpreq, err = req.NewHttpRequest(serviceEndpoint)
//		if nil != err {
//			return err
//		}
//	}
//	
//	var httpres *http.Response
//	{
//		var httpclient http.Client
//		
//		var err error
//		
//		httpresp, err =  httpclient.Do(httpreq)
//	}
type CreateSession struct {
	Identifier opt.Optional[string]
	Password   opt.Optional[string]
}

var _ json.Marshaler = CreateSession{}

func (receiver CreateSession) MarshalJSON() ([]byte, error) {
	if opt.Nothing[string]() == receiver.Identifier  {
		return nil, errIdentifierNotFound
	}
	if opt.Nothing[string]() == receiver.Password  {
		return nil, errPasswordNotFound
	}

	var data map[string]string = map[string]string{}
	receiver.Identifier.WhenSomething(func(value string){
		data["identifier"] = value
	})
	receiver.Password.WhenSomething(func(value string){
		data["password"] = value
	})

	return json.Marshal(data)
}

func (receiver CreateSession) NewHTTPRequest(serviceEndpoint string) (*http.Request, error) {
	const methodID string = "com.atproto.server.createSession"

	var body io.Reader
	{
		p, err := receiver.MarshalJSON()
		if nil != err {
			return nil, fmt.Errorf("atproto: problem creating JSON for HTTP POST body for AT-Protocol XRPC %q method call: %w", methodID, err)
		}

		body = bytes.NewReader(p)
	}

	const httpmethod string = http.MethodPost
	var url string = xrpcsrv.MethodURL(serviceEndpoint, methodID)

	var httpreq *http.Request
	{
		var err error

		httpreq, err = http.NewRequest(httpmethod, url, body)
		if nil != err {
			return nil, err
		}

		httpreq.Header.Add("Content-Type", "application/json")
	}

	return httpreq, nil
}

type CreateSessionResponse struct {
	DID    opt.Optional[string]
	EMail  opt.Optional[string]
	Handle opt.Optional[string]
	AccessJWT  opt.Optional[string]
	RefreshJWT opt.Optional[string]
}

func (receiver *CreateSessionResponse) Load(data map[string]string) {

	if value, found := data["did"]; found {
		receiver.DID = opt.Something[string](value)
	}
	if value, found := data["email"]; found {
		receiver.EMail = opt.Something[string](value)
	}
	if value, found := data["handle"]; found {
		receiver.Handle = opt.Something[string](value)
	}
	if value, found := data["accessJwt"]; found {
		receiver.AccessJWT = opt.Something[string](value)
	}
	if value, found := data["refreshJwt"]; found {
		receiver.RefreshJWT = opt.Something[string](value)
	}

}

func (receiver *CreateSessionResponse) ConsumeHTTPResponse(httpresp *http.Response) error {
	if nil == receiver {
		return errNilReceiver
	}
	if nil == httpresp {
		return errNilHttpResponse
	}

	if http.StatusOK != httpresp.StatusCode {
		return fck.Errorf("atproto: problem with HTTP request, HTTP status-code not 200 — %s", httpresp.Status)
	}

	var data map[string]string
	{
		defer httpresp.Body.Close()

		err := json.NewDecoder(httpresp.Body).Decode(&data)
		if nil != err {
			return fmt.Errorf("atproto: problem parsing JSON in response: %w", err)
		}
	}

	receiver.Load(data)

	return nil
}
