package report

import (
	"fmt"
	"strings"
)

// This function takes the analyzed body of a file and insert it into the
// final HTML file to be saved in to the drive
func generateIndexHTMLFile(fol ReportFolder) string {

	folderCoverage := fol.GetCoverage()

	folTotalCC := getCoverageClass(folderCoverage.Total)
	folFilesCC := getCoverageClass(folderCoverage.Files)
	folFoldersCC := getCoverageClass(folderCoverage.Folders)

	backButton := ""
	if fol.Name != "" {
		backButton = `<a href="../index.html"><-</a>`
	}

	files := make([]string, 0)
	for _, f := range fol.Files {
		sp := strings.Split(f.Name, "/")
		name := sp[len(sp)-1]
		fileCoverageClass := getCoverageClass(f.Coverage)
		files = append(files, fmt.Sprintf(`
		<tr>	
			<td class="file-td"><a href="./%v">%v</a></td>
			<td class="coverage-text coverage-%v">%v%%</td>
		</tr>
		`, strings.Replace(name, ".go", ".html", 1), name, fileCoverageClass, f.Coverage))
	}

	subFolders := make([]string, 0)
	for _, sub := range fol.Subfolders {
		cc := getCoverageClass(sub.GetCoverage().Total)
		subFolders = append(subFolders, fmt.Sprintf(`
		<tr>
			<td class="file-td"><a href="%v/index.html">%v</a></td>
			<td class="coverage-text coverage-%v">%v%%</td>
		</tr>
		`, sub.RelativePath, sub.Name, cc, sub.GetCoverage().Total))
	}

	subTable := ""
	if len(subFolders) > 0 {
		subTable = fmt.Sprintf(`
		<table>
			<tbody>
				<tr>
					<td class="file-td" >Subfolders</td>
					<td>Coverage</td>
				</tr>
				%v
			</tbody>
		</table>
		<br/>
		`, strings.Join(subFolders, ""))
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
		a {
			color: rgb(124, 152, 255);
			text-decoration: none;
		}
		table {
			width: 100%%;
		}

		.file-td {
			width: 350px;
		}

		.file-name {
			display: flex;
			align-items: center;
			column-gap: 10px;
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
		.coverage-none {
			display: none;
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
				%v
				<p>%v</p>
			</div>
			<div class="coverage-header">				
				<p>Coverage -> </p>
				<p class="coverage-text coverage-%v">Total: %v%%</p>
				<p class="coverage-text coverage-%v">Files: %v%%</p>
				<p class="coverage-text coverage-%v">Folders: %v%%</p>
			</div>

			%v

			<table>
				<tbody>
					
					<tr>
						<td>Files</td>
						<td>Coverage</td>
					</tr>
					
					%v
				</tbody>
			</table>
		</body>
	</html>
	`, backButton, fol.Name, folTotalCC, folderCoverage.Total, folFilesCC, folderCoverage.Files, folFoldersCC, folderCoverage.Folders, subTable, strings.Join(files, ""))

	return temp
}
