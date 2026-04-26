package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Credentials struct {
	IAMKeyID     string `json:"iam_key_id"`
	IAMSecretKey string `json:"iam_secret_key"`
	IAMEndpoint  string `json:"iam_endpoint"`
	ProjectID    string `json:"project_id,omitempty"`
	CustomerID   string `json:"customer_id,omitempty"`
	LastLogin    string `json:"last_login,omitempty"`
}

type CredentialsManager struct {
	credentialsPath string
}

func NewCredentialsManager() *CredentialsManager {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return &CredentialsManager{
		credentialsPath: filepath.Join(homeDir, ".ai-agents-skills-cli", "credentials.json"),
	}
}

func (cm *CredentialsManager) LoadCredentials() (*Credentials, error) {
	if _, err := os.Stat(cm.credentialsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("credentials not found")
	}

	data, err := os.ReadFile(cm.credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("read credentials file: %w", err)
	}

	var creds Credentials
	if err = json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("decode credentials JSON: %w", err)
	}
	return &creds, nil
}

func (cm *CredentialsManager) HasCredentials() bool {
	_, err := os.Stat(cm.credentialsPath)
	return !os.IsNotExist(err)
}

func (cm *CredentialsManager) SaveCredentials(creds *Credentials) error {
	dir := filepath.Dir(cm.credentialsPath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create credentials directory: %w", err)
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("encode credentials JSON: %w", err)
	}
	if err = os.WriteFile(cm.credentialsPath, data, 0o600); err != nil {
		return fmt.Errorf("write credentials file: %w", err)
	}
	return nil
}

func (cm *CredentialsManager) DeleteCredentials() error {
	if _, err := os.Stat(cm.credentialsPath); os.IsNotExist(err) {
		return fmt.Errorf("credentials not found")
	}
	if err := os.Remove(cm.credentialsPath); err != nil {
		return fmt.Errorf("delete credentials file: %w", err)
	}
	return nil
}

func (cm *CredentialsManager) SetEnvironmentVariables() error {
	creds, err := cm.LoadCredentials()
	if err != nil {
		return err
	}

	_ = os.Setenv("IAM_KEY_ID", creds.IAMKeyID)
	_ = os.Setenv("IAM_SECRET", creds.IAMSecretKey)
	_ = os.Setenv("IAM_ENDPOINT", creds.IAMEndpoint)
	if creds.ProjectID != "" {
		_ = os.Setenv("PROJECT_ID", creds.ProjectID)
	}
	if creds.CustomerID != "" {
		_ = os.Setenv("CUSTOMER_ID", creds.CustomerID)
	}

	return nil
}

func (cm *CredentialsManager) GetCredentialsPath() string {
	return cm.credentialsPath
}
