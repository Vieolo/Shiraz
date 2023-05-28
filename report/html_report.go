package report

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	filemanagement "github.com/vieolo/file-management"
	terminalutils "github.com/vieolo/terminal-utils"
	"golang.org/x/tools/cover"
)

// Entry point of report generation
//
// It takes the path of the `.out` file, analyze it, and generate the HTML reports
func GenHTMLReport(outPath string) error {

	// Parsing the `.out` file
	// This function is the default golang function
	profiles, pErr := cover.ParseProfiles(outPath)
	if pErr != nil {
		return pErr
	}

	// Preparing the folders
	//
	// These folders represent the structure of the project
	// The generated HTML files will be placed in a similar structure to the
	// actual project. This placement makes it easier for the developer to understand the
	// degree of the coverage
	//
	// Each project has its own coverage percentage which is the average coverage of its files
	folders := make([]ReportFolder, 0)

	// Getting the directories that will be used in coverage
	dirs, err := findPkgs(profiles)
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		fn := profile.FileName

		// Finding the file in the folders
		file, err := findFile(dirs, fn)
		if err != nil {
			return err
		}

		// Getting the relative path of the folder of the file
		folderPath, _ := getFolderRelativePath(file)

		// Reading the contents of the file and generating an HTML detail
		src, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("can't read %q: %v", fn, err)
		}
		var buf strings.Builder
		err = htmlGen(&buf, src, profile.Boundaries(src))
		if err != nil {
			return err
		}

		thisFolder, created, index := findOrCreateFolder(folders, folderPath)

		coverage, coveredBlock, totalBlock := percentCovered(profile)

		thisFolder.AddFile(ReportFile{
			Name:         fn,
			Path:         file,
			Body:         template.HTML(buf.String()),
			Coverage:     coverage,
			BlockCovered: coveredBlock,
			BlockTotal:   totalBlock,
		})

		if created {
			folders = append(folders, thisFolder)
		} else {
			folders[index] = thisFolder
		}
	}

	// Writing the generated HTML files
	outFolder := strings.Replace(outPath, "/coverage.out", "", 1)
	for _, fol := range folders {
		for _, file := range fol.Files {
			sp := strings.Split(file.Name, "/")

			prePath := outFolder + "/" + fol.Path

			filemanagement.CreateDirIfNotExists(prePath, 0777)

			newFileName := fmt.Sprintf("%v/%v.html", prePath, strings.Replace(sp[len(sp)-1], ".go", "", 1))
			we := os.WriteFile(newFileName, []byte(generateCompleteHTMLFile(file)), 0777)
			if we != nil {
				terminalutils.PrintError(we.Error())
			}
		}
	}

	return nil
}

func findOrCreateFolder(folders []ReportFolder, folderPath string) (ReportFolder, bool, int) {
	for i, f := range folders {
		if f.Path == folderPath {
			return f, false, i
		}
	}

	thisName := strings.Split(folderPath, "/")
	return ReportFolder{
		Name:  thisName[len(thisName)-1],
		Path:  folderPath,
		Files: make([]ReportFile, 0),
	}, true, 0
}

