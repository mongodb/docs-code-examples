// :snippet-start: line-items
// :state-remove-start: copy
// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
// :state-remove-end: [copy]
package main

import (
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/billing"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/data/export"
	"atlas-sdk-go/internal/errors"
	"atlas-sdk-go/internal/fileutils"
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not loaded: %v", err)
	}

	secrets, cfg, err := config.LoadAll("configs/ignore.config.json")
	if err != nil {
		errors.ExitWithError("Failed to load configuration", err)
	}

	client, err := auth.NewClient(cfg, secrets)
	if err != nil {
		errors.ExitWithError("Failed to initialize authentication client", err)
	}

	ctx := context.Background()
	OrgID := cfg.OrgID

	fmt.Printf("Fetching pending invoices for organization: %s\n", OrgID)

	details, err := billing.CollectLineItemBillingData(ctx, client.InvoicesApi, client.OrganizationsApi, OrgID, nil)
	if err != nil {
		errors.ExitWithError("Failed to fetch billing data", err)
	}

	fmt.Printf("Found %d line items in pending invoices\n", len(details))

	outDir := "invoices"
	prefix := fmt.Sprintf("pending_%s", OrgID)

	// Export to JSON
	jsonPath, err := fileutils.GenerateOutputPath(outDir, prefix, "json")
	if err != nil {
		errors.ExitWithError("Failed to generate JSON output path", err)
	}

	if err := export.ToJSON(details, jsonPath); err != nil {
		errors.ExitWithError("Failed to write JSON file", err)
	}
	fmt.Printf("Exported billing data to %s\n", jsonPath)

	// Export to CSV file
	csvPath, err := fileutils.GenerateOutputPath(outDir, prefix, "csv")
	if err != nil {
		errors.ExitWithError("Failed to generate CSV output path", err)
	}

	headers := []string{"Organization", "OrgID", "Project", "ProjectID", "Cluster",
		"SKU", "Cost", "Date", "Provider", "Instance", "Category"}

	err = export.ToCSVWithMapper(details, csvPath, headers, func(item billing.Detail) []string {
		return []string{
			item.Org.Name,
			item.Org.ID,
			item.Project.Name,
			item.Project.ID,
			item.Cluster,
			item.SKU,
			fmt.Sprintf("%.2f", item.Cost),
			item.Date.Format("2006-01-02"),
			item.Provider,
			item.Instance,
			item.Category,
		}
	})
	if err != nil {
		errors.ExitWithError("Failed to write CSV file", err)
	}
	fmt.Printf("Exported billing data to %s\n", csvPath)
}

// :snippet-end: [line-items]
