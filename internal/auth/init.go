package auth

import "os"

func InitCredentials() {
	manager := NewCredentialsManager()
	if !manager.HasCredentials() {
		return
	}

	creds, err := manager.LoadCredentials()
	if err != nil {
		return
	}

	if os.Getenv("IAM_KEY_ID") == "" && creds.IAMKeyID != "" {
		_ = os.Setenv("IAM_KEY_ID", creds.IAMKeyID)
	}
	if os.Getenv("IAM_SECRET") == "" && creds.IAMSecretKey != "" {
		_ = os.Setenv("IAM_SECRET", creds.IAMSecretKey)
	}
	if os.Getenv("IAM_ENDPOINT") == "" && creds.IAMEndpoint != "" {
		_ = os.Setenv("IAM_ENDPOINT", creds.IAMEndpoint)
	}
	if os.Getenv("PROJECT_ID") == "" && creds.ProjectID != "" {
		_ = os.Setenv("PROJECT_ID", creds.ProjectID)
	}
}
