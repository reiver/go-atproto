package space

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/reiver/go-atproto/com/atproto/managedspace"
	"github.com/reiver/go-atproto/com/atproto/memberspace"
	"github.com/reiver/go-atproto/com/atproto/personalspace"
)

// The lexicons in "testdata/lexicons" are checked against the Go types:
// every def has a Go type; field names; required (no omitempty) vs. optional (omitempty); nullable; field types; and union members.

type lexiconMethod struct {
	Params any // the params struct, if a query
	Input  any // the input struct, if the input is JSON
	Output any // the output struct, if the output is JSON (and has fields)

	// The Go function, to check the input and output kinds of methods whose input or output is not JSON.
	Func any
}

var lexiconMethods = map[string]lexiconMethod{
	NSIDApplyWrites:              {Input: ApplyWritesInput{}, Output: ApplyWritesOutput{}},
	NSIDAuthorizeSpaceSync:       {Input: AuthorizeSpaceSyncInput{}, Output: AuthorizeSpaceSyncOutput{}},
	NSIDClearRecords:             {Input: ClearRecordsInput{}, Func: ClearRecords},
	NSIDConfirmRegistration:      {Input: ConfirmRegistrationInput{}, Output: ConfirmRegistrationOutput{}},
	NSIDCreateRecord:             {Input: CreateRecordInput{}, Output: CreateRecordOutput{}},
	NSIDCreateSpace:              {Input: CreateSpaceInput{}, Output: CreateSpaceOutput{}},
	NSIDDeleteRecord:             {Input: DeleteRecordInput{}, Func: DeleteRecord},
	NSIDDeleteSpace:              {Input: DeleteSpaceInput{}, Func: DeleteSpace},
	NSIDExportSpaceMeta:          {Params: ExportSpaceMetaParams{}, Func: ExportSpaceMeta},
	NSIDGetBlob:                  {Params: GetBlobParams{}, Func: GetBlob},
	NSIDGetDelegationToken:       {Params: GetDelegationTokenParams{}, Output: GetDelegationTokenOutput{}},
	NSIDGetLatestCommit:          {Params: GetLatestCommitParams{}, Output: GetLatestCommitOutput{}},
	NSIDGetRecord:                {Params: GetRecordParams{}, Output: GetRecordOutput{}},
	NSIDGetRepo:                  {Params: GetRepoParams{}, Func: GetRepo},
	NSIDGetSpace:                 {Params: GetSpaceParams{}, Output: GetSpaceOutput{}},
	NSIDImportRepo:               {Func: ImportRepo},
	NSIDImportSpaceMeta:          {Func: ImportSpaceMeta},
	NSIDListBlobs:                {Params: ListBlobsParams{}, Output: ListBlobsOutput{}},
	NSIDListRecords:              {Params: ListRecordsParams{}, Output: ListRecordsOutput{}},
	NSIDListRepoOps:              {Params: ListRepoOpsParams{}, Output: ListRepoOpsOutput{}},
	NSIDListSpaces:               {Params: ListSpacesParams{}, Output: ListSpacesOutput{}},
	NSIDListWriters:              {Params: ListWritersParams{}, Output: ListWritersOutput{}},
	NSIDNotifyCredentialsRevoked: {Input: NotifyCredentialsRevokedInput{}, Func: NotifyCredentialsRevoked},
	NSIDNotifyRepoUpdate:         {Input: NotifyRepoUpdateInput{}, Func: NotifyRepoUpdate},
	NSIDNotifySpaceUpdate:        {Input: NotifySpaceUpdateInput{}, Func: NotifySpaceUpdate},
	NSIDPutRecord:                {Input: PutRecordInput{}, Output: PutRecordOutput{}},
	NSIDRemoveWriter:             {Input: RemoveWriterInput{}, Func: RemoveWriter},
	NSIDUpdateSpace:              {Input: UpdateSpaceInput{}, Func: UpdateSpace},
}

