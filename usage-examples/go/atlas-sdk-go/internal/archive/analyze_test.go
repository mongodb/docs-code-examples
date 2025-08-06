package archive

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

func TestDefaultOptions_ReturnsExpectedDefaults(t *testing.T) {
	t.Parallel()

	opts := DefaultOptions()

	assert.Equal(t, 2, opts.DefaultRetentionMultiplier)
	assert.Equal(t, 30, opts.MinimumRetentionDays)
	assert.True(t, opts.EnableDataExpiration)
	assert.Equal(t, "DAILY", opts.ArchiveSchedule)
}

func TestCollectionsForArchiving_ReturnsExpectedCandidates(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	var client *admin.APIClient

	candidates := CollectionsForArchiving(ctx, client, "project123", "cluster456")

	require.Len(t, candidates, 2)

	analyticsCandidate := candidates[0]
	assert.Equal(t, "sample_analytics", analyticsCandidate.DatabaseName)
	assert.Equal(t, "transactions", analyticsCandidate.CollectionName)
	assert.Equal(t, "transaction_date", analyticsCandidate.DateField)
	assert.Equal(t, "DATE", analyticsCandidate.DateFormat)
	assert.Equal(t, 90, analyticsCandidate.RetentionDays)
	assert.Equal(t, []string{"customer_id", "merchant"}, analyticsCandidate.PartitionFields)

	logsCandidate := candidates[1]
	assert.Equal(t, "sample_logs", logsCandidate.DatabaseName)
	assert.Equal(t, "application_logs", logsCandidate.CollectionName)
	assert.Equal(t, "timestamp", logsCandidate.DateField)
	assert.Equal(t, "EPOCH_MILLIS", logsCandidate.DateFormat)
	assert.Equal(t, 30, logsCandidate.RetentionDays)
	assert.Equal(t, []string{"service_name", "log_level"}, logsCandidate.PartitionFields)
}

func TestCollectionsForArchiving_HandlesNilClientGracefully(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	candidates := CollectionsForArchiving(ctx, nil, "project123", "cluster456")

	assert.Len(t, candidates, 2)
}

func TestCollectionsForArchiving_HandlesEmptyProjectIDAndClusterName(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	var client *admin.APIClient

	candidates := CollectionsForArchiving(ctx, client, "", "")

	assert.Len(t, candidates, 2)
}

func TestValidateCandidate_SucceedsWithValidCandidate(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   60,
		PartitionFields: []string{"created_at", "user_id"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.NoError(t, err)
}

func TestValidateCandidate_FailsWhenDatabaseNameIsEmpty(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   60,
		PartitionFields: []string{"user_id"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database name and collection name are required")
}

func TestValidateCandidate_FailsWhenCollectionNameIsEmpty(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   60,
		PartitionFields: []string{"user_id"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database name and collection name are required")
}

func TestValidateCandidate_FailsWhenBothDatabaseAndCollectionNamesAreEmpty(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "",
		CollectionName:  "",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   60,
		PartitionFields: []string{"user_id"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database name and collection name are required")
}

func TestValidateCandidate_FailsWhenRetentionDaysBelowMinimum(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   15,
		PartitionFields: []string{"user_id"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retention days must be at least 30")
}

func TestValidateCandidate_FailsWhenRetentionDaysEqualsZero(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   0,
		PartitionFields: []string{"user_id"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retention days must be at least 30")
}

func TestValidateCandidate_FailsWhenPartitionFieldsAreEmpty(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   60,
		PartitionFields: []string{},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one partition field is required")
}

func TestValidateCandidate_FailsWhenPartitionFieldsAreNil(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   60,
		PartitionFields: nil,
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one partition field is required")
}

func TestValidateCandidate_SucceedsWithCustomMinimumRetentionDays(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   45,
		PartitionFields: []string{"created_at", "user_id"},
	}
	opts := Options{
		MinimumRetentionDays: 40,
	}

	err := ValidateCandidate(candidate, opts)

	assert.NoError(t, err)
}

func TestValidateCandidate_FailsWhenRetentionDaysBelowCustomMinimum(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   35,
		PartitionFields: []string{"user_id"},
	}
	opts := Options{
		MinimumRetentionDays: 40,
	}

	err := ValidateCandidate(candidate, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retention days must be at least 40")
}

func TestValidateCandidate_SucceedsWithMultiplePartitionFields(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   60,
		PartitionFields: []string{"created_at", "user_id", "tenant_id", "category"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.NoError(t, err)
}

func TestValidateCandidate_SucceedsWithEpochMillisDateFormat(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "timestamp",
		DateFormat:      "EPOCH_MILLIS",
		RetentionDays:   90,
		PartitionFields: []string{"timestamp", "service_name"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.NoError(t, err)
}

func TestValidateCandidate_SucceedsWhenRetentionDaysEqualsMinimum(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   30,
		PartitionFields: []string{"created_at", "user_id"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.NoError(t, err)
}

func TestValidateCandidate_FailsWhenDateFieldNotInPartitionFields(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   60,
		PartitionFields: []string{"user_id", "tenant_id"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "date field created_at must be included in partition fields")
}

func TestValidateCandidate_FailsWithInvalidDateFormat(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "INVALID_FORMAT",
		RetentionDays:   60,
		PartitionFields: []string{"created_at", "user_id"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid date format: INVALID_FORMAT")
}

func TestValidateCandidate_SucceedsWithValidDateFormats(t *testing.T) {
	t.Parallel()
	validFormats := []string{"DATE", "EPOCH_SECONDS", "EPOCH_MILLIS", "EPOCH_NANOSECONDS", "OBJECTID"}

	for _, format := range validFormats {
		t.Run("format_"+format, func(t *testing.T) {
			candidate := Candidate{
				DatabaseName:    "testdb",
				CollectionName:  "testcoll",
				DateField:       "timestamp",
				DateFormat:      format,
				RetentionDays:   60,
				PartitionFields: []string{"timestamp", "user_id"},
			}
			opts := DefaultOptions()

			err := ValidateCandidate(candidate, opts)

			assert.NoError(t, err)
		})
	}
}

func TestValidateCandidate_SucceedsWhenDateFieldIsEmpty(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "",
		DateFormat:      "",
		RetentionDays:   60,
		PartitionFields: []string{"user_id"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.NoError(t, err)
}

func TestValidateCandidate_SucceedsWhenDateFieldIsFirstInPartitionFields(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   60,
		PartitionFields: []string{"created_at", "user_id", "tenant_id"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.NoError(t, err)
}

func TestValidateCandidate_SucceedsWhenDateFieldIsLastInPartitionFields(t *testing.T) {
	t.Parallel()
	candidate := Candidate{
		DatabaseName:    "testdb",
		CollectionName:  "testcoll",
		DateField:       "created_at",
		DateFormat:      "DATE",
		RetentionDays:   60,
		PartitionFields: []string{"user_id", "tenant_id", "created_at"},
	}
	opts := DefaultOptions()

	err := ValidateCandidate(candidate, opts)

	assert.NoError(t, err)
}
