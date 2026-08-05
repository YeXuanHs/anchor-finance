package model

// GenerateAPIKey generates a random hex API key (64 chars).
// Kept for backward compatibility; new code uses GenerateRandomPassword in user.go.
func GenerateAPIKey() string {
	return GenerateRandomPassword(16)
}
