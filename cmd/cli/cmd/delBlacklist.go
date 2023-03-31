package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var delBlacklistCmd = &cobra.Command{
	Use:   "delBlacklist <subnet>",
	Short: "Remove subnet from blacklist",
	Long:  `Remove given subnet from blacklist`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("usage: abf delBlacklist <subnet>")
		} else {
			ok, mess, err := DelBlacklist(ctx, client, args[0])
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			fmt.Printf("Ok: %s, message: %s", ok, mess)
		}
	},
}

func init() {
	rootCmd.AddCommand(delBlacklistCmd)
}
