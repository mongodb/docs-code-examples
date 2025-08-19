package main

import (
	"context"
	"fmt"
	"log"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/billing"
	"atlas-sdk-go/internal/config"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

func main() {
	_ = godotenv.Load() // or godotenv.Load(".env.development")

	ctx := context.Background()
	envName := config.Environment("test")    // Cast string to config.Environment
	configPath := "configs/config.test.json" // Optional explicit config file path; if empty, uses environment-based path
	secrets, cfg, err := config.LoadAll(envName, configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration %v", err)
	}

	client, err := auth.NewClient(ctx, cfg, secrets) // Pass pointers
	if err != nil {
		log.Fatalf("Failed to initialize authentication client: %v", err)
	}

	p := &admin.ListInvoicesApiParams{
		OrgId: cfg.OrgID,
	}

	fmt.Printf("Fetching linked organizations for billing organization: %s\n", p.OrgId)

	invoices, err := billing.GetCrossOrgBilling(ctx, client.InvoicesApi, p)
	if err != nil {
		log.Fatalf("Failed to retrieve cross-organization billing data for %s: %v", p.OrgId, err)
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

