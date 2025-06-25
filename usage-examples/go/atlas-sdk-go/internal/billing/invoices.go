package billing

import (
	"atlas-sdk-go/internal"
	"context"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"time"
)

// InvoiceOption defines a function type that modifies the parameters for listing invoices.
type InvoiceOption func(*admin.ListInvoicesApiParams)

// WithIncludeCount sets the optional includeCount parameter (default: true).
func WithIncludeCount(includeCount bool) InvoiceOption {
	return func(p *admin.ListInvoicesApiParams) {
		p.IncludeCount = &includeCount
	}
}

// WithItemsPerPage sets the optional itemsPerPage parameter (default: 100).
func WithItemsPerPage(itemsPerPage int) InvoiceOption {
	return func(p *admin.ListInvoicesApiParams) {
		p.ItemsPerPage = &itemsPerPage
	}
}

// WithPageNum sets the optional pageNum parameter (default: 1).
func WithPageNum(pageNum int) InvoiceOption {
	return func(p *admin.ListInvoicesApiParams) {
		p.PageNum = &pageNum
	}
}

// WithViewLinkedInvoices sets the optional viewLinkedInvoices parameter (default: true).
func WithViewLinkedInvoices(viewLinkedInvoices bool) InvoiceOption {
	return func(p *admin.ListInvoicesApiParams) {
		p.ViewLinkedInvoices = &viewLinkedInvoices
	}
}

// WithStatusNames sets the optional statusNames parameter (default: all statuses).
func WithStatusNames(statusNames []string) InvoiceOption {
	return func(p *admin.ListInvoicesApiParams) {
		p.StatusNames = &statusNames
	}
}

// WithDateRange sets the optional fromDate and toDate parameters (default: all possible dates).
func WithDateRange(fromDate, toDate time.Time) InvoiceOption {
	return func(p *admin.ListInvoicesApiParams) {
		from := fromDate.Format("2006-01-02")
		to := toDate.Format("2006-01-02")
		p.FromDate = &from
		p.ToDate = &to
	}
}

// WithSortBy sets the optional sortBy parameter (default: "END_DATE").
func WithSortBy(sortBy string) InvoiceOption {
	return func(p *admin.ListInvoicesApiParams) {
		p.SortBy = &sortBy
	}
}

// WithOrderBy sets the optional orderBy parameter (default: "desc").
func WithOrderBy(orderBy string) InvoiceOption {
	return func(p *admin.ListInvoicesApiParams) {
		p.OrderBy = &orderBy
	}
}

// ListInvoicesForOrg returns all eligible invoices for the given organization, including linked organizations when cross-organization billing is enabled. If optional parameters aren't specified, default values are used.
// NOTE: Organization Billing Admin or Organization Owner role required to view linked invoices.
func ListInvoicesForOrg(ctx context.Context, sdk admin.InvoicesApi, orgId string, opts ...InvoiceOption) (*admin.PaginatedApiInvoiceMetadata, error) {
	params := &admin.ListInvoicesApiParams{
		OrgId: orgId,
	}

	for _, opt := range opts {
		opt(params)
	}

	req := sdk.ListInvoices(ctx, params.OrgId)

	if params.IncludeCount != nil {
		req = req.IncludeCount(*params.IncludeCount)
	}
	if params.ItemsPerPage != nil {
		req = req.ItemsPerPage(*params.ItemsPerPage)
	}
	if params.PageNum != nil {
		req = req.PageNum(*params.PageNum)
	}
	if params.ViewLinkedInvoices != nil {
		req = req.ViewLinkedInvoices(*params.ViewLinkedInvoices)
	}
	if params.StatusNames != nil {
		req = req.StatusNames(*params.StatusNames)
	}
	if params.FromDate != nil {
		req = req.FromDate(*params.FromDate)
	}
	if params.ToDate != nil {
		req = req.ToDate(*params.ToDate)
	}
	if params.SortBy != nil {
		req = req.SortBy(*params.SortBy)
	}
	if params.OrderBy != nil {
		req = req.OrderBy(*params.OrderBy)
	}

	r, _, err := req.Execute()

	if err != nil {
		return nil, internal.FormatAPIError("list invoices", orgId, err)
	}

	return r, nil
}
