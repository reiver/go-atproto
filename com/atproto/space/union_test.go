package space_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/reiver/go-atproto/com/atproto/managedspace"
	"github.com/reiver/go-atproto/com/atproto/memberspace"
	"github.com/reiver/go-atproto/com/atproto/space"
)

func TestListWritersOutput_reference(t *testing.T) {
	// In the shape the reference PDS writes it.
	const data = `{"repos":[` +
		`{"$type":"com.atproto.space.defs#repoWrite","repo":"did:plc:aaaaaaaaaaaaaaaaaaaaaaaa","repoRev":"3lxyzabc2222","data":{"$bytes":"AAEC"},"spaceRev":"3lxyzabc3333"},` +
		`{"$type":"com.atproto.space.defs#repoRemoved","repo":"did:plc:bbbbbbbbbbbbbbbbbbbbbbbb","spaceRev":"3lxyzabc4444"},` +
		`{"$type":"com.example.somethingNew","x":1}` +
		`],"cursor":"3lxyzabc5555"}`

	var output space.ListWritersOutput
	if err := json.Unmarshal([]byte(data), &output); nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}

	if "3lxyzabc5555" != output.Cursor {
		t.Errorf("Unexpected cursor: %q", output.Cursor)
	}
	if 3 != len(output.Repos) {
		t.Fatalf("Expected 3 repos but actually got %d.", len(output.Repos))
	}

	write := output.Repos[0]
	if nil == write.RepoWrite || nil != write.RepoRemoved || nil != write.Unknown {
		t.Fatalf("Expected only RepoWrite to be set: %#v", write)
	}
	if "did:plc:aaaaaaaaaaaaaaaaaaaaaaaa" != write.RepoWrite.Repo || "3lxyzabc2222" != write.RepoWrite.RepoRev || "\x00\x01\x02" != string(write.RepoWrite.Data) || "3lxyzabc3333" != write.RepoWrite.SpaceRev {
		t.Errorf("Unexpected RepoWrite: %#v", write.RepoWrite)
	}
	if space.TypeRepoWrite != write.Type() {
		t.Errorf("Unexpected type: %q", write.Type())
	}

	removed := output.Repos[1]
	if nil == removed.RepoRemoved || nil != removed.RepoWrite || nil != removed.Unknown {
		t.Fatalf("Expected only RepoRemoved to be set: %#v", removed)
	}
	if "did:plc:bbbbbbbbbbbbbbbbbbbbbbbb" != removed.RepoRemoved.Repo || "3lxyzabc4444" != removed.RepoRemoved.SpaceRev {
		t.Errorf("Unexpected RepoRemoved: %#v", removed.RepoRemoved)
	}
	if space.TypeRepoRemoved != removed.Type() {
		t.Errorf("Unexpected type: %q", removed.Type())
	}

	unknown := output.Repos[2]
	if nil == unknown.Unknown || nil != unknown.RepoWrite || nil != unknown.RepoRemoved {
		t.Fatalf("Expected only Unknown to be set: %#v", unknown)
	}
	if "com.example.somethingNew" != unknown.Type() {
		t.Errorf("Unexpected type: %q", unknown.Type())
	}

	// And back.
	again, err := json.Marshal(output)
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}
	if data != string(again) {
		t.Errorf("The JSON did not round-trip.\nEXPECTED: %s\nACTUAL:   %s", data, again)
	}
}

