package space_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/reiver/go-xrpc"

	"github.com/reiver/go-atproto/com/atproto/space"
)

type recordedRequest struct {
	Method string
	Path   string
	Query  string
	Header http.Header
	Body   string
}

// newServer returns a test server that records each request, and answers with 'status', 'contentType', and 'body'.
func newServer(t *testing.T, status int, contentType string, body string) (*httptest.Server, *[]recordedRequest) {
	t.Helper()

	var requests []recordedRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		requests = append(requests, recordedRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.RawQuery,
			Header: r.Header.Clone(),
			Body:   string(data),
		})
		if "" != contentType {
			w.Header().Set("Content-Type", contentType)
		}
		w.Header().Set("Signature", "mac=:AAAA:")
		w.WriteHeader(status)
		io.WriteString(w, body)
	}))
	t.Cleanup(server.Close)
	return server, &requests
}

func TestGetRecord(t *testing.T) {
	server, requests := newServer(t, http.StatusOK, "application/json", `{"uri":"at://did:plc:bbbbbbbbbbbbbbbbbbbbbbbb/app.example.post/3lxyzabc3333","cid":"bafyreia","value":{"$type":"app.example.post","n":1.50}}`)

	var header = http.Header{}
	header.Set("Authorization", "Bearer credential")
	header.Set("Atproto-Space-Audience", "did:plc:bbbbbbbbbbbbbbbbbbbbbbbb")

	output, err := space.GetRecord(context.Background(), xrpc.Client{}, server.URL+"/", header, space.GetRecordParams{
		Space:      "at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/self",
		Repo:       "did:plc:bbbbbbbbbbbbbbbbbbbbbbbb",
		Collection: "app.example.post",
		Rkey:       "3lxyzabc3333",
	})
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}

	if 1 != len(*requests) {
		t.Fatalf("Expected 1 request but actually got %d.", len(*requests))
	}
	request := (*requests)[0]
	if http.MethodGet != request.Method {
		t.Errorf("Expected method %q but actually got %q.", http.MethodGet, request.Method)
	}
	if "/xrpc/com.atproto.space.getRecord" != request.Path {
		t.Errorf("Unexpected path: %q", request.Path)
	}
	const expectedQuery = "collection=app.example.post&repo=did%3Aplc%3Abbbbbbbbbbbbbbbbbbbbbbbb&rkey=3lxyzabc3333&space=at%3A%2F%2Fdid%3Aplc%3Aaaaaaaaaaaaaaaaaaaaaaaaa%2Fspace%2Fapp.example.group%2Fself"
	if expectedQuery != request.Query {
		t.Errorf("Unexpected query.\nEXPECTED: %s\nACTUAL:   %s", expectedQuery, request.Query)
	}
	if "Bearer credential" != request.Header.Get("Authorization") || "did:plc:bbbbbbbbbbbbbbbbbbbbbbbb" != request.Header.Get("Atproto-Space-Audience") {
		t.Errorf("Unexpected headers: %v", request.Header)
	}

	if "bafyreia" != output.CID || "at://did:plc:bbbbbbbbbbbbbbbbbbbbbbbb/app.example.post/3lxyzabc3333" != output.URI {
		t.Errorf("Unexpected output: %#v", output)
	}
	// The record is kept exactly as sent (ex: "1.50" is not changed to "1.5").
	if `{"$type":"app.example.post","n":1.50}` != string(output.Value) {
		t.Errorf("Unexpected value: %s", output.Value)
	}
}

func TestParams_URL(t *testing.T) {
	tests := []struct {
		Actual   string
		Expected string
	}{
		{
			Actual:   space.ListRepoOpsParams{Space: "s", Repo: "r"}.URL("https://pds.example.com"),
			Expected: "https://pds.example.com/xrpc/com.atproto.space.listRepoOps?repo=r&space=s",
		},
		{
			Actual:   space.ListRepoOpsParams{Space: "s", Repo: "r", Since: "3lxyzabc2222", Limit: 1000, Cursor: "3lxyzabc2222/4", ExcludeValues: true}.URL("https://pds.example.com/"),
			Expected: "https://pds.example.com/xrpc/com.atproto.space.listRepoOps?cursor=3lxyzabc2222%2F4&excludeValues=true&limit=1000&repo=r&since=3lxyzabc2222&space=s",
		},
		{
			Actual:   space.ListSpacesParams{}.URL("http://localhost:2583"),
			Expected: "http://localhost:2583/xrpc/com.atproto.space.listSpaces",
		},
		{
			Actual:   space.ListSpacesParams{SpaceDID: "did:plc:aaaaaaaaaaaaaaaaaaaaaaaa", SpaceType: "app.example.group"}.URL("https://pds.example.com"),
			Expected: "https://pds.example.com/xrpc/com.atproto.space.listSpaces?spaceDid=did%3Aplc%3Aaaaaaaaaaaaaaaaaaaaaaaaa&spaceType=app.example.group",
		},
		{
			Actual:   space.ListWritersParams{Space: "s", Limit: 5}.URL("https://host.example.com"),
			Expected: "https://host.example.com/xrpc/com.atproto.space.listWriters?limit=5&space=s",
		},
		{
			Actual:   space.GetBlobParams{Space: "s", Repo: "r", CID: "bafkreia"}.URL("https://pds.example.com"),
			Expected: "https://pds.example.com/xrpc/com.atproto.space.getBlob?cid=bafkreia&repo=r&space=s",
		},
	}

	for _, test := range tests {
		if test.Expected != test.Actual {
			t.Errorf("Unexpected URL.\nEXPECTED: %s\nACTUAL:   %s", test.Expected, test.Actual)
		}
	}
}

