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

	"golang.org/x/tools/cover"
)

type ReportFile struct {
	Name     string
	Path     string
	Body     template.HTML
	Coverage float64
}
type ReportFolder struct {
	Name     string
	Path     string
	Coverage float64
	Files    []ReportFile
}

func (f *ReportFolder) AddFile(p ReportFile) []ReportFile {
	f.Files = append(f.Files, p)
	return f.Files
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

func GenHTMLReport(outPath string) error {
	profiles, pErr := cover.ParseProfiles(outPath)
	if pErr != nil {
		return pErr
	}

	folders := make([]ReportFolder, 0)

	dirs, err := findPkgs(profiles)
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		fn := profile.FileName
		file, err := findFile(dirs, fn)
		if err != nil {
			return err
		}

		folderPath, _ := getFolderRelativePath(file)

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

		thisFolder.AddFile(ReportFile{
			Name:     fn,
			Path:     file,
			Body:     template.HTML(buf.String()),
			Coverage: percentCovered(profile),
		})

		if created {
			folders = append(folders, thisFolder)
		} else {
			folders[index] = thisFolder
		}
	}

	return nil
}

type Pkg struct {
	ImportPath string
	Dir        string
	Error      *struct {
		Err string
	}
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
func percentCovered(p *cover.Profile) float64 {
	var total, covered int64
	for _, b := range p.Blocks {
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	if total == 0 {
		return 0
	}
	return float64(covered) / float64(total) * 100
}

func getFolderRelativePath(p string) (string, string) {
	wd, _ := os.Getwd()

	newPath := strings.TrimPrefix(strings.Replace(p, wd, "", 1), "/")

	sp := strings.Split(newPath, "/")
	if strings.Contains(sp[len(sp)-1], ".go") {
		sp = sp[:len(sp)-1]
	}

	return strings.Join(sp, "/"), sp[len(sp)-1]
}
