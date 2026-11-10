package managedspace

import (
	"net/url"

	"github.com/reiver/go-xrpc"
)

// CheckUserAccess asks a ATProto space's managing app whether to authorize a user to read or write.
// This is called by the ATProto space authority (PDS) when processing credential requests, not by applications directly.
//
// NSID: com.atproto.managedspace.checkUserAccess
//
// # When This Is Called
//
// This endpoint is called by the space authority when:
//
//	• A user requests a space credential for a managed space (read access check)
//	• A user writes to a managed space (write access check)
//
// The managing app specified in the space policy receives these queries and returns whether the user is authorized based on its own logic.
//
// # Parameters
//
// The 'host' parameter should be the managing app's server URL.
//
// The 'serviceToken' parameter is the service authentication token from the space authority.
//
// The 'req' parameter contains:
//
//	• Space: The space reference URI (e.g., "at://did:plc:abc123/space/com.example.game/data")
//	• User: The DID of the user requesting access
//	• Access: The type of access, either "read" or "write"
//	• ClientID: Optional attested OAuth client ID (for read checks only)
//
// # Example (Go) - Implementing a managing app endpoint
//
//	// This is server-side code for a managing app
//	func handleCheckUserAccess(w http.ResponseWriter, r *http.Request) {
//		space := r.URL.Query().Get("space")
//		user := r.URL.Query().Get("user")
//		access := r.URL.Query().Get("access")
//		clientID := r.URL.Query().Get("clientId")
//
//		// Custom authorization logic
//		authorized := checkMyAppAuthorization(space, user, access, clientID)
//
//		json.NewEncoder(w).Encode(map[string]bool{
//			"authorized": authorized,
//		})
//	}
//
// # Example (Go) - Calling as a space authority
//
//	var resp managedspace.CheckUserAccessResponse
//	err := managedspace.CheckUserAccess(
//		&resp,
//		"https://game.example.com",
//		serviceAuthToken,            // service auth from space authority
//		managedspace.CheckUserAccessRequest{
//			Space:    "at://did:plc:abc123/space/com.example.game/data",
//			User:     "did:plc:user789",
//			Access:   "read",
//			ClientID: "did:web:app.example.com",  // optional
//		},
//	)
//	if nil != err {
//		log.Fatal(err)
//	}
//	if resp.Authorized {
//		// Issue credential
//	} else {
//		// Deny access
//	}
//
// # HTTP Request
//
//	GET /xrpc/com.atproto.managedspace.checkUserAccess?space=at%3A%2F%2Fdid%3Aplc%3Aabc123%2Fspace%2Fcom.example.game%2Fdata&user=did%3Aplc%3Auser789&access=read HTTP/1.1
//	Host: game.example.com
//	Authorization: Bearer <service-auth-token>
//
// # HTTP Response
//
//	HTTP/1.1 200 OK
//	Content-Type: application/json
//
//	{
//	  "authorized": true
//	}
//
// # Access Types
//
//   - "read": User is requesting a space credential to read data
//   - "write": User has written to their repo and the write needs authorization
//
// # Notes
//
// The managing app has full control over access decisions and can base them on:
//   - Subscription status
//   - Game progression
//   - User reputation or trust scores
//   - Time-based access windows
//   - Any other app-specific criteria
//
// This endpoint should respond quickly as it's in the critical path for
// credential issuance. Consider caching access decisions where appropriate.
func CheckUserAccess(dst any, host string, serviceToken string, req CheckUserAccessRequest) error {

	const nsid string = "com.atproto.managedspace.checkUserAccess"

	var query string = "space=" + url.QueryEscape(req.Space)
	query = query + "&user=" + url.QueryEscape(req.User)
	query = query + "&access=" + url.QueryEscape(req.Access)

	if "" != req.ClientID {
		query = query + "&clientId=" + url.QueryEscape(req.ClientID)
	}

	var requestURL string = xrpc.ConstructURL(host, nsid, query).String()

	return xrpc.AuthorizedQuery(dst, serviceToken, requestURL)
}

// CheckUserAccessRequest contains the parameters for CheckUserAccess.
type CheckUserAccessRequest struct {
	// Space is the space reference URI (e.g., "at://did:plc:abc123/space/com.example.game/data").
	Space string

	// User is the DID of the user requesting access.
	User string

	// Access is the type of access: "read" or "write".
	Access string

	// ClientID is the optional attested OAuth client ID (for read checks only).
	ClientID string
}

// CheckUserAccessResponse is the response from CheckUserAccess.
type CheckUserAccessResponse struct {
	// Authorized indicates whether the managing app grants the access request.
	Authorized bool `json:"authorized"`
}

