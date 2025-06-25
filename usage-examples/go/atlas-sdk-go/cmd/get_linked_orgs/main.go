// :snippet-start: get-linked-orgs
// :state-remove-start: copy
// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
// :state-remove-end: [copy]
package main

import (
	"context"
	"fmt"
	"log"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/billing"
	"atlas-sdk-go/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	secrets, cfg, err := config.LoadAll("configs/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	sdk, err := auth.NewClient(cfg, secrets)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	ctx := context.Background()

	fmt.Printf("Fetching cross-org billing for organization: %s\n", cfg.OrgID)
	results, err := billing.GetCrossOrgBilling(ctx, sdk.InvoicesApi, cfg.OrgID)
	if err != nil {
		log.Fatalf("Failed to retrieve invoices: %v", err)
	}
	if len(results) == 0 {
		fmt.Println("No invoices found for the billing organization")
		return
	}

	// Print the returned map of invoices grouped by organization ID
	fmt.Printf("Found %d organizations with invoices:\n", len(results))
	for orgID, invoices := range results {
		fmt.Printf("  Organization ID: %s\n", orgID)
		if len(invoices) == 0 {
			fmt.Println("    No invoices found for this organization")
			continue
		}
		for i, invoice := range invoices {
			fmt.Printf("    %d. Invoice #%s - Status: %s - Created: %s - Amount: $%.2f\n",
				i+1,
				invoice.GetId(),
				invoice.GetStatus(),
				invoice.GetCreatedDate(),
				invoice.GetAmountBilledCents()/100.0)
		}
	}
	// linkedOrgs, err := billing.ListLinkedOrgs(ctx, sdk.InvoicesApi, cfg.OrgID)
	// if err != nil {
	// 	log.Fatalf("Failed to retrieve linked organizations: %v", err)
	// }
	// if len(linkedOrgs) == 0 {
	// 	fmt.Println("No linked organizations found for the billing org")
	// 	return
	// }
	// fmt.Println("Linked organizations:")
	// for i, org := range linkedOrgs {
	// 	fmt.Printf("  %d. %v\n", i+1, org)
	// }
}

// :snippet-end: [get-linked-orgs]
