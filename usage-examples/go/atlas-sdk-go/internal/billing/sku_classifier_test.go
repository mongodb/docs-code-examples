package billing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineProvider(t *testing.T) {
	tests := []struct {
		name     string
		sku      string
		expected string
	}{
		{"AWS SKU", "MONGODB_ATLAS_AWS_INSTANCE_M10", "AWS"},
		{"AZURE SKU", "MONGODB_ATLAS_AZURE_INSTANCE_M20", "AZURE"},
		{"GCP SKU", "MONGODB_ATLAS_GCP_INSTANCE_M30", "GCP"},
		{"Unknown provider", "MONGODB_ATLAS_INSTANCE_M40", "n/a"},
		{"Empty SKU", "", "n/a"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := determineProvider(tc.sku)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDetermineInstance(t *testing.T) {
	tests := []struct {
		name     string
		sku      string
		expected string
	}{
		{"Basic instance", "MONGODB_ATLAS_AWS_INSTANCE_M10", "M10"},
		{"Complex instance name", "MONGODB_ATLAS_AWS_INSTANCE_M30_NVME", "M30_NVME"},
		{"No instance marker", "MONGODB_ATLAS_BACKUP", "non-instance"},
		{"Empty SKU", "", "non-instance"},
		{"Multiple instance markers", "INSTANCE_M10_INSTANCE_M20", "M20"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := determineInstance(tc.sku)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDetermineCategory(t *testing.T) {
	tests := []struct {
		name     string
		sku      string
		expected string
	}{
		{"Instance category", "MONGODB_ATLAS_AWS_INSTANCE_M10", "instances"},
		{"Backup category", "MONGODB_ATLAS_BACKUP", "backup"},
		{"PIT Restore", "MONGODB_ATLAS_PIT_RESTORE", "backup"},
		{"Data Transfer", "MONGODB_ATLAS_DATA_TRANSFER", "data xfer"},
		{"Storage", "MONGODB_ATLAS_STORAGE", "storage"},
		{"BI Connector", "MONGODB_ATLAS_BI_CONNECTOR", "bi-connector"},
		{"Data Lake", "MONGODB_ATLAS_DATA_LAKE", "data lake"},
		{"Auditing", "MONGODB_ATLAS_AUDITING", "audit"},
		{"Atlas Support", "MONGODB_ATLAS_SUPPORT", "support"},
		{"Free Support", "MONGODB_ATLAS_FREE_SUPPORT", "free support"},
		{"Charts", "MONGODB_ATLAS_CHARTS", "charts"},
		{"Serverless", "MONGODB_ATLAS_SERVERLESS", "serverless"},
		{"Security", "MONGODB_ATLAS_SECURITY", "security"},
		{"Private Endpoint", "MONGODB_ATLAS_PRIVATE_ENDPOINT", "private endpoint"},
		{"Other category", "MONGODB_ATLAS_UNKNOWN", "other"},
		{"Empty SKU", "", "other"},
		{"Multiple patterns", "MONGODB_ATLAS_BACKUP_STORAGE", "backup"}, // First match should win
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := determineCategory(tc.sku)
			assert.Equal(t, tc.expected, result)
		})
	}
}
