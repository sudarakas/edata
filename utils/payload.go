package utils

func GetOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func GetOrDefaultFloat(value, fallback float64) float64 {
	if value == 0 {
		return fallback
	}
	return value
}

func GetOrDefaultInt(value, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}