// lexiconObjects maps each object def (other than a method's main def) to its Go type.
var lexiconObjects = map[string]any{
	TypeRepoWrite:                            RepoWrite{},
	TypeRepoRemoved:                          RepoRemoved{},
	NSIDDefs + "#commit":                     Commit{},
	TypeApplyWritesCreate:                    ApplyWritesCreate{},
	TypeApplyWritesUpdate:                    ApplyWritesUpdate{},
	TypeApplyWritesDelete:                    ApplyWritesDelete{},
	TypeApplyWritesCreateResult:              ApplyWritesCreateResult{},
	TypeApplyWritesUpdateResult:              ApplyWritesUpdateResult{},
	TypeApplyWritesDeleteResult:              ApplyWritesDeleteResult{},
	NSIDAuthorizeSpaceSync + "#notify":       AuthorizeSpaceSyncNotify{},
	NSIDListRecords + "#record":              ListRecordsRecord{},
	NSIDListRepoOps + "#opEntry":             ListRepoOpsOpEntry{},
	NSIDListSpaces + "#spaceView":            ListSpacesSpaceView{},
	TypeNotifySpaceUpdateSpaceDeleted:        NotifySpaceUpdateSpaceDeleted{},
	TypeNotifySpaceUpdateRegistrationExpired: NotifySpaceUpdateRegistrationExpired{},
	TypePersonalSpacePolicy:                  personalspace.Policy{},
	TypeMemberSpacePolicy:                    memberspace.Policy{},
	TypeManagedSpacePolicy:                   managedspace.Policy{},
}

// lexiconRecords maps each record def to the Go type of its "record" object.
var lexiconRecords = map[string]any{
	NSIDPolicy: PolicyRecord{},
}

type lexiconSchema struct {
	Type       string                   `json:"type"`
	Ref        string                   `json:"ref"`
	Refs       []string                 `json:"refs"`
	Closed     bool                     `json:"closed"`
	Required   []string                 `json:"required"`
	Nullable   []string                 `json:"nullable"`
	Properties map[string]lexiconSchema `json:"properties"`
	Items      *lexiconSchema           `json:"items"`
	Encoding   string                   `json:"encoding"`
	Schema     *lexiconSchema           `json:"schema"`
	Parameters *lexiconSchema           `json:"parameters"`
	Input      *lexiconSchema           `json:"input"`
	Output     *lexiconSchema           `json:"output"`
	Record     *lexiconSchema           `json:"record"`
}

type lexiconDocument struct {
	ID   string                   `json:"id"`
	Defs map[string]lexiconSchema `json:"defs"`
}

func loadLexicons(t *testing.T) []lexiconDocument {
	t.Helper()

	paths, err := filepath.Glob(filepath.Join("testdata", "lexicons", "*.json"))
	if nil != err {
		t.Fatalf("Could not list the lexicons: %s", err)
	}
	if 0 == len(paths) {
		t.Fatal("No lexicons found.")
	}

	var documents []lexiconDocument
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if nil != err {
			t.Fatalf("Could not read %q: %s", path, err)
		}
		var document lexiconDocument
		if err := json.Unmarshal(data, &document); nil != err {
			t.Fatalf("Could not parse %q: %s", path, err)
		}
		documents = append(documents, document)
	}
	return documents
}

