// :snippet-start: line-items
// :state-remove-start: copy
// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
// :state-remove-end: [copy]
package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/billing"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/data/export"
	"atlas-sdk-go/internal/errors"
	"atlas-sdk-go/internal/fileutils"

	"github.com/joho/godotenv"
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

	fmt.Printf("Fetching pending invoices for organization: %s\n", p.OrgId)

	details, err := billing.CollectLineItemBillingData(ctx, client.InvoicesApi, client.OrganizationsApi, p.OrgId, nil)
	if err != nil {
		errors.ExitWithError(fmt.Sprintf("Failed to retrieve pending invoices for %s", p.OrgId), err)
	}

	fmt.Printf("Found %d line items in pending invoices\n", len(details))

	// Export invoice data to be used in other systems or for reporting
	outDir := "invoices"
	prefix := fmt.Sprintf("pending_%s", p.OrgId)

	exportInvoicesToJSON(details, outDir, prefix)
	exportInvoicesToCSV(details, outDir, prefix)
	// :remove-start:
	// Clean up (internal-only function)
	if err := fileutils.SafeDelete(outDir); err != nil {
		log.Printf("Cleanup error: %v", err)
	}
	fmt.Println("Deleted generated files from", outDir)
	// :remove-end:
}

func exportInvoicesToJSON(details []billing.Detail, outDir, prefix string) {
	jsonPath, err := fileutils.GenerateOutputPath(outDir, prefix, "json")
	if err != nil {
		errors.ExitWithError("Failed to generate JSON output path", err)
	}

	if err := export.ToJSON(details, jsonPath); err != nil {
		errors.ExitWithError("Failed to write JSON file", err)
	}
	fmt.Printf("Exported billing data to %s\n", jsonPath)
}

func exportInvoicesToCSV(details []billing.Detail, outDir, prefix string) {
	csvPath, err := fileutils.GenerateOutputPath(outDir, prefix, "csv")
	if err != nil {
		errors.ExitWithError("Failed to generate CSV output path", err)
	}

	// Set the headers and mapped rows for the CSV export
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
// :state-remove-start: copy
// NOTE: INTERNAL
// ** OUTPUT EXAMPLE **
//
// Fetching pending invoices for organization: 5f7a9ec7d78fc03b42959328
//
// Found 3 line items in pending invoices
// Exported billing data to invoices/pending_5f7a9ec7d78fc03b42959328.json
// Exported billing data to invoices/pending_5f7a9ec7d78fc03b42959328.csv
// :state-remove-end: [copy]
