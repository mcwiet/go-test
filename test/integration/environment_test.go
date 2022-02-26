package integration_test

import "os"

// Attempt to load environment variable; panic if not found (fail fast)
func GetRequiredEnv(name string) string {
	val, exists := os.LookupEnv(name)
	if !exists {
		panic("Could not load environment variable: " + name)
	}
	return val
}
