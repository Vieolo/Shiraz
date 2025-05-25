/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	fm "github.com/vieolo/file-management"
	"github.com/vieolo/shiraz/browser"
	"github.com/vieolo/shiraz/report"
	"github.com/vieolo/shiraz/utils"
	tu "github.com/vieolo/terminal-utils"
	"strings"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Produces a report for the coverage",
	Long:  `This command runs the tests and generate the out file (via standard go tool) and generates a report in the coverage folder`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := utils.GetConfigOrDefault()

		projPath := "./..."
		if conf.ProjectPath != "" && conf.ProjectPath != "." {
			projPath = conf.ProjectPath
		}

		outPath := fmt.Sprintf("%vcoverage.out", conf.CoverageFolderPath)
		re := os.RemoveAll(conf.CoverageFolderPath)
		if re != nil {
			tu.PrintError(re.Error())
		}
		fm.CreateDirIfNotExists(conf.CoverageFolderPath, 0777)

		// go test -v -coverpkg=./... -coverprofile=coverage/coverage.out ./...
		cArgs := []string{
			"test",
			"-v",
			fmt.Sprintf("-coverpkg=%v", projPath),
			fmt.Sprintf("-coverprofile=%v", outPath),
			fmt.Sprintf("%v/...", projPath),
		}
		cmdString := strings.Join(cArgs, " ")
	        fmt.Println(cmdString)

		stdout, stderr, commandErr := tu.RunCommand(tu.CommandConfig{
			Command: "go",
			Args:    cArgs,
			Env:     conf.Env,
		})

		fmt.Println(stdout.String())
		if len(stderr.String()) > 0 {
			fmt.Println(stderr.String())
		}

		if commandErr != nil {
			tu.PrintError(commandErr.Error())
		}

		genErr := report.GenHTMLReport(outPath, conf)
		if genErr != nil {
			tu.PrintError(genErr.Error())
			return
		}

		browser.Open(conf.CoverageFolderPath + "/index.html")
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