func findPkgs(profiles []*cover.Profile) (map[string]*Pkg, error) {
	// Run go list to find the location of every package we care about.
	pkgs := make(map[string]*Pkg)
	var list []string
	for _, profile := range profiles {
		if strings.HasPrefix(profile.FileName, ".") || filepath.IsAbs(profile.FileName) {
			// Relative or absolute path.
			continue
		}
		pkg := path.Dir(profile.FileName)
		if _, ok := pkgs[pkg]; !ok {
			pkgs[pkg] = nil
			list = append(list, pkg)
		}
	}

	if len(list) == 0 {
		return pkgs, nil
	}

	// Note: usually run as "go tool cover" in which case $GOROOT is set,
	// in which case runtime.GOROOT() does exactly what we want.
	goTool := filepath.Join(runtime.GOROOT(), "bin/go")
	cmd := exec.Command(goTool, append([]string{"list", "-e", "-json"}, list...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	stdout, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("cannot run go list: %v\n%s", err, stderr.Bytes())
	}
	dec := json.NewDecoder(bytes.NewReader(stdout))
	for {
		var pkg Pkg
		err := dec.Decode(&pkg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decoding go list json: %v", err)
		}
		pkgs[pkg.ImportPath] = &pkg
	}
	return pkgs, nil
}

// findFile finds the location of the named file in GOROOT, GOPATH etc.
func findFile(pkgs map[string]*Pkg, file string) (string, error) {
	if strings.HasPrefix(file, ".") || filepath.IsAbs(file) {
		// Relative or absolute path.
		return file, nil
	}
	pkg := pkgs[path.Dir(file)]
	if pkg != nil {
		if pkg.Dir != "" {
			return filepath.Join(pkg.Dir, path.Base(file)), nil
		}
		if pkg.Error != nil {
			return "", errors.New(pkg.Error.Err)
		}
	}
	return "", fmt.Errorf("did not find package for %s in go list output", file)
}

// htmlGen generates an HTML coverage report with the provided filename,
// source code, and tokens, and writes it to the given Writer.
func htmlGen(w io.Writer, src []byte, boundaries []cover.Boundary) error {
	dst := bufio.NewWriter(w)
	for i := range src {
		for len(boundaries) > 0 && boundaries[0].Offset == i {
			b := boundaries[0]
			if b.Start {
				n := 0
				if b.Count > 0 {
					n = int(math.Floor(b.Norm*9)) + 1
				}
				fmt.Fprintf(dst, `<span class="cov%v" title="%v">`, n, b.Count)
			} else {
				dst.WriteString("</span>")
			}
			boundaries = boundaries[1:]
		}
		switch b := src[i]; b {
		case '>':
			dst.WriteString("&gt;")
		case '<':
			dst.WriteString("&lt;")
		case '&':
			dst.WriteString("&amp;")
		case '\t':
			dst.WriteString("        ")
		default:
			dst.WriteByte(b)
		}
	}
	return dst.Flush()
}

// percentCovered returns, as a percentage, the fraction of the statements in
// the profile covered by the test run.
// In effect, it reports the coverage of a given source file.
func percentCovered(p *cover.Profile) (float64, int64, int64) {
	var total, covered int64
	for _, b := range p.Blocks {
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	if total == 0 {
		return 0, 0, 0
	}
	return float64(covered) / float64(total) * 100, covered, total
}

func getFolderRelativePath(p string) (string, string) {
	wd, _ := os.Getwd()

	newPath := strings.TrimPrefix(strings.Replace(p, wd, "", 1), "/")

	sp := strings.Split(newPath, "/")
	if strings.Contains(sp[len(sp)-1], ".go") {
		sp = sp[:len(sp)-1]
	}

	j := strings.Join(sp, "/")

	if len(sp) == 0 {
		return j, ""
	}

	return j, sp[len(sp)-1]
}

// This function takes the analyzed body of a file and insert it into the
// final HTML file to be saved in to the drive
func generateCompleteHTMLFile(file ReportFile) string {

	// Adding line number to the `pre` tags
	n := strings.Split(string(file.Body), "\n")
	for i := 0; i < len(n); i++ {
		n[i] = fmt.Sprintf("%v    %v", i+1, n[i])
	}
	f := strings.Join(n, "\n")

	coverageClass := "success"
	if file.Coverage > 30 && file.Coverage < 80 {
		coverageClass = "alert"
	} else if file.Coverage <= 30 {
		coverageClass = "error"
	}

	temp := fmt.Sprintf(`
	<html>

		<head>
		<style>
		body {
			background: rgb(29, 29, 29);
			color: rgb(113, 113, 113);
		}
		body, pre, #legend span {
			font-family: Menlo, monospace;
			font-weight: bold;
		}
		.file-name p {
			margin: 0;
			font-size: 12px;
		}
		.coverage-header {
			height: 40px;
			display: flex;
			align-items: center;
			column-gap: 10px;
			border-bottom: 1px solid rgb(113, 113, 113);
		}
		.coverage-text {
			color: black;
			padding: 2px 5px;
		}
		.coverage-error {
			background-color: rgb(229, 85, 85);			
		}
		.coverage-alert {
			background-color: rgb(220, 207, 104);
		}
		.coverage-success {
			background-color: rgb(57, 220, 57);
		}
		#topbar {
			background: black;
			position: fixed;
			top: 0; left: 0; right: 0;
			height: 42px;
			border-bottom: 1px solid rgb(80, 80, 80);
		}
		#content {
			margin-top: 50px;
		}
		#nav, #legend {
			float: left;
			margin-left: 10px;
		}
		#legend {
			margin-top: 12px;
		}
		#nav {
			margin-top: 10px;
		}
		#legend span {
			margin: 0 5px;
		}
		.cov0 { color: rgb(192, 0, 0) }
		.cov1 { color: rgb(128, 128, 128) }
		.cov2 { color: rgb(116, 140, 131) }
		.cov3 { color: rgb(104, 152, 134) }
		.cov4 { color: rgb(92, 164, 137) }
		.cov5 { color: rgb(80, 176, 140) }
		.cov6 { color: rgb(68, 188, 143) }
		.cov7 { color: rgb(56, 200, 146) }
		.cov8 { color: rgb(44, 212, 149) }
		.cov9 { color: rgb(32, 224, 152) }
		.cov10 { color: rgb(20, 236, 155) }

	</style>
		</head>

		<body>
			<div class="file-name">
				<p>%v</p>
			</div>
			<div class="coverage-header">				
				<p class="coverage-text coverage-%v">Coverage: %v%%</p>
			</div>
			<pre>%v
			</pre>
		</body>
	</html>
	`, file.Name, coverageClass, file.Coverage, f)

	return temp
}