func TestParams_URL_nsid(t *testing.T) {
	const endpoint = "https://pds.example.com"

	tests := []struct {
		URL  string
		NSID string
	}{
		{URL: space.ExportSpaceMetaParams{}.URL(endpoint), NSID: space.NSIDExportSpaceMeta},
		{URL: space.GetBlobParams{}.URL(endpoint), NSID: space.NSIDGetBlob},
		{URL: space.GetDelegationTokenParams{}.URL(endpoint), NSID: space.NSIDGetDelegationToken},
		{URL: space.GetLatestCommitParams{}.URL(endpoint), NSID: space.NSIDGetLatestCommit},
		{URL: space.GetRecordParams{}.URL(endpoint), NSID: space.NSIDGetRecord},
		{URL: space.GetRepoParams{}.URL(endpoint), NSID: space.NSIDGetRepo},
		{URL: space.GetSpaceParams{}.URL(endpoint), NSID: space.NSIDGetSpace},
		{URL: space.ListBlobsParams{}.URL(endpoint), NSID: space.NSIDListBlobs},
		{URL: space.ListRecordsParams{}.URL(endpoint), NSID: space.NSIDListRecords},
		{URL: space.ListRepoOpsParams{}.URL(endpoint), NSID: space.NSIDListRepoOps},
		{URL: space.ListSpacesParams{}.URL(endpoint), NSID: space.NSIDListSpaces},
		{URL: space.ListWritersParams{}.URL(endpoint), NSID: space.NSIDListWriters},
	}

	for _, test := range tests {
		if expected := endpoint + "/xrpc/" + test.NSID; expected != test.URL {
			t.Errorf("Expected %q but actually got %q.", expected, test.URL)
		}
	}
}

func TestAuthorizeSpaceSync(t *testing.T) {
	server, requests := newServer(t, http.StatusOK, "application/json", `{"credential":"eyJ.cred.sig","spaceRev":"3lxyzabc4444"}`)

	output, err := space.AuthorizeSpaceSync(context.Background(), xrpc.Client{}, server.URL, nil, space.AuthorizeSpaceSyncInput{
		Space:           "at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/self",
		CredentialKey:   "did:key:zDnaexample",
		DelegationToken: "eyJ.dt.sig",
		Notify:          &space.AuthorizeSpaceSyncNotify{Endpoint: "https://mirror.example.com", Token: "secret"},
	})
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}

	request := (*requests)[0]
	if http.MethodPost != request.Method || "/xrpc/com.atproto.space.authorizeSpaceSync" != request.Path || "" != request.Query {
		t.Errorf("Unexpected request: %s %s?%s", request.Method, request.Path, request.Query)
	}
	if "application/json" != request.Header.Get("Content-Type") {
		t.Errorf("Unexpected Content-Type: %q", request.Header.Get("Content-Type"))
	}
	const expectedBody = `{"space":"at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/self","credentialKey":"did:key:zDnaexample","delegationToken":"eyJ.dt.sig","notify":{"endpoint":"https://mirror.example.com","token":"secret"}}`
	if expectedBody != strings.TrimSpace(request.Body) {
		t.Errorf("Unexpected body.\nEXPECTED: %s\nACTUAL:   %s", expectedBody, request.Body)
	}

	if "eyJ.cred.sig" != output.Credential || "3lxyzabc4444" != output.SpaceRev {
		t.Errorf("Unexpected output: %#v", output)
	}
}

func TestNoOutput(t *testing.T) {
	server, requests := newServer(t, http.StatusOK, "", "")

	err := space.NotifySpaceUpdate(context.Background(), xrpc.Client{}, server.URL, nil, space.NotifySpaceUpdateInput{
		Space:  "s",
		Token:  "secret",
		Update: space.NotifySpaceUpdateUpdate{SpaceDeleted: &space.NotifySpaceUpdateSpaceDeleted{}},
	})
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}

	request := (*requests)[0]
	if "/xrpc/com.atproto.space.notifySpaceUpdate" != request.Path {
		t.Errorf("Unexpected path: %q", request.Path)
	}
	const expectedBody = `{"space":"s","token":"secret","update":{"$type":"com.atproto.space.notifySpaceUpdate#spaceDeleted"}}`
	if expectedBody != strings.TrimSpace(request.Body) {
		t.Errorf("Unexpected body.\nEXPECTED: %s\nACTUAL:   %s", expectedBody, request.Body)
	}
}

