package auth

import (
	"fmt"
	"os"

	authinternal "github.com/cloud-ru/evo-ai-agents-skills-cli/internal/auth"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current auth configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := authinternal.NewCredentialsManager()

		fmt.Println("📁 Credentials file:")
		fmt.Printf("%s\n\n", manager.GetCredentialsPath())

		if manager.HasCredentials() {
			creds, err := manager.LoadCredentials()
			if err != nil {
				return err
			}
			fmt.Println("✅ Saved credentials:")
			fmt.Printf("IAM Key ID: %s\n", maskString(creds.IAMKeyID))
			fmt.Printf("IAM Endpoint: %s\n", creds.IAMEndpoint)
			fmt.Printf("Project ID: %s\n\n", creds.ProjectID)
		} else {
			fmt.Println("❌ Saved credentials: not found")
			fmt.Println()
		}

		fmt.Println("🔎 Active environment:")
		printEnv("IAM_KEY_ID", true)
		printEnv("IAM_SECRET", true)
		printEnv("IAM_ENDPOINT", false)
		printEnv("PROJECT_ID", false)
		printEnv("CUSTOMER_ID", false)
		return nil
	},
}

func printEnv(name string, mask bool) {
	value := os.Getenv(name)
	if value == "" {
		fmt.Printf("%s: <empty>\n", name)
		return
	}
	if mask {
		fmt.Printf("%s: %s\n", name, maskString(value))
		return
	}
	fmt.Printf("%s: %s\n", name, value)
}

func maskString(value string) string {
	if len(value) <= 8 {
		return "********"
	}
	return value[:4] + "****" + value[len(value)-4:]
}
