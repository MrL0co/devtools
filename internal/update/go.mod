module update

go 1.17

require internal/logging v1.0.0

replace internal/logging => ../logging

require github.com/mouuff/go-rocket-update v1.5.2

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
)
