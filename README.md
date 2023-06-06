# Shiraz
CLI for testing and coverage of Go projects

At the moment, shiraz runs unit tests and generates an improved version of HTML coverage report by running the `report` command.

You can define the configuration of your project by creating a `shiraz.json` in your main directory.


## shiraz.json

All fields in the `shiraz.json` are optional. Here are the possible fields.

- `projectPath`: The path to the go project. Useful if the config file is not in the project being tested.
- `coverageFolderPath`: The path to the folder where the coverage files are generated and saved at
- `env`: The environmental variables to be added when running the test command.
- `ignore`: An array of files of folders you wish to ignore from the report. You need to include the package name as well. e.g. `github.com/example/dir_1` or `github.com/example/dir_2/file_1.go`