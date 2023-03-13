/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear <login> <ip>",
	Short: "Clear bucket",
	Long:  `Clear bucket for login and ip`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("usage: abf clear <login> <ip>")
		} else {
			ok, err := Clear(ctx, client, args[0], args[1])
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			fmt.Printf("Ok: %s", ok)
		}
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
