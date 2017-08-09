package controller

import (
	"fmt"

	api_v1 "k8s.io/client-go/pkg/api/v1"
)

// ValidateTLSSecret validates the secret. If it is valid, the function returns nil.
func ValidateTLSSecret(secret *api_v1.Secret) error {
	if _, exists := secret.Data[api_v1.TLSCertKey]; !exists {
		return fmt.Errorf("Secret doesn't have %v", api_v1.TLSCertKey)
	}

	if _, exists := secret.Data[api_v1.TLSPrivateKeyKey]; !exists {
		return fmt.Errorf("Secret doesn't have %v", api_v1.TLSPrivateKeyKey)
	}

	return nil
}
