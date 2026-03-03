package config

// coalesce returns the first non-empty string from the provided arguments
func coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
