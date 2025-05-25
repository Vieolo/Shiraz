package report

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/tools/cover"
)

func findOrCreateFolder(folders []ReportFolder, relativePath string, absolutePath string) (ReportFolder, bool, int) {
	for i, f := range folders {
		if f.RelativePath == relativePath {
			return f, false, i
		}
	}

	thisName := strings.Split(relativePath, "/")
	return ReportFolder{
		Name:         thisName[len(thisName)-1],
		RelativePath: relativePath,
		AbsolutePath: absolutePath,
		Files:        make([]ReportFile, 0),
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

func getFolderPath(p string) (string, string, string) {
	absPathArr := strings.Split(p, "/")
	absPathArr = absPathArr[:len(absPathArr)-1]
	absPath := strings.Join(absPathArr, "/")

	wd, _ := os.Getwd()

	newPath := strings.TrimPrefix(strings.Replace(p, wd, "", 1), "/")

	sp := strings.Split(newPath, "/")
	if strings.Contains(sp[len(sp)-1], ".go") {
		sp = sp[:len(sp)-1]
	}

	j := strings.Join(sp, "/")

	if len(sp) == 0 {
		return j, absPath, ""
	}

	return j, absPath, sp[len(sp)-1]
}

func getCoverageClass(cov float64) string {
	coverageClass := "success"
	if cov == -1 {
		coverageClass = "none"
	} else if cov > 30 && cov < 80 {
		coverageClass = "alert"
	} else if cov <= 30 {
		coverageClass = "error"
	}

	return coverageClass
}

func getSubfolders(currentFolder ReportFolder, allFolders []ReportFolder) []ReportFolder {
	subs := make([]ReportFolder, 0)
	mainDivCount := len(strings.Split(currentFolder.AbsolutePath, "/"))

	for _, fol := range allFolders {
		if fol.AbsolutePath == currentFolder.AbsolutePath {
			continue
		}

		divCount := len(strings.Split(fol.AbsolutePath, "/"))

		if divCount <= mainDivCount {
			continue
		}

		if strings.Contains(fol.AbsolutePath, currentFolder.AbsolutePath) {
			subs = append(subs, fol)
		}
	}
	return subs
}

