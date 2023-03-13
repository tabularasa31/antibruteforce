package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var requestCmd = &cobra.Command{
	Use:   "request <login> <pass> <ip>",
	Short: "Send request for login, password and ip",
	Long:  `Send request for login, password and ip if there are bruteforce or not`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Println("usage: abf request <login> <pass> <ip>")
		} else {
			ok, mess, err := Request(ctx, client, args[0], args[1], args[2])
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			fmt.Printf("Ok: %s, message: %s", ok, mess)
		}
	},
}

func init() {
	rootCmd.AddCommand(requestCmd)
}
