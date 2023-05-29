package report

import (
	"fmt"
	"strings"
)

// This function takes the analyzed body of a file and insert it into the
// final HTML file to be saved in to the drive
func generateContentHTMLFile(file ReportFile) string {

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
