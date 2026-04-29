package skill_marketplace

import "github.com/spf13/cobra"

var RootCMD = &cobra.Command{
	Use:   "skill-marketplace",
	Short: "Marketplace skill operations",
}

func init() {
	RootCMD.AddCommand(AddCMD)
}
