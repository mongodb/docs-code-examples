// :snippet-start: historical-billing
// :state-remove-start: copy
// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
// :state-remove-end: [copy]
package main

import (
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/data/export"
	"atlas-sdk-go/internal/fileutils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"

	"atlas-sdk-go/internal/billing"
)

// :remove-start:
// TODO: QUESTION FOR REVIEWER: currently set to pull the past 6 months from current; do we want any additional configurations in this example?
// :remove-end:

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not loaded: %v", err)
	}

	secrets, cfg, err := config.LoadAll("configs/config.json")
	if err != nil {
		log.Fatalf("config: failed to load file: %v", err)
	}

	client, err := auth.NewClient(cfg, secrets)
	if err != nil {
		log.Fatalf("auth: failed client init: %v", err)
	}

	ctx := context.Background()

	fmt.Printf("Fetching historical invoices for organization: %s\n", cfg.OrgID)

	invoices, err := billing.ListInvoicesForOrg(ctx, client.InvoicesApi, cfg.OrgID,
		billing.WithViewLinkedInvoices(true),
		billing.WithIncludeCount(true),
		billing.WithDateRange(time.Now().AddDate(0, -6, 0), time.Now())) // previous six months
	if err != nil {
		log.Fatalf("billing: cannot retrieve invoices: %v", err)
	}

	if invoices.GetTotalCount() > 0 {
		fmt.Printf("Total count of invoices: %d\n", invoices.GetTotalCount())
	} else {
		fmt.Println("No invoices found for the specified date range.")
		return
	}

	// Export invoice data to JSON and CSV file formats
	outDir := "invoices"
	prefix := fmt.Sprintf("historical_%s", cfg.OrgID)

	jsonPath, err := fileutils.GenerateOutputPath(outDir, prefix, "json")
	if err != nil {
		log.Fatalf("common: generate output path: %v", err)
	}
	if err := export.ToJSON(invoices.GetResults(), jsonPath); err != nil {
		log.Fatalf("json: write file: %v", err)
	}
	fmt.Printf("Exported invoice data to %s\n", jsonPath)

	csvPath, err := fileutils.GenerateOutputPath(outDir, prefix, "csv")
	if err != nil {
		log.Fatalf("common: generate output path: %v", err)
	}

	headers := []string{"InvoiceID", "Status", "Created", "AmountBilled"}

	err = export.ToCSVWithMapper(invoices.GetResults(), csvPath, headers, func(invoice billing.InvoiceOption) []string {
		return []string{
			invoice.GetId(),
			invoice.GetStatusName(),
			invoice.GetCreated().Format(time.RFC3339),
			fmt.Sprintf("%.2f", float64(invoice.GetAmountBilledCents())/100.0),
		}
	})
	if err != nil {
		log.Fatalf("export: failed to write CSV file: %v", err)
	}
	fmt.Printf("Exported invoice data to %s\n", csvPath)
}

// :snippet-end: [historical-billing]
