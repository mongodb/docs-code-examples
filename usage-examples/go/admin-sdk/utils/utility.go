package utils

// Utility functions
func FirstNonEmpty(cli, config string) string {
	if cli != "" {
		return cli
	}
	return config
}
func FirstNonEmptyArray(cli, config []string) []string {
	if len(cli) > 0 {
		return cli
	}
	return config
}

func FirstNonZero(cli, config int) int {
	if cli != 0 {
		return cli
	}
	return config
}
