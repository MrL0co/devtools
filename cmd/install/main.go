package install

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"
)

var command = cli.Command{
	Name:    "install",
	Aliases: []string{"i"},
	Usage:   "Install from a predefined list of suggested packages",
	Action:  action,
}

var cout chan []byte = make(chan []byte)
var cin chan []byte = make(chan []byte)
var exit chan bool = make(chan bool)

func Foo(x byte) byte { return call_port([]byte{1, x}) }
func Bar(y byte) byte { return call_port([]byte{2, y}) }
func Exit() byte      { return call_port([]byte{0, 0}) }
func call_port(s []byte) byte {
	cout <- s
	s = <-cin
	return s[1]
}

func start(command string) {
	fmt.Println("start")
	cmd := exec.Command(command)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdout, err2 := cmd.StdoutPipe()
	if err2 != nil {
		log.Fatal(err2)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	defer stdin.Close()
	defer stdout.Close()
	for {
		select {
		case s := <-cout:
			stdin.Write(s)
			buf := make([]byte, 2)
			runtime.Gosched()
			time.Sleep(100 * time.Millisecond)
			stdout.Read(buf)
			cin <- buf
		case b := <-exit:
			if b {
				fmt.Printf("Exit")
				return //os.Exit(0)
			}
		}
	}
}

func action(c *cli.Context) error {
	output, err := exec.Command("apt", "update").CombinedOutput()
	log.Println(string(output[:])) // write the output with ResponseWriter
	//if err != nil {
	//	return err
	//}

	output, err = exec.Command("apt", "install", "-y", "zsh", "git").CombinedOutput()
	log.Println(string(output[:])) // write the output with ResponseWriter
	//if err != nil {
	//	return err
	//}

	current, err := user.Current()
	log.Println(os.Getenv("SUDO_USER"))
	log.Println(current.Username)

	zsh, err := exec.Command("which", "zsh").Output()

	cmd := exec.Command("chsh", "-s", strings.Trim(string(zsh), " \n"), current.Username)
	//stdin, err := cmd.StdinPipe()

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err != nil {
		log.Fatal(err)
	}

	//go func() {
	//	defer stdin.Close()
	//	io.WriteString(stdin, "test123")
	//}()

	log.Println(cmd.String())
	//output, err = cmd.CombinedOutput()
	//if err != nil {
	//	return err
	//}
	//log.Println(string(output[:])) // write the output with ResponseWriter
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	cmd = exec.Command("zsh", "-c", "wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh")
	log.Println(cmd.String())
	output, err = cmd.CombinedOutput()
	log.Println(string(output[:])) // write the output with ResponseWriter
	if err != nil {
		return err
	}

	//go start()
	//runtime.Gosched()
	//fmt.Println("30+1=", Foo(30)) //30+1= 31
	//fmt.Println("2*40=", Bar(40)) //2*40= 80
	//Exit()
	//exit <- true
	return nil
}

func Cmd() *cli.Command {
	return &command
}
