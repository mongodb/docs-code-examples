package billing

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// GetCrossOrgBilling returns a map of all billing invoices for the given organization and any linked organizations, grouped by organization ID.
// NOTE: Organization Billing Admin or Organization Owner role required to view linked invoices.
func GetCrossOrgBilling(ctx context.Context, sdk admin.InvoicesApi, orgId string, opts ...InvoiceOption) (map[string][]admin.BillingInvoiceMetadata, error) {
	r, err := ListInvoicesForOrg(ctx, sdk, orgId, opts...)
	if err != nil {
		return nil, err
	}

	crossOrgBilling := make(map[string][]admin.BillingInvoiceMetadata)
	if r == nil || !r.HasResults() || len(r.GetResults()) == 0 {
		return crossOrgBilling, nil
	}

	crossOrgBilling[orgId] = r.GetResults()
	for _, invoice := range r.GetResults() {
		if !invoice.HasLinkedInvoices() || len(invoice.GetLinkedInvoices()) == 0 {
			continue
		}

		for _, linkedInvoice := range invoice.GetLinkedInvoices() {
			orgID := linkedInvoice.GetOrgId()
			if orgID == "" {
				continue
			}
			crossOrgBilling[orgID] = append(crossOrgBilling[orgID], linkedInvoice)
		}
	}

	return crossOrgBilling, nil
}
