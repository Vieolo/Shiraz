/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/vieolo/shiraz/report"
	"github.com/vieolo/shiraz/utils"
	tu "github.com/vieolo/terminal-utils"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Produces a report for the coverage",
	Long:  `This command runs the tests and generate the out file (via standard go tool) and generates a report in the coverage folder`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := utils.GetConfigOrDefault()

		projPath := "./..."
		if conf.ProjectPath != "." {
			projPath = conf.ProjectPath
		}

		outPath := fmt.Sprintf("%vcoverage.out", conf.CoverageFolderPath)

		cArgs := []string{
			"test",
			fmt.Sprintf("-coverprofile=%v", outPath),
			projPath,
		}

		c := exec.Command("go", cArgs...)
		runErr := c.Run()
		if runErr != nil {
			tu.PrintError(runErr.Error())
			return
		}

		genErr := report.GenHTMLReport(outPath)
		if genErr != nil {
			tu.PrintError(genErr.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
