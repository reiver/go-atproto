package repo

import (
	"github.com/reiver/go-xrpc"

	"github.com/reiver/go-atproto/internal/config"
)

func CreateRecord(dst any, bearerToken string, repoName string, collection string, record any) error {

	const nsid string = "com.atproto.repo.createRecord"
	const host string = config.CreateRecordHost

	var url string = xrpc.ConstructURL(host, nsid, "").String()

	var src = map[string]any{
		"repo":repoName,
		"collection":collection,
		"record":record,
	}

	return xrpc.AuthorizedExecute(dst, bearerToken, url, src)
}

type CreateRecordResponse struct {
	URI              string         `json:"uri"`
	CID              string         `json:"cid"`
	ValidationStatus string         `json:"validationStatus"`
	Commit           map[string]any `json:"commit"`
}
