/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	filemanagement "github.com/vieolo/file-management"
	"github.com/vieolo/shiraz/browser"
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
		re := os.RemoveAll(conf.CoverageFolderPath)
		if re != nil {
			tu.PrintError(re.Error())
		}
		filemanagement.CreateDirIfNotExists(conf.CoverageFolderPath, 0777)

		// go test -v -coverpkg=./... ./... -coverprofile=coverage.out ./...
		cArgs := []string{
			"test",
			"-v",
			fmt.Sprintf("-coverpkg=%v", projPath),
			projPath,
			fmt.Sprintf("-coverprofile=%v", outPath),
			projPath,
		}

		c := exec.Command("go", cArgs...)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		c.Stdout = &stdout
		c.Stderr = &stderr

		if len(conf.Env) > 0 {
			c.Env = os.Environ()

			for k, v := range conf.Env {
				c.Env = append(c.Env, fmt.Sprintf("%v=%v", k, v))
			}
		}

		runErr := c.Run()
		if runErr != nil {
			tu.PrintError(runErr.Error())
			return
		}
		fmt.Println(stdout.String())
		if len(stderr.String()) > 0 {
			fmt.Println(stderr.String())
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
