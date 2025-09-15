package recovery

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"atlas-sdk-go/internal/typeutils"
)

const (
	defaultAddNodes        = 1
	scenarioRegionalOutage = "regional-outage"
	scenarioDataDeletion   = "data-deletion"
)

// DrOptions holds the scenario and configuration parameters used by the
// disaster recovery example. Values are typically loaded from environment
// variables. Only the fields relevant to the chosen Scenario are required.
//
//	Scenario values:
//	  - "regional-outage" : simulate adding capacity to a healthy region
//	  - "data-deletion"   : submit a snapshot restore job
//	Required per scenario:
//	  regional-outage: ProjectID, ClusterName, TargetRegion
//	  data-deletion:   ProjectID, ClusterName, SnapshotID
//	Optional:
//	  OutageRegion (regional-outage) region to zero electable nodes
//	  AddNodes (regional-outage) number of electable nodes to add (default 1)
//	  DryRun when true prints intended actions only.
type DrOptions struct {
	Scenario     string
	ProjectID    string
	ClusterName  string
	TargetRegion string
	OutageRegion string
	AddNodes     int
	SnapshotID   string
	DryRun       bool
}

// LoadDROptionsFromEnv reads environment variables and validates scenario-specific requirements.
// Defaults are applied first, then overridden if env vars are present:
//
//	DR_SCENARIO          (req) regional-outage | data-deletion
//	ATLAS_PROJECT_ID     (red unless provided via config loader)
//	ATLAS_CLUSTER_NAME   (req) target cluster name
//	DR_TARGET_REGION     (regional-outage req) region receiving added capacity
//	DR_OUTAGE_REGION     (regional-outage opt) region considered impaired (its electable nodes set to 0)
//	DR_ADD_NODES         (regional-outage opt) number of electable nodes to add (default: 1)
//	DR_SNAPSHOT_ID       (data-deletion req) snapshot ID to restore
//	DR_DRY_RUN           (opt bool) if true, only log intended actions (default: false)
func LoadDROptionsFromEnv(fallbackProjectID string) (DrOptions, error) {
	o := DrOptions{
		AddNodes: defaultAddNodes,
	}

	o.Scenario = strings.ToLower(strings.TrimSpace(os.Getenv("DR_SCENARIO")))
	o.ProjectID = typeutils.FirstNonEmpty(os.Getenv("ATLAS_PROJECT_ID"), fallbackProjectID)
	o.ClusterName = strings.TrimSpace(os.Getenv("ATLAS_CLUSTER_NAME"))
	o.TargetRegion = strings.TrimSpace(os.Getenv("DR_TARGET_REGION"))
	o.OutageRegion = strings.TrimSpace(os.Getenv("DR_OUTAGE_REGION"))
	o.SnapshotID = strings.TrimSpace(os.Getenv("DR_SNAPSHOT_ID"))

	if v, ok := os.LookupEnv("DR_ADD_NODES"); ok {
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return o, fmt.Errorf("invalid DR_ADD_NODES value '%s': must be a positive integer", v)
		}
		if n <= 0 {
			return o, fmt.Errorf("DR_ADD_NODES must be a positive integer, got %d", n)
		}
		o.AddNodes = n
	}

	if v, ok := os.LookupEnv("DR_DRY_RUN"); ok {
		o.DryRun = typeutils.ParseBool(v)
	}
	if err := validateRequiredFields(o); err != nil {
		return o, err
	}
	if err := validateScenarioRequirements(o); err != nil {
		return o, err
	}

	return o, nil
}

func validateRequiredFields(o DrOptions) error {
	if o.Scenario == "" {
		return fmt.Errorf("DR_SCENARIO is required")
	}
	if o.ProjectID == "" {
		return fmt.Errorf("ATLAS_PROJECT_ID is required")
	}
	if o.ClusterName == "" {
		return fmt.Errorf("ATLAS_CLUSTER_NAME is required")
	}
	return nil
}

// validateScenarioRequirements checks that scenario-specific required fields are set.
func validateScenarioRequirements(o DrOptions) error {
	switch o.Scenario {
	case scenarioRegionalOutage:
		if o.TargetRegion == "" {
			return fmt.Errorf("DR_TARGET_REGION is required for %s scenario", scenarioRegionalOutage)
		}
	case scenarioDataDeletion:
		if o.SnapshotID == "" {
			return fmt.Errorf("DR_SNAPSHOT_ID is required for %s scenario", scenarioDataDeletion)
		}
	default:
		return fmt.Errorf("unsupported DR_SCENARIO '%s': valid options are %s, %s",
			o.Scenario, scenarioRegionalOutage, scenarioDataDeletion)
	}
	return nil
}
