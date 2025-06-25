package main

import (
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/billing"
	"atlas-sdk-go/internal/config"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"time"
)

// view invoices in the past six months for a given organization, including linked invoices

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

	fmt.Printf("Fetching historical invoices for organization: %s\n", cfg.OrgID)

	// TODO: confirm if we want to filter statuses or date range in this example
	nonPendingStatuses := []string{
		// "PENDING",
		"CLOSED", "FORGIVEN", "FAILED", "PAID", "FREE", "PREPAID", "INVOICED"}
	invoices, err := billing.ListInvoicesForOrg(ctx, sdk.InvoicesApi, cfg.OrgID,
		billing.WithStatusNames(nonPendingStatuses),
		billing.WithViewLinkedInvoices(true),
		billing.WithIncludeCount(true),
		billing.WithDateRange(time.Now().AddDate(0, -6, 0), time.Now())) // past six months

	if err != nil {
		log.Fatalf("Failed to retrieve invoices: %v", err)
	}

	if invoices == nil || !invoices.HasResults() || len(invoices.GetResults()) == 0 {
		fmt.Println("No invoices found")
		return
	}

	fmt.Printf("Found %d invoices\n", len(invoices.GetResults()))
	for i, invoice := range invoices.GetResults() {
		fmt.Printf("  %d. Invoice #%s - Status: %s - Created: %s - Amount: $%.d\n",
			i+1,
			invoice.GetId(),
			invoice.GetStatusName(),
			invoice.GetCreated(),
			invoice.GetAmountBilledCents()/100.0)
	}
}