func TestGetRepo(t *testing.T) {
	server, requests := newServer(t, http.StatusOK, space.ContentTypeSTAR, "\x00STAR bytes\xff")

	response, err := space.GetRepo(context.Background(), xrpc.Client{}, server.URL, nil, space.GetRepoParams{Space: "s", Repo: "r"})
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}
	if "\x00STAR bytes\xff" != string(data) {
		t.Errorf("Unexpected body: %q", data)
	}
	if "mac=:AAAA:" != response.Header.Get("Signature") {
		t.Errorf("Expected the response headers to be kept, but actually got: %v", response.Header)
	}
	if "/xrpc/com.atproto.space.getRepo" != (*requests)[0].Path || "repo=r&space=s" != (*requests)[0].Query {
		t.Errorf("Unexpected request: %s?%s", (*requests)[0].Path, (*requests)[0].Query)
	}
}

func TestImportRepo(t *testing.T) {
	server, requests := newServer(t, http.StatusOK, "", "")

	err := space.ImportRepo(context.Background(), xrpc.Client{}, server.URL, nil, strings.NewReader("\x00STAR bytes\xff"))
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}

	request := (*requests)[0]
	if http.MethodPost != request.Method || "/xrpc/com.atproto.space.importRepo" != request.Path {
		t.Errorf("Unexpected request: %s %s", request.Method, request.Path)
	}
	if space.ContentTypeSTAR != request.Header.Get("Content-Type") {
		t.Errorf("Unexpected Content-Type: %q", request.Header.Get("Content-Type"))
	}
	if "\x00STAR bytes\xff" != request.Body {
		t.Errorf("Unexpected body: %q", request.Body)
	}
}

func TestXRPCError(t *testing.T) {
	server, _ := newServer(t, http.StatusBadRequest, "application/json", `{"error":"RecordNotFound","message":"Could not locate record"}`)

	_, err := space.GetRecord(context.Background(), xrpc.Client{}, server.URL, nil, space.GetRecordParams{Space: "s", Repo: "r", Collection: "c", Rkey: "k"})
	if space.ErrorNameRecordNotFound != xrpc.ErrorName(err) {
		t.Errorf("Expected the XRPC error %q but actually got: %v", space.ErrorNameRecordNotFound, err)
	}

	err = space.DeleteSpace(context.Background(), xrpc.Client{}, server.URL, nil, space.DeleteSpaceInput{Space: "s"})
	if space.ErrorNameRecordNotFound != xrpc.ErrorName(err) {
		t.Errorf("Expected the XRPC error %q but actually got: %v", space.ErrorNameRecordNotFound, err)
	}

	_, err = space.GetBlob(context.Background(), xrpc.Client{}, server.URL, nil, space.GetBlobParams{Space: "s", Repo: "r", CID: "c"})
	if space.ErrorNameRecordNotFound != xrpc.ErrorName(err) {
		t.Errorf("Expected the XRPC error %q but actually got: %v", space.ErrorNameRecordNotFound, err)
	}
}

func TestListWriters_call(t *testing.T) {
	server, requests := newServer(t, http.StatusOK, "application/json", `{"repos":[{"$type":"com.atproto.space.defs#repoRemoved","repo":"did:plc:bbbbbbbbbbbbbbbbbbbbbbbb","spaceRev":"3lxyzabc4444"}],"cursor":"3lxyzabc4444"}`)

	output, err := space.ListWriters(context.Background(), xrpc.Client{}, server.URL, nil, space.ListWritersParams{Space: "s", Cursor: "3lxyzabc2222"})
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}
	if "cursor=3lxyzabc2222&space=s" != (*requests)[0].Query {
		t.Errorf("Unexpected query: %q", (*requests)[0].Query)
	}
	if 1 != len(output.Repos) || nil == output.Repos[0].RepoRemoved || "3lxyzabc4444" != output.Cursor {
		t.Errorf("Unexpected output: %#v", output)
	}
}

func TestOutputBadJSON(t *testing.T) {
	server, _ := newServer(t, http.StatusOK, "application/json", `{"repos":[{"repo":"x"}]}`)

	_, err := space.ListWriters(context.Background(), xrpc.Client{}, server.URL, nil, space.ListWritersParams{Space: "s"})
	if nil == err {
		t.Error("Expected an error (a union member with no $type), but did not get one.")
	}
}

func TestCreateSpace_policy(t *testing.T) {
	server, requests := newServer(t, http.StatusOK, "application/json", `{"space":"at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/3lxyzabc2222"}`)

	var policy space.Policy
	if err := json.Unmarshal([]byte(`{"$type":"com.atproto.personalspace.policy"}`), &policy); nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}

	output, err := space.CreateSpace(context.Background(), xrpc.Client{}, server.URL, nil, space.CreateSpaceInput{SpaceType: "app.example.group", Policy: policy})
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}
	const expectedBody = `{"spaceType":"app.example.group","policy":{"$type":"com.atproto.personalspace.policy"}}`
	if expectedBody != strings.TrimSpace((*requests)[0].Body) {
		t.Errorf("Unexpected body.\nEXPECTED: %s\nACTUAL:   %s", expectedBody, (*requests)[0].Body)
	}
	if "at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/3lxyzabc2222" != output.Space {
		t.Errorf("Unexpected output: %#v", output)
	}
}
