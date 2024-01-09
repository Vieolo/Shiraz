/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vieolo/shiraz/output"
	"github.com/vieolo/shiraz/utils"
	terminalutils "github.com/vieolo/terminal-utils"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs the tests",
	Long: `Runs the unit tests using the command provided in the shiraz.json file.
	If no command is provided, a standard test command is run -> go test -v ./...`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := utils.GetConfigOrDefault()

		stdout, stderr, _ := terminalutils.RunRawCommand(conf.Test.Command)

		if len(stderr.String()) > 0 {
			fmt.Println(stderr.String())
			return
		}

		// if commandErr != nil {
		// 	terminalutils.PrintError(commandErr.Error())
		// 	return
		// }

		outputType := output.PackageName
		if conf.Test.Output == "testname" {
			outputType = output.TestName
		}

		output.ParseTestOutput(stdout.String(), outputType)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