func TestNotifySpaceUpdateInput(t *testing.T) {
	tests := []struct {
		JSON  string
		Type  string
		Check func(space.NotifySpaceUpdateUpdate) bool
	}{
		{
			JSON: `{"space":"at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/self","token":"secret","update":{"$type":"com.atproto.space.defs#repoWrite","repo":"did:plc:bbbbbbbbbbbbbbbbbbbbbbbb","repoRev":"3lxyzabc2222","data":{"$bytes":"AAEC"},"spaceRev":"3lxyzabc4444","prevSpaceRev":"3lxyzabc3333"}}`,
			Type: space.TypeRepoWrite,
			Check: func(u space.NotifySpaceUpdateUpdate) bool {
				return nil != u.RepoWrite && "3lxyzabc3333" == u.RepoWrite.PrevSpaceRev
			},
		},
		{
			JSON: `{"space":"at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/self","token":"secret","update":{"$type":"com.atproto.space.defs#repoRemoved","repo":"did:plc:bbbbbbbbbbbbbbbbbbbbbbbb","spaceRev":"3lxyzabc4444"}}`,
			Type: space.TypeRepoRemoved,
			Check: func(u space.NotifySpaceUpdateUpdate) bool {
				return nil != u.RepoRemoved && "" == u.RepoRemoved.PrevSpaceRev
			},
		},
		{
			JSON:  `{"space":"at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/self","token":"secret","update":{"$type":"com.atproto.space.notifySpaceUpdate#spaceDeleted"}}`,
			Type:  space.TypeNotifySpaceUpdateSpaceDeleted,
			Check: func(u space.NotifySpaceUpdateUpdate) bool { return nil != u.SpaceDeleted },
		},
		{
			JSON: `{"space":"at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/self","token":"secret","update":{"$type":"com.atproto.space.notifySpaceUpdate#registrationExpired","spaceRev":"3lxyzabc4444"}}`,
			Type: space.TypeNotifySpaceUpdateRegistrationExpired,
			Check: func(u space.NotifySpaceUpdateUpdate) bool {
				return nil != u.RegistrationExpired && "3lxyzabc4444" == u.RegistrationExpired.SpaceRev
			},
		},
		{
			JSON:  `{"space":"at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/self","token":"secret","update":{"$type":"com.example.later","y":[1,2]}}`,
			Type:  "com.example.later",
			Check: func(u space.NotifySpaceUpdateUpdate) bool { return nil != u.Unknown },
		},
	}

	for _, test := range tests {
		var input space.NotifySpaceUpdateInput
		if err := json.Unmarshal([]byte(test.JSON), &input); nil != err {
			t.Errorf("For %s, did not expect an error but actually got one: %s", test.JSON, err)
			continue
		}
		if "secret" != input.Token || "at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/self" != input.Space {
			t.Errorf("For %s, unexpected space or token: %#v", test.JSON, input)
		}
		if test.Type != input.Update.Type() {
			t.Errorf("For %s, expected type %q but actually got %q.", test.JSON, test.Type, input.Update.Type())
		}
		if !test.Check(input.Update) {
			t.Errorf("For %s, unexpected update: %#v", test.JSON, input.Update)
		}

		again, err := json.Marshal(input)
		if nil != err {
			t.Errorf("For %s, did not expect an error marshaling but actually got one: %s", test.JSON, err)
			continue
		}
		if test.JSON != string(again) {
			t.Errorf("The JSON did not round-trip.\nEXPECTED: %s\nACTUAL:   %s", test.JSON, again)
		}
	}
}

func TestUnion_missingType(t *testing.T) {
	var update space.NotifySpaceUpdateUpdate
	if err := json.Unmarshal([]byte(`{"repo":"did:plc:bbbbbbbbbbbbbbbbbbbbbbbb","spaceRev":"3lxyzabc4444"}`), &update); !errors.Is(err, space.ErrMissingType) {
		t.Errorf("Expected an error wrapping %q but actually got: %v", space.ErrMissingType, err)
	}

	var policy space.Policy
	if err := json.Unmarshal([]byte(`null`), &policy); !errors.Is(err, space.ErrMissingType) {
		t.Errorf("Expected an error wrapping %q but actually got: %v", space.ErrMissingType, err)
	}

	var repo space.ListWritersRepo
	if err := json.Unmarshal([]byte(`"not an object"`), &repo); !errors.Is(err, space.ErrBadUnion) {
		t.Errorf("Expected an error wrapping %q but actually got: %v", space.ErrBadUnion, err)
	}
}

func TestUnion_marshalNotOneMember(t *testing.T) {
	tests := []any{
		space.ListWritersRepo{},
		space.ListWritersRepo{RepoWrite: &space.RepoWrite{}, RepoRemoved: &space.RepoRemoved{}},
		space.NotifySpaceUpdateUpdate{},
		space.NotifySpaceUpdateUpdate{SpaceDeleted: &space.NotifySpaceUpdateSpaceDeleted{}, Unknown: json.RawMessage(`{"$type":"x"}`)},
		space.ApplyWritesWrite{},
		space.ApplyWritesWrite{Create: &space.ApplyWritesCreate{}, Delete: &space.ApplyWritesDelete{}},
		space.ApplyWritesResult{},
		space.Policy{},
		space.Policy{Personal: nil, Member: &memberspace.Policy{}, Managed: &managedspace.Policy{}},
	}

	for _, test := range tests {
		_, err := json.Marshal(test)
		if !errors.Is(err, space.ErrBadUnion) {
			t.Errorf("For %#v, expected an error wrapping %q but actually got: %v", test, space.ErrBadUnion, err)
		}
	}
}

