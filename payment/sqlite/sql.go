package sqlite

import (
	"github.com/ifraixedes/go-payments-api-example/payment"
)

func selectionToSQL(s payment.Selection) (string, pymt, []interface{}) {
	var (
		sf      = make([]string, 1)
		p       pymt
		cols    = []interface{}{&p.ID}
		hasData bool
	)

	sf[0] = "id"

	if s.Version {
		sf = append(sf, "version")
		cols = append(cols, &p.Version)
	}

	if s.OrgID {
		sf = append(sf, "organisation_id")
		cols = append(cols, &p.OrgID)
	}

	if s.Type {
		hasData = true
		sf = append(sf, "json_extract(data")
		cols = append(cols, &p.Data)
	}

	//	if s.Attributes.Amount

	if hasData {
		sf = append(sf, ")")
	}

	return "", p, cols
}
