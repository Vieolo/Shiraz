# Shiraz
CLI for testing and coverage of Go projects

## HTML Report
The generated HTML files are of two types; content and index.

The content files display the coverage of a single file.

The index files display the files and nested folders in a folder. The index files allow the viewer to navigate the files and provides an average coverage of the nested files and folders.

TODO:

- Generate index files
- Place the generated files in the right folders
- Style the generated files
- Add number of execution to the content files
- Add support for a config file in the target project
- Add multiple types of test runs (normal, with coverage, etc.)
- Process and refactor the output of the tests to be more readable
- Add support for the features of `gotestsum`
- Add a coverage summary in the terminal output
