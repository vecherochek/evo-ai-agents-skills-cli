package auth

import (
	"fmt"
	"os"

	authinternal "github.com/cloud-ru/evo-ai-agents-skills-cli/internal/auth"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out and remove saved credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := authinternal.NewCredentialsManager()
		if !manager.HasCredentials() {
			fmt.Println("ℹ️  Credentials were not found.")
			return nil
		}
		if err := manager.DeleteCredentials(); err != nil {
			return err
		}

		_ = os.Unsetenv("IAM_KEY_ID")
		_ = os.Unsetenv("IAM_SECRET")
		_ = os.Unsetenv("IAM_ENDPOINT")
		_ = os.Unsetenv("PROJECT_ID")
		_ = os.Unsetenv("CUSTOMER_ID")

		fmt.Println("✅ Logout successful. Credentials removed.")
		return nil
	},
}
