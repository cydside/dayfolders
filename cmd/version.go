package cmd

//______________________________________________________________________________

import (
	"fmt"

	"github.com/spf13/cobra"
)

//______________________________________________________________________________

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of dayfolders",
	Long:  `All software has versions. This is 'dayfolders'`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("dayfolders args: %v\n", args)
		fmt.Println("dayfolders a CLI that creates daily folders v1.0.0")
	},
}
