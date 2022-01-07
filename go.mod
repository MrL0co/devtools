module devtools

go 1.17

require internal/update v1.0.0
replace internal/update => ./internal/update

require (
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/mouuff/go-rocket-update v1.5.2 // indirect
	github.com/sanbornm/go-selfupdate v0.0.0-20210106163404-c9b625feac49 // indirect
)
