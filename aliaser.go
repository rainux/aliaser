package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"syscall"

	shellwords "github.com/mattn/go-shellwords"
	ini "gopkg.in/ini.v1"
)

var backtickRE = regexp.MustCompile("`[^`]*`")
var aliasFile = path.Join(os.Getenv("HOME"), ".docker/alias")

func init() {
	shellwords.ParseBacktick = true
}

func main() {
	if _, err := os.Stat(aliasFile); os.IsNotExist(err) {
		execDocker()
	}

	cfg, err := ini.Load(aliasFile)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		section := cfg.Section("")
		if section.HasKey(os.Args[1]) {
			execAlias(section.Key(os.Args[1]).String())
		}
	}

	execDocker()
}

func execAlias(commandStr string) {
	result := backtickRE.ReplaceAllFunc([]byte(commandStr), parseBacktick)
	args, err := shellwords.Parse(string(result))
	if err != nil {
		log.Fatal(err)
	}

	execDocker(args)
}

func execDocker(args ...[]string) {
	var dockerArgs []string

	docker, _ := exec.LookPath("docker")

	if len(args) > 0 {
		dockerArgs = append([]string{docker}, args[0]...)
	} else {
		dockerArgs = os.Args
	}

	err := syscall.Exec(docker, dockerArgs, os.Environ())
	if err != nil {
		log.Fatalf("Error exec docker: %v", err)
	}
}

func parseBacktick(backtickStr []byte) []byte {
	result, _ := shellwords.Parse(string(backtickStr))
	return []byte(result[0])
}
