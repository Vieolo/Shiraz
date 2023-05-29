package report

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	filemanagement "github.com/vieolo/file-management"
	terminalutils "github.com/vieolo/terminal-utils"
	"golang.org/x/tools/cover"
)

// Entry point of report generation
//
// It takes the path of the `.out` file, analyze it, and generate the HTML reports
func GenHTMLReport(outPath string) error {

	// Parsing the `.out` file. Each profile is the representation of the analysis of a single file
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
		relativePath, absolutePath, _ := getFolderPath(file)

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

		thisFolder, created, index := findOrCreateFolder(folders, relativePath, absolutePath)

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

	for i := 0; i < len(folders); i++ {
		folder := folders[i]
		folder.Subfolders = append(folder.Subfolders, getSubfolders(folder, folders)...)
		folders[i] = folder
	}

	// Writing the generated HTML files
	outFolder := strings.Replace(outPath, "/coverage.out", "", 1)
	for _, fol := range folders {

		prePath := outFolder + "/" + fol.RelativePath
		filemanagement.CreateDirIfNotExists(prePath, 0777)

		newFileName := fmt.Sprintf("%v/index.html", prePath)
		iwe := os.WriteFile(newFileName, []byte(generateIndexHTMLFile(fol)), 0777)
		if iwe != nil {
			terminalutils.PrintError(iwe.Error())
		}

		for _, file := range fol.Files {
			sp := strings.Split(file.Name, "/")

			newFileName := fmt.Sprintf("%v/%v.html", prePath, strings.Replace(sp[len(sp)-1], ".go", "", 1))
			we := os.WriteFile(newFileName, []byte(generateContentHTMLFile(file)), 0777)
			if we != nil {
				terminalutils.PrintError(we.Error())
			}
		}
	}

	return nil
}
