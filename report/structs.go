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
	RelativePath string
	AbsolutePath string
	Coverage     float64
	BlockCovered int64
	BlockTotal   int64
	Subfolders   []ReportFolder
	Files        []ReportFile
}

type ReportFolderCoverage struct {
	Total          float64
	UndividedTotal float64
	Files          float64
	Folders        float64
	TotalFileCount int
}

func (f ReportFolder) GetCoverage() ReportFolderCoverage {

	var total float64 = 0
	var fileCount int = 0

	var filesCoverage float64 = 0
	for _, file := range f.Files {
		filesCoverage += file.Coverage
		total += file.Coverage
		fileCount += 1
	}
	filesCoverage = filesCoverage / float64(len(f.Files))

	var folCoverage float64 = 0
	if len(f.Subfolders) > 0 {
		for _, sub := range f.Subfolders {
			sc := sub.GetCoverage()
			folCoverage += sc.Total
			total += sc.UndividedTotal
			fileCount += sc.TotalFileCount
		}
		folCoverage = folCoverage / float64(len(f.Subfolders))
	} else {
		folCoverage = -1
	}

	return ReportFolderCoverage{
		Files:          filesCoverage,
		Total:          total / float64(fileCount),
		UndividedTotal: total,
		Folders:        folCoverage,
		TotalFileCount: fileCount,
	}
}

type Pkg struct {
	ImportPath string
	Dir        string
	Error      *struct {
		Err string
	}
}
