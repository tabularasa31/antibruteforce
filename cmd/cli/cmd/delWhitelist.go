package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var delWhitelistCmd = &cobra.Command{
	Use:   "delWhitelist <subnet>",
	Short: "Remove subnet from whitelist",
	Long:  `Remove given subnet from whitelist`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("usage: abf delWhitelist <subnet>")
		} else {
			ok, mess, err := DelWhitelist(ctx, client, args[0])
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			fmt.Printf("Ok: %s, message: %s", ok, mess)
		}
	},
}

func init() {
	rootCmd.AddCommand(delWhitelistCmd)
}
