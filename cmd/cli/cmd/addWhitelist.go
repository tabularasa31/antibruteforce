package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addWhitelistCmd = &cobra.Command{
	Use:   "addWhitelist <subnet>",
	Short: "Add subnet to whitelist",
	Long:  `Add given subnet to whitelist`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("usage: abf addWhitelist <subnet>")
		} else {
			ok, mess, err := AddWhitelist(ctx, client, args[0])
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			fmt.Printf("Ok: %s, message: %s", ok, mess)
		}
	},
}

func init() {
	rootCmd.AddCommand(addWhitelistCmd)
}
