package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var composeDirName, baseCompose, env string
var debug bool

func main() {
	flag.StringVar(&composeDirName, "composeDirName", ".compose", "Name of directory containing docker compose files")
	flag.StringVar(&baseCompose, "baseCompose", "docker-compose.default.yml", "The base docker compose file")
	flag.StringVar(&env, "env", "devel", "Environment that docker compose is running in")
	flag.BoolVar(&debug, "debug", false, "Print debug messages")
	flag.Parse()

	printDebug("Using flags [composeDirName=%v, baseCompose=%v, env=%v]", composeDirName, baseCompose, env)

	pwd, err := os.Getwd()
	pwd = filepath.ToSlash(pwd)
	if err != nil {
		handleError(err)
	}

	printDebug("Looking for compose directory")
	composeDir, err := findComposeDIr(pwd, composeDirName)
	if err != nil {
		handleError(err)
	}
	printDebug("Using directory [%s]", composeDir)

	printDebug("Looking for [%s]", baseCompose)
	baseComposePath := path.Join(composeDir, baseCompose)
	if _, err := os.Stat(baseComposePath); err != nil {
		handleError(errors.New(baseCompose + " does not exist"))
	}

	envCompose := composeFile(env)
	printDebug("Looking for [%s]", envCompose)
	envComposePath := path.Join(composeDir, envCompose)
	if _, err := os.Stat(envComposePath); err != nil {
		handleError(errors.New(envCompose + " does not exist"))
	}

	args := append([]string{"-f", baseComposePath, "-f", envComposePath}, flag.Args()...)
	cmd := exec.Command("docker-compose", args...)
	projectName := filepath.Base(filepath.Dir(composeDir)) + "-" + env
	printDebug("Using project name [%s]", projectName)
	cmd.Env = append(os.Environ(), fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", projectName))
	printDebug("Executing [%s]", strings.Join(append([]string{"docker-compose"}, args...), " "))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println("")
	cmd.Run()
}

func printDebug(s string, args ...interface{}) {
	if debug {
		fmt.Printf("debug | "+s+"\n", args...)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error | %v\n", err)
		os.Exit(1)
	}
}

func getEnv() string {
	env := os.Getenv("ENVIRONMENT")

	if env != "" {
		return env
	}

	return "devel"
}

func findComposeDIr(cwd, name string) (string, error) {
	// Assuming compose dir cannot be at root
	if cwd == "" || cwd == "." || cwd == "/" {
		return "", fmt.Errorf("cannot find the compose directory")
	}

	p := path.Join(cwd, name)
	printDebug("Trying path [%s]", p)
	_, err := os.Stat(p)
	if err != nil {
		// Keep searching up
		up, _ := path.Split(cwd)
		up = path.Clean(up)
		return findComposeDIr(up, name)
	}

	return p, nil
}

func composeFile(env string) string {
	return "docker-compose." + env + ".yml"
}