func TestLexicons(t *testing.T) {
	var seenMethods = map[string]bool{}
	var seenObjects = map[string]bool{}
	var seenRecords = map[string]bool{}

	for _, document := range loadLexicons(t) {
		for name, def := range document.Defs {
			var id string = document.ID
			if "main" != name {
				id += "#" + name
			}

			switch def.Type {
			case "query", "procedure":
				seenMethods[id] = true
				method, found := lexiconMethods[id]
				if !found {
					t.Errorf("No Go types for the method %q.", id)
					continue
				}
				checkMethod(t, document.ID, id, def, method)
			case "object":
				seenObjects[id] = true
				goValue, found := lexiconObjects[id]
				if !found {
					t.Errorf("No Go type for the object %q.", id)
					continue
				}
				checkObject(t, document.ID, id, def, reflect.TypeOf(goValue))
			case "record":
				seenRecords[id] = true
				goValue, found := lexiconRecords[id]
				if !found {
					t.Errorf("No Go type for the record %q.", id)
					continue
				}
				checkObject(t, document.ID, id, *def.Record, reflect.TypeOf(goValue))
			default:
				t.Errorf("The def %q has a type (%q) this test does not know.", id, def.Type)
			}
		}
	}

	for id := range lexiconMethods {
		if !seenMethods[id] {
			t.Errorf("The method %q is not in the lexicons.", id)
		}
	}
	for id := range lexiconRecords {
		if !seenRecords[id] {
			t.Errorf("The record %q is not in the lexicons.", id)
		}
	}
	for id := range lexiconObjects {
		switch id {
		case TypePersonalSpacePolicy, TypeMemberSpacePolicy, TypeManagedSpacePolicy:
			// Defined in other lexicons.
		default:
			if !seenObjects[id] {
				t.Errorf("The object %q is not in the lexicons.", id)
			}
		}
	}
}

func checkMethod(t *testing.T, nsid string, id string, def lexiconSchema, method lexiconMethod) {
	t.Helper()

	// Params.
	switch {
	case nil == def.Parameters && nil != method.Params:
		t.Errorf("For %q, there are Go params but the lexicon has none.", id)
	case nil != def.Parameters && nil == method.Params:
		t.Errorf("For %q, the lexicon has params but there is no Go params type.", id)
	case nil != def.Parameters:
		checkParams(t, id, *def.Parameters, reflect.TypeOf(method.Params))
	}

	var funcType reflect.Type
	if nil != method.Func {
		funcType = reflect.TypeOf(method.Func)
	}

	// Input.
	switch {
	case nil == def.Input:
		if nil != method.Input {
			t.Errorf("For %q, there is a Go input type but the lexicon has no input.", id)
		}
	case "application/json" == def.Input.Encoding:
		if nil == method.Input {
			t.Errorf("For %q, the lexicon has a JSON input but there is no Go input type.", id)
		} else {
			checkObject(t, nsid, id+" input", *def.Input.Schema, reflect.TypeOf(method.Input))
		}
	case ContentTypeSTAR == def.Input.Encoding:
		if nil == funcType || "io.Reader" != funcType.In(funcType.NumIn()-1).String() {
			t.Errorf("For %q, the input is a STAR file, but the Go function does not take an io.Reader.", id)
		}
	default:
		t.Errorf("For %q, the input encoding (%q) is not one this test knows.", id, def.Input.Encoding)
	}

	// Output.
	switch {
	case nil == def.Output:
		if nil != method.Output {
			t.Errorf("For %q, there is a Go output type but the lexicon has no output.", id)
		}
		checkReturnsOnlyError(t, id, funcType)
	case "application/json" == def.Output.Encoding && 0 == len(def.Output.Schema.Properties):
		if nil != method.Output {
			t.Errorf("For %q, there is a Go output type but the lexicon's output has no fields.", id)
		}
		checkReturnsOnlyError(t, id, funcType)
	case "application/json" == def.Output.Encoding:
		if nil == method.Output {
			t.Errorf("For %q, the lexicon has a JSON output but there is no Go output type.", id)
		} else {
			checkObject(t, nsid, id+" output", *def.Output.Schema, reflect.TypeOf(method.Output))
		}
	default:
		if nil == funcType || 2 != funcType.NumOut() || reflect.TypeOf((*http.Response)(nil)) != funcType.Out(0) {
			t.Errorf("For %q, the output is not JSON (%q), but the Go function does not return an *http.Response.", id, def.Output.Encoding)
		}
	}
}

func checkReturnsOnlyError(t *testing.T, id string, funcType reflect.Type) {
	t.Helper()

	if nil == funcType || 1 != funcType.NumOut() || "error" != funcType.Out(0).String() {
		t.Errorf("For %q, there is no output, but the Go function does not return only an error.", id)
	}
}

