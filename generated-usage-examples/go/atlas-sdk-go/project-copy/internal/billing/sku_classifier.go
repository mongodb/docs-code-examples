package billing

import (
	"strings"
)

// determineProvider identifies the cloud provider based on SKU
func determineProvider(sku string) string {
	if strings.Contains(sku, "AWS") {
		return "AWS"
	} else if strings.Contains(sku, "AZURE") {
		return "AZURE"
	} else if strings.Contains(sku, "GCP") {
		return "GCP"
	}
	return "n/a"
}

// determineInstance extracts the instance type from SKU
func determineInstance(sku string) string {
	parts := strings.Split(sku, "_INSTANCE_")
	if len(parts) > 1 {
		return parts[1]
	}
	return "non-instance"
}

// determineCategory categorizes the SKU
func determineCategory(sku string) string {
	categoryPatterns := map[string]string{
		"_INSTANCE":        "instances",
		"BACKUP":           "backup",
		"PIT_RESTORE":      "backup",
		"DATA_TRANSFER":    "data xfer",
		"STORAGE":          "storage",
		"BI_CONNECTOR":     "bi-connector",
		"DATA_LAKE":        "data lake",
		"AUDITING":         "audit",
		"ATLAS_SUPPORT":    "support",
		"FREE_SUPPORT":     "free support",
		"CHARTS":           "charts",
		"SERVERLESS":       "serverless",
		"SECURITY":         "security",
		"PRIVATE_ENDPOINT": "private endpoint",
	}

	for pattern, category := range categoryPatterns {
		if strings.Contains(sku, pattern) {
			return category
		}
	}
	return "other"
}
