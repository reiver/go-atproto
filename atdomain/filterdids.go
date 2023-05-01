package atdomain

import (
	"strings"
)

func filterDIDs(txtresponses []string) (dids []string) {

	const prefix string = "did="

	for _, resp := range txtresponses {

		if strings.HasPrefix(resp, prefix) {
			var did string = resp[len(prefix):]

			if "" != did {
				dids = append(dids, did)
			}
		}
	}

	return dids
}