func checkParams(t *testing.T, id string, def lexiconSchema, goType reflect.Type) {
	t.Helper()

	var goFields = map[string]reflect.StructField{}
	for i := 0; i < goType.NumField(); i++ {
		var structField reflect.StructField = goType.Field(i)
		if name, found := structField.Tag.Lookup("query"); found {
			goFields[name] = structField
		}
	}

	checkSameNames(t, id+" params", def.Properties, goFields)

	for name, property := range def.Properties {
		structField, found := goFields[name]
		if !found {
			continue
		}
		var expected reflect.Kind
		switch property.Type {
		case "string":
			expected = reflect.String
		case "integer":
			expected = reflect.Int64
		case "boolean":
			expected = reflect.Bool
		default:
			t.Errorf("For %q params %q, the type (%q) is not one this test knows.", id, name, property.Type)
			continue
		}
		if expected != structField.Type.Kind() {
			t.Errorf("For %q params %q, expected a Go %s but actually got %s.", id, name, expected, structField.Type)
		}
	}
}

func checkObject(t *testing.T, nsid string, id string, def lexiconSchema, goType reflect.Type) {
	t.Helper()

	if reflect.Struct != goType.Kind() {
		t.Errorf("For %q, expected a Go struct but actually got %s.", id, goType)
		return
	}

	var goFields = map[string]reflect.StructField{}
	var omitEmpty = map[string]bool{}
	for i := 0; i < goType.NumField(); i++ {
		var structField reflect.StructField = goType.Field(i)
		tag, found := structField.Tag.Lookup("json")
		if !found {
			continue
		}
		parts := strings.Split(tag, ",")
		goFields[parts[0]] = structField
		for _, option := range parts[1:] {
			if "omitempty" == option {
				omitEmpty[parts[0]] = true
			}
		}
	}

	// The policies in the other packages marshal through internal types, so only their own lexicons check them.
	switch goType {
	case reflect.TypeOf(personalspace.Policy{}), reflect.TypeOf(memberspace.Policy{}), reflect.TypeOf(managedspace.Policy{}):
		return
	}

	checkSameNames(t, id, def.Properties, goFields)

	var required = map[string]bool{}
	for _, name := range def.Required {
		required[name] = true
	}
	var nullable = map[string]bool{}
	for _, name := range def.Nullable {
		nullable[name] = true
	}

	for name, property := range def.Properties {
		structField, found := goFields[name]
		if !found {
			continue
		}

		if required[name] && omitEmpty[name] {
			t.Errorf("For %q field %q, the field is required, but the Go field has omitempty.", id, name)
		}
		if !required[name] && !omitEmpty[name] {
			t.Errorf("For %q field %q, the field is optional, but the Go field does not have omitempty.", id, name)
		}

		var fieldType reflect.Type = structField.Type
		if nullable[name] {
			if reflect.Pointer != fieldType.Kind() {
				t.Errorf("For %q field %q, the field is nullable, but the Go field (%s) is not a pointer.", id, name, fieldType)
				continue
			}
			fieldType = fieldType.Elem()
		}

		checkType(t, nsid, id+" field "+name, property, fieldType, required[name])
	}
}

func checkSameNames[T any](t *testing.T, id string, properties map[string]lexiconSchema, goFields map[string]T) {
	t.Helper()

	var lexiconNames []string
	for name := range properties {
		lexiconNames = append(lexiconNames, name)
	}
	var goNames []string
	for name := range goFields {
		goNames = append(goNames, name)
	}
	sort.Strings(lexiconNames)
	sort.Strings(goNames)

	if strings.Join(lexiconNames, ",") != strings.Join(goNames, ",") {
		t.Errorf("For %q, the field names are not the same.\nLEXICON: %v\nGO:      %v", id, lexiconNames, goNames)
	}
}

