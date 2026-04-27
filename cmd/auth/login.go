package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	authinternal "github.com/vecherochek/evo-ai-agents-skills-cli/internal/auth"
	"github.com/vecherochek/evo-ai-agents-skills-cli/internal/config"
)

var loginOptions struct {
	keyID      string
	secret     string
	endpoint   string
	projectID  string
	customerID string
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in and save IAM credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)

		keyID := strings.TrimSpace(loginOptions.keyID)
		if keyID == "" {
			keyID = strings.TrimSpace(cfg.IAMKeyID)
		}
		if keyID == "" {
			fmt.Print("IAM Key ID: ")
			value, _ := reader.ReadString('\n')
			keyID = strings.TrimSpace(value)
		}
		if keyID == "" {
			return fmt.Errorf("IAM Key ID is required")
		}

		secret := strings.TrimSpace(loginOptions.secret)
		if secret == "" {
			secret = strings.TrimSpace(cfg.IAMSecret)
		}
		if secret == "" {
			fmt.Print("IAM Secret: ")
			value, _ := reader.ReadString('\n')
			secret = strings.TrimSpace(value)
		}
		if secret == "" {
			return fmt.Errorf("IAM Secret is required")
		}

		endpoint := strings.TrimSpace(loginOptions.endpoint)
		if endpoint == "" {
			endpoint = strings.TrimSpace(cfg.IAMEndpoint)
		}
		if endpoint == "" {
			endpoint = "https://iam.api.cloud.ru"
		}

		projectID := strings.TrimSpace(loginOptions.projectID)
		if projectID == "" {
			projectID = strings.TrimSpace(cfg.ProjectID)
		}
		if projectID == "" {
			fmt.Print("Project ID (optional): ")
			value, _ := reader.ReadString('\n')
			projectID = strings.TrimSpace(value)
		}

		customerID := strings.TrimSpace(loginOptions.customerID)
		if customerID == "" {
			customerID = strings.TrimSpace(cfg.CustomerID)
		}
		if customerID == "" {
			fmt.Print("Customer ID (optional): ")
			value, _ := reader.ReadString('\n')
			customerID = strings.TrimSpace(value)
		}

		manager := authinternal.NewCredentialsManager()
		if err := manager.SaveCredentials(&authinternal.Credentials{
			IAMKeyID:     keyID,
			IAMSecretKey: secret,
			IAMEndpoint:  endpoint,
			ProjectID:    projectID,
			CustomerID:   customerID,
			LastLogin:    time.Now().Format("2006-01-02 15:04:05"),
		}); err != nil {
			return err
		}

		if err := manager.SetEnvironmentVariables(); err != nil {
			return err
		}

		fmt.Println("✅ Login successful. Credentials saved.")
		fmt.Printf("📁 %s\n", manager.GetCredentialsPath())
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVar(&loginOptions.keyID, "iam-key-id", "", "IAM Key ID")
	loginCmd.Flags().StringVar(&loginOptions.secret, "iam-secret", "", "IAM Secret")
	loginCmd.Flags().StringVar(&loginOptions.endpoint, "iam-endpoint", "", "IAM Endpoint (default: IAM_ENDPOINT env or https://iam.api.cloud.ru)")
	loginCmd.Flags().StringVar(&loginOptions.projectID, "project-id", "", "Project ID")
	loginCmd.Flags().StringVar(&loginOptions.customerID, "customer-id", "", "Customer ID")
}
