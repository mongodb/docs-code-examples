package billing

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// ListLinkedOrgs returns all linked organizations for a given billing organization.
func ListLinkedOrgs(ctx context.Context, sdk admin.InvoicesApi, orgID string, opts ...InvoiceOption) ([]string, error) {
	invoices, err := GetCrossOrgBilling(ctx, sdk, orgID, opts...)
	if err != nil {
		return nil, fmt.Errorf("get cross-org billing: %w", err)
	}

	var linkedOrgs []string
	for orgID := range invoices {
		if orgID != orgID {
			linkedOrgs = append(linkedOrgs, orgID)
		}
	}

	return linkedOrgs, nil
}
