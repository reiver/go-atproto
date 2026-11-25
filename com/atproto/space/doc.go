/*
Package space provides the "com.atproto.space.*" XRPC methods — ATProto permissioned spaces.

For each method there are:
the input / output (or params / output) types, which are the method's JSON;
and a function that calls the method using an [xrpc.Client].

For example:

	var client xrpc.Client

	var header = http.Header{}
	header.Set("Authorization", "Bearer "+spaceCredential)

	output, err := space.GetRecord(ctx, client, "https://pds.example.com", header, space.GetRecordParams{
		Space:      "at://did:plc:ewvi7nxzyoun6zhxrhs64oiz/space/app.example.group/3lxyzabc2222",
		Repo:       "did:plc:44ybard66vv44zksje25o7dz",
		Collection: "app.example.post",
		Rkey:       "3lxyzabc3333",
	})

'endpoint' is the base URL of the service (ex: "https://pds.example.com") — such as the "#atproto_pds" service endpoint of a DID document.

The functions return an [*xrpc.Error] when the service answers with an XRPC error. Its name is one of the ErrorName constants (ex: [ErrorNameSpaceNotFound]):

	if xrpc.ErrorName(err) == space.ErrorNameRecordNotFound {
		// ...
	}

Some responses are authenticated with a MAC (requested with the "Accept-Signature" header) —
getRecord, getRepo, getLatestCommit, and listRepoOps.
To check the MAC the caller needs the response's headers and body as is.
For that, use the params' URL method with [xrpc.Client.Query], and decode the body into the output type.

The lexicons these types follow are in the "testdata/lexicons" directory.
*/
package space
