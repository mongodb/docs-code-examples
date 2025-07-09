package main

import (
	"context"
	"fmt"
	"log"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/billing"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/errors"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not loaded: %v", err)
	}

	secrets, cfg, err := config.LoadAll("configs/config.json")
	if err != nil {
		errors.ExitWithError("Failed to load configuration", err)
	}

	client, err := auth.NewClient(cfg, secrets)
	if err != nil {
		errors.ExitWithError("Failed to initialize authentication client", err)
	}

	ctx := context.Background()
	p := &admin.ListInvoicesApiParams{
		OrgId: cfg.OrgID,
	}

	fmt.Printf("Fetching linked organizations for billing organization: %s\n", p.OrgId)

	invoices, err := billing.GetCrossOrgBilling(ctx, client.InvoicesApi, p)
	if err != nil {
		errors.ExitWithError(fmt.Sprintf("Failed to retrieve cross-organization billing data for %s", p.OrgId), err)
	}

	displayLinkedOrganizations(invoices, p.OrgId)
}

func displayLinkedOrganizations(invoices map[string][]admin.BillingInvoiceMetadata, primaryOrgID string) {
	var linkedOrgs []string
	for orgID := range invoices {
		if orgID != primaryOrgID {
			linkedOrgs = append(linkedOrgs, orgID)
		}
	}

	if len(linkedOrgs) == 0 {
		fmt.Println("No linked organizations found for the billing organization")
		return
	}

	fmt.Printf("Found %d linked organizations:\n", len(linkedOrgs))
	for i, orgID := range linkedOrgs {
		fmt.Printf("  %d. Organization ID: %s\n", i+1, orgID)
	}
}