func TestUnion_marshalUnknownWithoutType(t *testing.T) {
	tests := []any{
		space.ListWritersRepo{Unknown: json.RawMessage(`{"repo":"x"}`)},
		space.NotifySpaceUpdateUpdate{Unknown: json.RawMessage(`{"repo":"x"}`)},
		space.Policy{Unknown: json.RawMessage(`{"publicRead":true}`)},
	}

	for _, test := range tests {
		_, err := json.Marshal(test)
		if !errors.Is(err, space.ErrMissingType) {
			t.Errorf("For %#v, expected an error wrapping %q but actually got: %v", test, space.ErrMissingType, err)
		}
	}
}

func TestUnion_unknownIsCopied(t *testing.T) {
	var data = []byte(`{"$type":"com.example.later","y":1}`)

	var repo space.ListWritersRepo
	if err := repo.UnmarshalJSON(data); nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}
	var update space.NotifySpaceUpdateUpdate
	if err := update.UnmarshalJSON(data); nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}
	var policy space.Policy
	if err := policy.UnmarshalJSON(data); nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}

	// The caller may reuse its buffer.
	for i := range data {
		data[i] = ' '
	}

	const expected = `{"$type":"com.example.later","y":1}`
	for _, actual := range []json.RawMessage{repo.Unknown, update.Unknown, policy.Unknown} {
		if expected != string(actual) {
			t.Errorf("Expected Unknown to be a copy, %s, but actually got %q.", expected, actual)
		}
	}
}

func TestApplyWrites_unions(t *testing.T) {
	var input = space.ApplyWritesInput{
		Space: "at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/self",
		Repo:  "did:plc:aaaaaaaaaaaaaaaaaaaaaaaa",
		Writes: []space.ApplyWritesWrite{
			{Create: &space.ApplyWritesCreate{Collection: "app.example.post", Value: json.RawMessage(`{"$type":"app.example.post","text":"hi"}`)}},
			{Update: &space.ApplyWritesUpdate{Collection: "app.example.post", Rkey: "3lxyzabc2222", Value: json.RawMessage(`{"$type":"app.example.post","text":"hello"}`)}},
			{Delete: &space.ApplyWritesDelete{Collection: "app.example.post", Rkey: "3lxyzabc3333"}},
		},
	}

	const expected = `{"space":"at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/self","repo":"did:plc:aaaaaaaaaaaaaaaaaaaaaaaa","writes":[` +
		`{"$type":"com.atproto.space.applyWrites#create","collection":"app.example.post","value":{"$type":"app.example.post","text":"hi"}},` +
		`{"$type":"com.atproto.space.applyWrites#update","collection":"app.example.post","rkey":"3lxyzabc2222","value":{"$type":"app.example.post","text":"hello"}},` +
		`{"$type":"com.atproto.space.applyWrites#delete","collection":"app.example.post","rkey":"3lxyzabc3333"}` +
		`]}`

	actual, err := json.Marshal(input)
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}
	if expected != string(actual) {
		t.Errorf("The actual JSON is not what was expected.\nEXPECTED: %s\nACTUAL:   %s", expected, actual)
	}

	var again space.ApplyWritesInput
	if err := json.Unmarshal(actual, &again); nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}
	if 3 != len(again.Writes) || nil == again.Writes[0].Create || nil == again.Writes[1].Update || nil == again.Writes[2].Delete {
		t.Errorf("Unexpected writes: %#v", again.Writes)
	}

	const results = `{"results":[` +
		`{"$type":"com.atproto.space.applyWrites#createResult","uri":"at://x/y/z","cid":"bafyreia","validationStatus":"valid"},` +
		`{"$type":"com.atproto.space.applyWrites#updateResult","uri":"at://x/y/w","cid":"bafyreib"},` +
		`{"$type":"com.atproto.space.applyWrites#deleteResult"}` +
		`]}`
	var output space.ApplyWritesOutput
	if err := json.Unmarshal([]byte(results), &output); nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}
	if 3 != len(output.Results) || nil == output.Results[0].CreateResult || "valid" != output.Results[0].CreateResult.ValidationStatus || nil == output.Results[1].UpdateResult || nil == output.Results[2].DeleteResult {
		t.Errorf("Unexpected results: %#v", output.Results)
	}
	if again, err := json.Marshal(output); nil != err || results != string(again) {
		t.Errorf("The JSON did not round-trip (error: %v).\nEXPECTED: %s\nACTUAL:   %s", err, results, again)
	}
}

func TestApplyWrites_closedUnion(t *testing.T) {
	var write space.ApplyWritesWrite
	if err := json.Unmarshal([]byte(`{"$type":"com.atproto.space.applyWrites#upsert","collection":"x"}`), &write); !errors.Is(err, space.ErrUnknownType) {
		t.Errorf("Expected an error wrapping %q but actually got: %v", space.ErrUnknownType, err)
	}

	var result space.ApplyWritesResult
	if err := json.Unmarshal([]byte(`{"$type":"com.atproto.space.applyWrites#create"}`), &result); !errors.Is(err, space.ErrUnknownType) {
		t.Errorf("Expected an error wrapping %q but actually got: %v", space.ErrUnknownType, err)
	}
}

