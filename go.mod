module devtools

go 1.17

require internal/update v1.0.0

require (
	github.com/adrg/xdg v0.4.0
	github.com/urfave/cli/v2 v2.3.0
	internal/logging v1.0.0
)

replace internal/update => ./internal/update

replace internal/logging => ./internal/logging

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mouuff/go-rocket-update v1.5.2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
)
