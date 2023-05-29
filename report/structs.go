package report

import "html/template"

type ReportFile struct {
	Name         string
	Path         string
	Body         template.HTML
	Coverage     float64
	BlockCovered int64
	BlockTotal   int64
}

func (f *ReportFolder) AddFile(p ReportFile) []ReportFile {
	f.Files = append(f.Files, p)
	return f.Files
}

type ReportFolder struct {
	Name         string
	Path         string
	Coverage     float64
	BlockCovered int64
	BlockTotal   int64
	Subfolders   []ReportFolder
	Files        []ReportFile
}

type ReportFolderCoverage struct {
	Total   float64
	Files   float64
	Folders float64
}

func (f ReportFolder) GetCoverage() ReportFolderCoverage {
	var filesCoverage float64 = 0
	for _, file := range f.Files {
		filesCoverage += file.Coverage
	}
	filesCoverage = filesCoverage / float64(len(f.Files))

	return ReportFolderCoverage{
		Files: filesCoverage,
		Total: filesCoverage,
	}
}

type Pkg struct {
	ImportPath string
	Dir        string
	Error      *struct {
		Err string
	}
}