var (
	typeBytes      = reflect.TypeOf(Bytes(nil))
	typeRawMessage = reflect.TypeOf(json.RawMessage(nil))
)

func checkType(t *testing.T, nsid string, id string, property lexiconSchema, goType reflect.Type, required bool) {
	t.Helper()

	switch property.Type {
	case "string":
		if reflect.String != goType.Kind() {
			t.Errorf("For %q, expected a Go string but actually got %s.", id, goType)
		}
	case "integer":
		if reflect.Int64 != goType.Kind() {
			t.Errorf("For %q, expected a Go int64 but actually got %s.", id, goType)
		}
	case "boolean":
		// An optional boolean is a *bool, so that "not set" differs from false.
		if required && reflect.Bool != goType.Kind() {
			t.Errorf("For %q, expected a Go bool but actually got %s.", id, goType)
		}
		if !required && reflect.TypeOf((*bool)(nil)) != goType {
			t.Errorf("For %q, expected a Go *bool but actually got %s.", id, goType)
		}
	case "bytes":
		if typeBytes != goType {
			t.Errorf("For %q, expected %s but actually got %s.", id, typeBytes, goType)
		}
	case "unknown":
		if typeRawMessage != goType {
			t.Errorf("For %q, expected %s but actually got %s.", id, typeRawMessage, goType)
		}
	case "array":
		if reflect.Slice != goType.Kind() {
			t.Errorf("For %q, expected a Go slice but actually got %s.", id, goType)
			return
		}
		checkType(t, nsid, id+" items", *property.Items, goType.Elem(), true)
	case "ref":
		var ref string = fullRef(nsid, property.Ref)
		expected, found := lexiconObjects[ref]
		if !found {
			t.Errorf("For %q, the ref %q has no Go type.", id, ref)
			return
		}
		if reflect.TypeOf(expected) != goType && reflect.PointerTo(reflect.TypeOf(expected)) != goType {
			t.Errorf("For %q, expected %T (or a pointer to it) but actually got %s.", id, expected, goType)
		}
	case "union":
		checkUnion(t, nsid, id, property, goType)
	default:
		t.Errorf("For %q, the type (%q) is not one this test knows.", id, property.Type)
	}
}

func checkUnion(t *testing.T, nsid string, id string, property lexiconSchema, goType reflect.Type) {
	t.Helper()

	if reflect.Struct != goType.Kind() {
		t.Errorf("For %q, expected a Go union struct but actually got %s.", id, goType)
		return
	}

	var expected []string
	for _, ref := range property.Refs {
		expected = append(expected, reflect.TypeOf(lexiconObjects[fullRef(nsid, ref)]).String())
	}

	var actual []string
	var hasUnknown bool
	for i := 0; i < goType.NumField(); i++ {
		var structField reflect.StructField = goType.Field(i)
		switch {
		case "Unknown" == structField.Name && typeRawMessage == structField.Type:
			hasUnknown = true
		case reflect.Pointer == structField.Type.Kind():
			actual = append(actual, structField.Type.Elem().String())
		default:
			t.Errorf("For %q, the union's Go field %q (%s) is neither a pointer nor Unknown.", id, structField.Name, structField.Type)
		}
	}

	sort.Strings(expected)
	sort.Strings(actual)
	if strings.Join(expected, ",") != strings.Join(actual, ",") {
		t.Errorf("For %q, the union members are not the same.\nLEXICON: %v\nGO:      %v", id, expected, actual)
	}

	if property.Closed && hasUnknown {
		t.Errorf("For %q, the union is closed, but the Go type has an Unknown field.", id)
	}
	if !property.Closed && !hasUnknown {
		t.Errorf("For %q, the union is open, but the Go type has no Unknown field.", id)
	}
}

// fullRef returns the ref 'ref' (ex: "#create", "com.atproto.space.defs#commit") in full, for a ref in the lexicon 'nsid'.
func fullRef(nsid string, ref string) string {
	if strings.HasPrefix(ref, "#") {
		return nsid + ref
	}
	return ref
}
