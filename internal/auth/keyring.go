package auth

import (
	"github.com/zalando/go-keyring"
)

const (
	serviceName = "git-linear"
	username    = "api-key"
)

// StoreAPIKey validates and stores the API key in the system keyring.
func StoreAPIKey(key string) error {
	return keyring.Set(serviceName, username, key)
}

// GetAPIKey retrieves the stored API key from the system keyring.
// Returns keyring.ErrNotFound if the key does not exist.
func GetAPIKey() (string, error) {
	return keyring.Get(serviceName, username)
}

// DeleteAPIKey removes the stored API key from the system keyring.
// Returns nil even if the key does not exist.
func DeleteAPIKey() error {
	err := keyring.Delete(serviceName, username)
	// Ignore ErrNotFound - it's not an error if the key doesn't exist
	if err == keyring.ErrNotFound {
		return nil
	}
	return err
}

// HasAPIKey checks if an API key exists in the system keyring.
func HasAPIKey() bool {
	_, err := GetAPIKey()
	return err == nil
}
