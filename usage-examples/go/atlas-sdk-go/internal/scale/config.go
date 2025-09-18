package scale

import (
	"atlas-sdk-go/internal/config"
)

// Expose ScalingConfig within scale package for tests and callers while reusing config.ScalingConfig.
type ScalingConfig = config.ScalingConfig

const (
	defaultTargetTier    = "M50"
	defaultCPUThreshold  = 75.0
	defaultPeriodMinutes = 60
)

// LoadScalingConfig loads programmatic scaling configuration with sensible defaults.
// Defaults are applied for missing optional fields to align with Atlas auto-scaling guidance.
func LoadScalingConfig(cfg config.Config) config.ScalingConfig {
	sc := cfg.Scaling

	// Apply defaults for missing values
	if sc.TargetTier == "" {
		sc.TargetTier = defaultTargetTier
	}

	if sc.CPUThreshold == 0 {
		sc.CPUThreshold = defaultCPUThreshold
	}

	if sc.PeriodMinutes == 0 {
		sc.PeriodMinutes = defaultPeriodMinutes
	}

	return sc
}

// DefaultScalingConfig returns ScalingConfig with sensible defaults aligned with Atlas auto-scaling guidance
func DefaultScalingConfig() config.ScalingConfig {
	return config.ScalingConfig{
		TargetTier:    defaultTargetTier,
		PreScale:      false,
		CPUThreshold:  defaultCPUThreshold,
		PeriodMinutes: defaultPeriodMinutes,
		DryRun:        true, // Safe default
	}
}
