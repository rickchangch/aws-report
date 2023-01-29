/*
Copyright © 2022 Rick Chang <medo972283@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/rickchangch/aws-report/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws-report",
	Short: "Provide CLI tool to generate report of AWS resources.",
	Long:  `Provide CLI tool to generate report of AWS resources.`,
}

func init() {
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return utils.ErrInvalidFlag
	})
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
