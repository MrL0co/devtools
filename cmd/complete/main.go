package complete

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	. "internal/logging"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v2"
)

var command = cli.Command{
	Name:   "complete",
	Usage:  "install autocompletion (bash, zsh)",
	Action: action,
}

func action(c *cli.Context) error {
	err := installBashComplete(c.App)
	err2 := installZshComplete(err, c.App)
	if err2 != nil {
		return err2
	}

	return nil
}

func Cmd() *cli.Command {
	return &command
}

func installZshComplete(err error, app *cli.App) error {
	s, err := downloadAutoCompleteFile("zsh_autocomplete", "", app)

	zshrc := xdg.Home + "/.zshrc"
	if _, err = os.Stat(zshrc); !os.IsNotExist(err) {
		// read the whole file at once
		b, err := ioutil.ReadFile(zshrc)
		if err != nil {
			panic(err)
		}
		contents := string(b)

		line := "PROG=" + app.Name + " _CLI_ZSH_AUTOCOMPLETE_HACK=1 source " + s
		if !strings.Contains(contents, line) {
			f, err := os.OpenFile(zshrc, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := f.WriteString("\n" + line + "\n"); err != nil {
				return err
			}

			fmt.Println("installed in " + s)
		} else {
			fmt.Println("already installed in " + s)
		}
	}
	return nil
}

func installBashComplete(app *cli.App) error {
	bashPath := "/etc/bash_completion.d/"
	file := bashPath + app.Name
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		fmt.Println("already installed in " + file)
		return nil
	}

	s, err := downloadAutoCompleteFile("bash_autocomplete", "", app)
	if err != nil {
		return err
	}

	if _, err = os.Stat(bashPath); !os.IsNotExist(err) {
		_, err = copyFile(s, file)
	}

	if os.IsPermission(err) {
		fmt.Println("could not install bash autocompletion. Rerun as sudo OR")
		fmt.Println("use it manually via: ")
		fmt.Println()
		fmt.Println("PROG=" + app.Name + "")
		fmt.Println("source " + s)
		fmt.Println()
		fmt.Println("You can add this to your .bash_rc")
		return err
	} else if err != nil {
		fmt.Println("could not automatically install bash autocompletion.")
		fmt.Println("use it manually via: ")
		fmt.Println()
		fmt.Println("PROG=" + app.Name + "")
		fmt.Println("source " + s)
		fmt.Println()
		fmt.Println("You can add this to your .bash_rc")
		return err
	}
	return nil
}

func downloadAutoCompleteFile(scriptName string, renameFileAs string, app *cli.App) (string, error) {
	if renameFileAs == "" {
		renameFileAs = scriptName
	}

	configFilePath, err := xdg.ConfigFile(app.Name + "/autocomplete/" + renameFileAs)
	if err != nil {
		return "", err
	}

	fileUrl := "https://raw.githubusercontent.com/urfave/cli/master/autocomplete/" + scriptName

	err = DownloadFile(configFilePath, fileUrl)
	if err != nil {
		return "", err
	}
	Log.Info("Downloaded: " + scriptName)

	return configFilePath, nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
