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
	Files        []ReportFile
}

type Pkg struct {
	ImportPath string
	Dir        string
	Error      *struct {
		Err string
	}
}
