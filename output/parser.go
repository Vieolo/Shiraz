package output

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	terminalutils "github.com/vieolo/terminal-utils"
)

type TestTrace struct {
	TestName   string
	FileName   string
	LineNumber string
	ErrorName  string
	Expected   string
	Actual     string
}

type SingleTestResult struct {
	Name         string
	IsSuccessful bool
	Time         string
}

type SinglePackageResult struct {
	IsSuccessful bool
	Name         string
	Time         string
	Tests        []SingleTestResult
}

const (
	PackageName = 1 << iota
	TestName
)

func ParseTestOutput(raw string, outputType int) {
	lines := strings.Split(raw, "\n")
	units := []SingleTestResult{}
	traces := []TestTrace{}
	results := []SinglePackageResult{}
	packageSuccessCount := 0
	packageFailCount := 0
	unitSuccessCount := 0
	unitFailCount := 0

	for i, l := range lines {
		splited := strings.Split(l, "\t")
		if len(splited) == 3 && slices.Contains([]string{"ok", "FAIL"}, strings.TrimSpace(splited[0])) {
			s := strings.TrimSpace(splited[0]) == "ok"
			results = append(results, SinglePackageResult{
				IsSuccessful: s,
				Name:         strings.TrimSpace(splited[1]),
				Time:         strings.TrimSpace(splited[2]),
				Tests:        units,
			})

			if s {
				packageSuccessCount += 1
			} else {
				packageFailCount += 1
			}

			units = []SingleTestResult{}
		} else {
			if strings.Contains(l, "--- PASS:") || strings.Contains(l, "--- FAIL:") {
				splited = strings.Split(l, " ")
				s := splited[1] == "PASS:"
				units = append(units, SingleTestResult{
					Name:         splited[2],
					IsSuccessful: s,
					Time:         splited[3],
				})

				if s {
					unitSuccessCount += 1
				} else {
					unitFailCount += 1
				}
			} else {
				if strings.Contains(l, "Error Trace:") {
					thisTrace := TestTrace{}
					thisTrace.FileName = strings.Split(splited[2], ":")[0]
					thisTrace.LineNumber = strings.Split(splited[2], ":")[1]

					for k := i; k < len(lines); k++ {
						kp := strings.Split(lines[k], "\t")
						if len(kp) < 3 {
							continue
						}
						kp1 := strings.TrimSpace(kp[1])

						if kp1 == "Error:" {
							thisTrace.ErrorName = strings.Replace(kp[2], ":", "", -1)
						} else if kp1 == "Test:" {
							thisTrace.TestName = kp[2]
							break
						} else if strings.Contains(kp[2], "expected:") {
							thisTrace.Expected = strings.Split(kp[2], ": ")[1]
						} else if strings.Contains(kp[2], "actual  :") {
							thisTrace.Actual = strings.Split(kp[2], ": ")[1]
						}
					}
					traces = append(traces, thisTrace)

				}
			}
		}
	}
	const padding = 3
	packageWriter := tabwriter.NewWriter(
		os.Stdout,
		0,
		0, padding,
		' ',
		0,
	)

	for _, res := range results {
		statusText := "✓"
		statusColor := "\u001b[32m"
		if !res.IsSuccessful {
			statusText = "x"
			statusColor = "\u001b[31m"
		}
		fmt.Fprintf(packageWriter, "%v%v\033[0m\t%v\t%v\n", statusColor, statusText, res.Name, res.Time)

		if len(res.Tests) > 0 && outputType == TestName {
			for _, unit := range res.Tests {
				statusText := "PASS"
				statusColor := "\u001b[32m"
				if !unit.IsSuccessful {
					statusText = "FAIL"
					statusColor = "\u001b[31m"
				}
				fmt.Fprintf(packageWriter, "|___ %v%v\033[0m\t> %v\t%v\n", statusColor, statusText, unit.Name, unit.Time)
			}
		}
	}

	packageWriter.Flush()

	if len(traces) > 0 {
		fmt.Println("-----------------------")
		fmt.Println(" ")
		fmt.Println("Error Traces")

		for _, t := range traces {
			terminalutils.PrintError(fmt.Sprintf(" - %v -> %v", t.TestName, t.ErrorName))
			terminalutils.PrintColorln(fmt.Sprintf("\tExpected\t%v", t.Expected), terminalutils.Yellow)
			terminalutils.PrintError(fmt.Sprintf("\tActual  \t%v", t.Actual))
			fmt.Printf("\tFile    \t%v\n", t.FileName)
			fmt.Printf("\tLine    \t%v\n", t.LineNumber)
			fmt.Println(" ")
		}

		fmt.Println(" ")
		fmt.Println("-----------------------")
	}

	fmt.Println("--------------------")
	fmt.Println("Summary")
	if packageFailCount == 0 {
		terminalutils.PrintSuccess("All Passed")
	} else {
		if outputType == PackageName {
			terminalutils.PrintError(fmt.Sprintf("%v test(s) failed out of %v", packageFailCount, packageFailCount+packageSuccessCount))
		} else {
			terminalutils.PrintError(fmt.Sprintf("%v test(s) failed out of %v", unitFailCount, unitFailCount+unitSuccessCount))
		}
	}
	fmt.Println(" ")
}
