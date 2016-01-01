package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"syscall"

	"github.com/mattn/go-shellwords"
	"gopkg.in/ini.v1"
)

var (
	backtickRE = regexp.MustCompile("`[^`]*`")
	aliasFile  = path.Join(os.Getenv("HOME"), ".aliaser")
	target     string
)

func init() {
	shellwords.ParseBacktick = true

	if _, err := os.Stat(aliasFile); os.IsNotExist(err) {
		log.Fatal("Configure your aliases in ~/.aliaser first.")
	}
}

func main() {
	loadTarget()

	if len(os.Args) > 1 {
		config, err := ini.Load(aliasFile)
		if err != nil {
			log.Fatal(err)
		}

		section, err := config.GetSection(target)
		if err == nil {
			if section.HasKey(os.Args[1]) {
				execAlias(section.Key(os.Args[1]).String())
			}
		}
	}

	execTarget()
}

func execAlias(commandStr string) {
	result := backtickRE.ReplaceAllFunc([]byte(commandStr), parseBacktick)
	args, err := shellwords.Parse(string(result))
	if err != nil {
		log.Fatal(err)
	}

	execTarget(args)
}

func execTarget(args ...[]string) {
	targetPath, _ := exec.LookPath(target)
	targetArgs := []string{targetPath}

	if len(args) > 0 {
		targetArgs = append(targetArgs, args[0]...)
	} else {
		targetArgs = append(targetArgs, os.Args[1:]...)
	}

	err := syscall.Exec(targetPath, targetArgs, os.Environ())
	if err != nil {
		log.Fatalf("Error exec %v: %v", target, err)
	}
}

func parseBacktick(backtickStr []byte) []byte {
	result, _ := shellwords.Parse(string(backtickStr))
	return []byte(strings.Join(result, ""))
}

func loadTarget() {
	config, err := ini.Load(aliasFile)
	if err != nil {
		log.Fatal(err)
	}

	myName := path.Base(os.Args[0])
	section := config.Section("core")
	if section.HasKey(myName) {
		target = section.Key(myName).String()
	}
}
