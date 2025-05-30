/*
Copyright © 2025
*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vieolo/shiraz/utils"
	tu "github.com/vieolo/terminal-utils"
	"strings"
	"fmt"
)

var runCmd = &cobra.Command{
	Use: "run",
	Short: "Runs scripts defined in the config",
	Long: "This command runs scripts in the scripts section of shiraz.json, addressed by the key value",
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		conf := utils.GetConfigOrDefault()

		projPath := "./..."
		if conf.ProjectPath != "" && conf.ProjectPath != "." {
			projPath = conf.ProjectPath
		}
		var _ = projPath

		cArgs := []string{
			"-c",
			conf.Scripts[args[0]],
		}
		cmdString := strings.Join(cArgs, " ")
	        fmt.Println("sh", cmdString)

		stdout, stderr, commandErr := tu.RunCommand(tu.CommandConfig{
			Command: "sh",
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
	},
}


func init() {
	rootCmd.AddCommand(runCmd)
}
