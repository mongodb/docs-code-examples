// :snippet-start: historical-billing
// :state-remove-start: copy
// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
// :state-remove-end: [copy]
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/billing"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/data/export"
	"atlas-sdk-go/internal/errors"
	"atlas-sdk-go/internal/fileutils"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

func main() {
	configPath := ""  // Use default config path for environment
	explicitEnv := "" // Use default environment
	secrets, cfg, err := config.LoadAll(configPath, explicitEnv)
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

	fmt.Printf("Fetching historical invoices for organization: %s\n", p.OrgId)

	// Fetch invoices from the previous six months with the provided options
	invoices, err := billing.ListInvoicesForOrg(ctx, client.InvoicesApi, p,
		billing.WithViewLinkedInvoices(true),
		billing.WithIncludeCount(true),
		billing.WithDateRange(time.Now().AddDate(0, -6, 0), time.Now()))
	if err != nil {
		errors.ExitWithError("Failed to retrieve invoices", err)
	}

	if invoices.GetTotalCount() > 0 {
		fmt.Printf("Total count of invoices: %d\n", invoices.GetTotalCount())
	} else {
		fmt.Println("No invoices found for the specified date range")
		return
	}

	// Export invoice data to be used in other systems or for reporting
	outDir := "invoices"
	prefix := fmt.Sprintf("historical_%s", p.OrgId)

	exportInvoicesToJSON(invoices, outDir, prefix)
	exportInvoicesToCSV(invoices, outDir, prefix)
	// :remove-start:
	// Clean up (internal-only function)
	if err = fileutils.SafeDelete(outDir); err != nil {
		log.Printf("Cleanup error: %v", err)
	}
	fmt.Println("Deleted generated files from", outDir)
	// :remove-end:
}

func exportInvoicesToJSON(invoices *admin.PaginatedApiInvoiceMetadata, outDir, prefix string) {
	jsonPath, err := fileutils.GenerateOutputPath(outDir, prefix, "json")
	if err != nil {
		errors.ExitWithError("Failed to generate JSON output path", err)
	}
	if err := export.ToJSON(invoices.GetResults(), jsonPath); err != nil {
		errors.ExitWithError("Failed to write JSON file", err)
	}
	fmt.Printf("Exported invoice data to %s\n", jsonPath)
}

func exportInvoicesToCSV(invoices *admin.PaginatedApiInvoiceMetadata, outDir, prefix string) {
	csvPath, err := fileutils.GenerateOutputPath(outDir, prefix, "csv")
	if err != nil {
		errors.ExitWithError("Failed to generate CSV output path", err)
	}

	// Set the headers and mapped rows for the CSV export
	headers := []string{"InvoiceID", "Status", "Created", "AmountBilled"}
	err = export.ToCSVWithMapper(invoices.GetResults(), csvPath, headers, func(invoice admin.BillingInvoiceMetadata) []string {
		return []string{
			invoice.GetId(),
			invoice.GetStatusName(),
			invoice.GetCreated().Format(time.RFC3339),
			fmt.Sprintf("%.2f", float64(invoice.GetAmountBilledCents())/100.0),
		}
	})
	if err != nil {
		errors.ExitWithError("Failed to write CSV file", err)
	}

	fmt.Printf("Exported invoice data to %s\n", csvPath)
}

// :snippet-end: [historical-billing]
// :state-remove-start: copy
// NOTE: INTERNAL
// ** OUTPUT EXAMPLE **
//
// Fetching historical invoices for organization: 5f7a9aec7d78fc03b42959328
// Total count of invoices: 12
// Exported invoice data to invoices/historical_5f7a9aec7d78fc03b42959328.json
// Exported invoice data to invoices/historical_5f7a9aec7d78fc03b42959328.csv
// :state-remove-end: [copy]
