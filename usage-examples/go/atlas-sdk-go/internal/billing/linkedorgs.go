package billing

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// GetLinkedOrgs returns all linked organizations for a given billing organization.
func GetLinkedOrgs(ctx context.Context, sdk admin.InvoicesApi, p *admin.ListInvoicesApiParams) ([]string, error) {
	invoices, err := GetCrossOrgBilling(ctx, sdk, p)
	if err != nil {
		return nil, fmt.Errorf("get cross-org billing: %w", err)
	}

	var linkedOrgs []string
	for orgID := range invoices {
		if orgID != p.OrgId {
			linkedOrgs = append(linkedOrgs, orgID)
		}
	}

	return linkedOrgs, nil
}