func TestPolicy(t *testing.T) {
	tests := []struct {
		JSON string
		Type string
	}{
		{JSON: `{"$type":"com.atproto.personalspace.policy"}`, Type: space.TypePersonalSpacePolicy},
		{JSON: `{"$type":"com.atproto.memberspace.policy","publicRead":true,"allowedApps":["https://app.example.com/client-metadata.json"]}`, Type: space.TypeMemberSpacePolicy},
		{JSON: `{"$type":"com.atproto.managedspace.policy","managingApp":"did:web:app.example.com#manager"}`, Type: space.TypeManagedSpacePolicy},
		{JSON: `{"$type":"com.example.otherspace.policy","z":true}`, Type: "com.example.otherspace.policy"},
	}

	for _, test := range tests {
		var policy space.Policy
		if err := json.Unmarshal([]byte(test.JSON), &policy); nil != err {
			t.Errorf("For %s, did not expect an error but actually got one: %s", test.JSON, err)
			continue
		}
		if test.Type != policy.Type() {
			t.Errorf("For %s, expected type %q but actually got %q.", test.JSON, test.Type, policy.Type())
		}

		again, err := json.Marshal(policy)
		if nil != err {
			t.Errorf("For %s, did not expect an error marshaling but actually got one: %s", test.JSON, err)
			continue
		}
		if test.JSON != string(again) {
			t.Errorf("The JSON did not round-trip.\nEXPECTED: %s\nACTUAL:   %s", test.JSON, again)
		}
	}

	var member space.Policy
	_ = json.Unmarshal([]byte(tests[1].JSON), &member)
	if nil == member.Member || !member.Member.PublicRead || 1 != len(member.Member.AllowedApps) {
		t.Errorf("Unexpected memberspace policy: %#v", member.Member)
	}
}

func TestPolicyRecord(t *testing.T) {
	var record = space.PolicyRecord{Policy: space.Policy{Managed: &managedspace.Policy{ManagingApp: "did:web:app.example.com#manager"}}}

	const expected = `{"$type":"com.atproto.space.policy","policy":{"$type":"com.atproto.managedspace.policy","managingApp":"did:web:app.example.com#manager"}}`

	actual, err := json.Marshal(record)
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}
	if expected != string(actual) {
		t.Errorf("The actual JSON is not what was expected.\nEXPECTED: %s\nACTUAL:   %s", expected, actual)
	}

	var again space.PolicyRecord
	if err := json.Unmarshal(actual, &again); nil != err || nil == again.Policy.Managed {
		t.Errorf("Did not round-trip (error: %v): %#v", err, again)
	}
}

func TestListRepoOpsOutput_nullable(t *testing.T) {
	const data = `{"ops":[` +
		`{"rev":"3lxyzabc2222","collection":"app.example.post","rkey":"a","cid":"bafyreia","prev":null,"value":{"text":"hi"}},` +
		`{"rev":"3lxyzabc3333","collection":"app.example.post","rkey":"a","cid":null,"prev":"bafyreia"}` +
		`],"commit":{"ver":1,"repo":"did:plc:aaaaaaaaaaaaaaaaaaaaaaaa","space":"at://did:plc:aaaaaaaaaaaaaaaaaaaaaaaa/space/app.example.group/self","rev":"3lxyzabc3333","data":{"$bytes":"AAEC"}}}`

	var output space.ListRepoOpsOutput
	if err := json.Unmarshal([]byte(data), &output); nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}

	if 2 != len(output.Ops) {
		t.Fatalf("Expected 2 ops but actually got %d.", len(output.Ops))
	}
	create, del := output.Ops[0], output.Ops[1]
	if nil == create.CID || "bafyreia" != *create.CID || nil != create.Prev || nil == create.Value {
		t.Errorf("Unexpected create: %#v", create)
	}
	if nil != del.CID || nil == del.Prev || "bafyreia" != *del.Prev || nil != del.Value {
		t.Errorf("Unexpected delete: %#v", del)
	}
	if 1 != output.Commit.Ver || "\x00\x01\x02" != string(output.Commit.Data) || "" != output.Cursor {
		t.Errorf("Unexpected commit or cursor: %#v %q", output.Commit, output.Cursor)
	}

	again, err := json.Marshal(output)
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: %s", err)
	}
	if data != string(again) {
		t.Errorf("The JSON did not round-trip.\nEXPECTED: %s\nACTUAL:   %s", data, again)
	}
}
