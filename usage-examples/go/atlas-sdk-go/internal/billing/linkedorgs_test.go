package billing_test

import (
	"atlas-sdk-go/internal/billing"
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"go.mongodb.org/atlas-sdk/v20250219001/mockadmin"
)

func TestGetLinkedOrgs_Success(t *testing.T) {
	billingOrgID := "billingOrgABC"
	linkedOrgID1 := "linkedOrgABC"
	linkedOrgID2 := "linkedOrgDEF"
	invoiceID1 := "inv_123"
	invoiceID2 := "inv_456"

	// Create mock response with linked invoices
	mockResponse := &admin.PaginatedApiInvoiceMetadata{
		Results: &[]admin.BillingInvoiceMetadata{
			{
				Id: &invoiceID1,
				LinkedInvoices: &[]admin.BillingInvoiceMetadata{
					{OrgId: &linkedOrgID1, AmountBilledCents: admin.PtrInt64(1000)},
				},
			},
			{
				Id: &invoiceID2,
				LinkedInvoices: &[]admin.BillingInvoiceMetadata{
					{OrgId: &linkedOrgID1, AmountBilledCents: admin.PtrInt64(2000)},
					{OrgId: &linkedOrgID2, AmountBilledCents: admin.PtrInt64(500)},
				},
			},
		},
	}

	mockSvc := mockadmin.NewInvoicesApi(t)
	mockSvc.EXPECT().
		ListInvoices(mock.Anything, billingOrgID).
		Return(admin.ListInvoicesApiRequest{ApiService: mockSvc}).Once()
	mockSvc.EXPECT().
		ListInvoicesExecute(mock.Anything).
		Return(mockResponse, nil, nil).Once()

	params := &admin.ListInvoicesApiParams{OrgId: billingOrgID}
	linkedOrgs, err := billing.GetLinkedOrgs(context.Background(), mockSvc, params)

	require.NoError(t, err)
	assert.Len(t, linkedOrgs, 2, "Should return two linked organizations")
	sort.Strings(linkedOrgs)
	assert.Contains(t, linkedOrgs, linkedOrgID1, "Should contain linkedOrgID1")
	assert.Contains(t, linkedOrgs, linkedOrgID2, "Should contain linkedOrgID2")
}

func TestGetLinkedOrgs_ApiError(t *testing.T) {
	billingOrgID := "billingOrgErr"
	expectedError := errors.New("API error")

	mockSvc := mockadmin.NewInvoicesApi(t)
	mockSvc.EXPECT().
		ListInvoices(mock.Anything, billingOrgID).
		Return(admin.ListInvoicesApiRequest{ApiService: mockSvc}).Once()
	mockSvc.EXPECT().
		ListInvoicesExecute(mock.Anything).
		Return(nil, nil, expectedError).Once()

	params := &admin.ListInvoicesApiParams{OrgId: billingOrgID}
	_, err := billing.GetLinkedOrgs(context.Background(), mockSvc, params)

	// Verify error handling
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get cross-org billing")
	assert.ErrorContains(t, err, expectedError.Error(), "Should return API error")
}

func TestGetLinkedOrgs_NoLinkedOrgs(t *testing.T) {
	billingOrgID := "billingOrg123"
	invoiceID := "no_links"

	// Create mock response with no linked invoices
	mockResponse := &admin.PaginatedApiInvoiceMetadata{
		Results: &[]admin.BillingInvoiceMetadata{
			{
				Id:             &invoiceID,
				LinkedInvoices: &[]admin.BillingInvoiceMetadata{},
			},
		},
	}

	mockSvc := mockadmin.NewInvoicesApi(t)
	mockSvc.EXPECT().
		ListInvoices(mock.Anything, billingOrgID).
		Return(admin.ListInvoicesApiRequest{ApiService: mockSvc}).Once()
	mockSvc.EXPECT().
		ListInvoicesExecute(mock.Anything).
		Return(mockResponse, nil, nil).Once()

	params := &admin.ListInvoicesApiParams{OrgId: billingOrgID}
	linkedOrgs, err := billing.GetLinkedOrgs(context.Background(), mockSvc, params)

	require.NoError(t, err)
	assert.Empty(t, linkedOrgs, "Should return empty when no linked orgs exist")
}

func TestGetLinkedOrgs_MissingOrgID(t *testing.T) {
	billingOrgID := "billingOrg123"
	linkedOrgID := "validOrgID"
	invoiceID := "inv_missing_org"

	// Create mock response with one valid and one invalid linked invoice
	mockResponse := &admin.PaginatedApiInvoiceMetadata{
		Results: &[]admin.BillingInvoiceMetadata{
			{
				Id: &invoiceID,
				LinkedInvoices: &[]admin.BillingInvoiceMetadata{
					{OrgId: nil, AmountBilledCents: admin.PtrInt64(500)},
					{OrgId: &linkedOrgID, AmountBilledCents: admin.PtrInt64(1000)},
				},
			},
		},
	}

	mockSvc := mockadmin.NewInvoicesApi(t)
	mockSvc.EXPECT().
		ListInvoices(mock.Anything, billingOrgID).
		Return(admin.ListInvoicesApiRequest{ApiService: mockSvc}).Once()
	mockSvc.EXPECT().
		ListInvoicesExecute(mock.Anything).
		Return(mockResponse, nil, nil).Once()

	// Run test
	params := &admin.ListInvoicesApiParams{OrgId: billingOrgID}
	linkedOrgs, err := billing.GetLinkedOrgs(context.Background(), mockSvc, params)

	require.NoError(t, err)
	assert.Len(t, linkedOrgs, 1, "Should return one linked organization")
	assert.Equal(t, linkedOrgID, linkedOrgs[0], "Should return the valid linked organization ID")
}
