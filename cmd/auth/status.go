package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	authinternal "github.com/vecherochek/evo-ai-agents-skills-cli/internal/auth"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := authinternal.NewCredentialsManager()
		if !manager.HasCredentials() {
			fmt.Println("❌ Credentials are not configured.")
			fmt.Println("💡 Run: skills-cli auth login")
			return nil
		}

		creds, err := manager.LoadCredentials()
		if err != nil {
			return err
		}

		fmt.Println("✅ Credentials are configured:")
		fmt.Printf("🔑 IAM Key ID: %s\n", maskString(creds.IAMKeyID))
		fmt.Printf("🌐 IAM Endpoint: %s\n", creds.IAMEndpoint)
		if creds.ProjectID != "" {
			fmt.Printf("📋 Project ID: %s\n", creds.ProjectID)
		}
		if creds.LastLogin != "" {
			fmt.Printf("⏰ Last login: %s\n", creds.LastLogin)
		}
		return nil
	},
}
