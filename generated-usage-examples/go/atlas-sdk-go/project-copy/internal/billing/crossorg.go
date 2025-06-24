package billing

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"atlas-sdk-go/internal"
)

// GetCrossOrgBilling returns all invoices for the billing organization and any linked organizations.
// NOTE: Organization Billing Admin or Organization Owner role required to view linked invoices.
func GetCrossOrgBilling(ctx context.Context, sdk admin.InvoicesApi, p *admin.ListInvoicesApiParams) (map[string][]admin.BillingInvoiceMetadata, error) {
	req := sdk.ListInvoices(ctx, p.OrgId)

	r, _, err := req.Execute()

	if err != nil {
		return nil, internal.FormatAPIError("list invoices", p.OrgId, err)
	}

	crossOrgBilling := make(map[string][]admin.BillingInvoiceMetadata)
	if r == nil || !r.HasResults() || len(r.GetResults()) == 0 {
		return crossOrgBilling, nil
	}

	crossOrgBilling[p.OrgId] = r.GetResults()
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
