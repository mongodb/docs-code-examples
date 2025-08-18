package archive

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"atlas-sdk-go/internal/errors"
)

// Candidate represents a collection eligible for archiving
type Candidate struct {
	DatabaseName    string
	CollectionName  string
	DateField       string
	DateFormat      string
	RetentionDays   int
	PartitionFields []string
}

// Options defines configuration settings for archive operations
type Options struct {
	// Default data retention period multiplier
	DefaultRetentionMultiplier int
	// Minimum retention days required before archiving
	MinimumRetentionDays int
	// Whether to enable data expiration
	EnableDataExpiration bool
	// Schedule for archive operations
	ArchiveSchedule string
}

// DefaultOptions provides sensible defaults for archiving
func DefaultOptions() Options {
	return Options{
		DefaultRetentionMultiplier: 2,
		MinimumRetentionDays:       30,
		EnableDataExpiration:       true,
		ArchiveSchedule:            "DAILY",
	}
}

type ExpireAfterDays struct {
	// NOTE: this placeholder struct can be extended to include more complex rules if needed
	ExpireAfterDays int `json:"expireAfterDays,omitempty"`
}

// CollectionsForArchiving identifies collections suitable for archiving as a simplified example for demonstration purposes.
// This function returns a list of Candidates that meet the archiving criteria
// NOTE: In a real implementation, you would determine which collections are eligible based on criteria analysis such as size, age, and access patterns.
func CollectionsForArchiving(ctx context.Context, sdk *admin.APIClient,
	projectID, clusterName string) []Candidate {
	// For demonstration purposes, we specify example Candidates
	return []Candidate{
		{
			DatabaseName:    "sample_analytics",
			CollectionName:  "transactions",
			DateField:       "transaction_date",
			DateFormat:      "DATE",
			RetentionDays:   90,
			PartitionFields: []string{"customer_id", "merchant"},
		},
		{
			DatabaseName:    "sample_logs",
			CollectionName:  "application_logs",
			DateField:       "timestamp",
			DateFormat:      "EPOCH_MILLIS",
			RetentionDays:   30,
			PartitionFields: []string{"service_name", "log_level"},
		},
	}
}

// ValidateCandidate ensures the archiving candidate meets requirements
func ValidateCandidate(candidate Candidate, opts Options) error {
	if candidate.DatabaseName == "" || candidate.CollectionName == "" {
		return fmt.Errorf("database name and collection name are required")
	}

	if candidate.RetentionDays < opts.MinimumRetentionDays {
		return fmt.Errorf("retention days must be at least %d", opts.MinimumRetentionDays)
	}

	if len(candidate.PartitionFields) == 0 {
		return fmt.Errorf("at least one partition field is required")
	}

	// For date-based archiving, validate date field settings
	if candidate.DateField != "" {
		validFormats := map[string]bool{
			"DATE":              true,
			"EPOCH_SECONDS":     true,
			"EPOCH_MILLIS":      true,
			"EPOCH_NANOSECONDS": true,
			"OBJECTID":          true,
		}
		if !validFormats[candidate.DateFormat] {
			return fmt.Errorf("invalid date format: %s", candidate.DateFormat)
		}

		// Check if date field is included in partition fields
		dateFieldIncluded := false
		for _, field := range candidate.PartitionFields {
			if field == candidate.DateField {
				dateFieldIncluded = true
				break
			}
		}
		if !dateFieldIncluded {
			return fmt.Errorf("date field %s must be included in partition fields", candidate.DateField)
		}
	}

	return nil
}

// ConfigureOnlineArchive configures online archive for a collection
func ConfigureOnlineArchive(ctx context.Context, sdk *admin.APIClient,
	projectID, clusterName string, candidate Candidate) error {

	opts := DefaultOptions()

	if err := ValidateCandidate(candidate, opts); err != nil {
		return errors.FormatError("validate archive candidate",
			fmt.Sprintf("%s.%s", candidate.DatabaseName, candidate.CollectionName),
			err)
	}

	// Create partition fields configuration
	var partitionFields []admin.PartitionField
	for idx, field := range candidate.PartitionFields {
		partitionFields = append(partitionFields, admin.PartitionField{
			FieldName: field,
			Order:     idx + 1,
		})
	}

	// Setup data expiration if enabled
	var dataExpiration *admin.OnlineArchiveSchedule
	if opts.EnableDataExpiration && opts.DefaultRetentionMultiplier > 0 {
		expirationDays := candidate.RetentionDays * opts.DefaultRetentionMultiplier
		dataExpiration = &admin.OnlineArchiveSchedule{
			Type: opts.ArchiveSchedule,
		}

		// Define request body
		archiveReq := &admin.BackupOnlineArchiveCreate{
			CollName:        candidate.CollectionName,
			DbName:          candidate.DatabaseName,
			PartitionFields: &partitionFields,
		}

		// Set expiration if configured
		if dataExpiration != nil {
			archiveReq.DataExpirationRule = &admin.DataExpirationRule{
				ExpireAfterDays: admin.PtrInt(expirationDays),
			}
		}

		// Configure date criteria if present
		if candidate.DateField != "" {
			archiveReq.Criteria = admin.Criteria{
				DateField:       admin.PtrString(candidate.DateField),
				DateFormat:      admin.PtrString(candidate.DateFormat),
				ExpireAfterDays: admin.PtrInt(candidate.RetentionDays),
			}
		}

		// Execute the request
		_, _, err := sdk.OnlineArchiveApi.CreateOnlineArchive(ctx, projectID, clusterName, archiveReq).Execute()

		if err != nil {
			return errors.FormatError("create online archive",
				fmt.Sprintf("%s.%s", candidate.DatabaseName, candidate.CollectionName),
				err)
		}
	}

	return nil
}
