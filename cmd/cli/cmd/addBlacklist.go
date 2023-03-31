package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addBlacklistCmd = &cobra.Command{
	Use:   "addBlacklist <subnet>",
	Short: "Add subnet to blacklist",
	Long:  `Add given subnet to blacklist`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("usage: abf addBlacklist <subnet>")
		} else {
			ok, mess, err := AddBlacklist(ctx, client, args[0])
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			fmt.Printf("Ok: %s, message: %s", ok, mess)
		}
	},
}

func init() {
	rootCmd.AddCommand(addBlacklistCmd)
}
