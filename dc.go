package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var composeDirName, baseCompose, env string
var verbose bool

var flags = [...]string{"-composeDirName", "-baseCompose", "-env", "-v"}

func main() {
	flag.StringVar(&composeDirName, "composeDirName", ".compose", "Name of directory containing docker compose files")
	flag.StringVar(&baseCompose, "baseCompose", "docker-compose.default.yml", "The base docker compose file")
	flag.StringVar(&env, "env", "devel", "Environment that docker compose is running in")
	flag.BoolVar(&verbose, "v", false, "Environment that docker compose is running in")
	flag.Parse()

	pwd, err := os.Getwd()
	pwd = filepath.ToSlash(pwd)
	if err != nil {
		handleError(err)
	}

	composeDir, err := findComposeDIr(pwd, composeDirName)
	if err != nil {
		handleError(err)
	}
	info("Use compose directory %s", composeDir)

	debug("Looking for %s", baseCompose)
	baseComposePath := path.Join(composeDir, baseCompose)
	if _, err := os.Stat(baseComposePath); err != nil {
		handleError(errors.New(baseCompose + " does not exist"))
	}

	envCompose := composeFile(env)
	debug("Looking for %s", envCompose)
	envComposePath := path.Join(composeDir, envCompose)
	if _, err := os.Stat(envComposePath); err != nil {
		handleError(errors.New(envCompose + " does not exist"))
	}

	args := []string{"-f", baseComposePath, "-f", envComposePath}
	args = append(args, getArgs()...)
	info("Executing: %s", strings.Join(append([]string{"docker-compose"}, args...), " "))
	cmd := exec.Command("docker-compose", args...)
	projectName := filepath.Base(filepath.Dir(composeDir)) + "-" + env
	info("COMPOSE_PROJECT_NAME=%s", projectName)
	cmd.Env = append(os.Environ(), fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", projectName))

	stdout, err := cmd.StdoutPipe()
	handleError(err)
	scannerStd := bufio.NewScanner(stdout)
	go func() {
		for scannerStd.Scan() {
			fmt.Println(scannerStd.Text())
		}
	}()

	stderr, err := cmd.StderrPipe()
	handleError(err)
	scannerErr := bufio.NewScanner(stderr)
	go func() {
		for scannerErr.Scan() {
			fmt.Println(scannerErr.Text())
		}
	}()

	fmt.Println("")
	_ = cmd.Start()
	// handleError(err)

	_ = cmd.Wait()
	// handleError(err)
}

func info(s string, args ...interface{}) {
	fmt.Printf("[INFO] "+s+"\n", args...)
}

func debug(s string, args ...interface{}) {
	if verbose {
		fmt.Printf("[DEBUG] "+s+"\n", args...)
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

func getArgs() []string {
	args := make([]string, 0)
	values := os.Args[1:]

	for _, v := range values {
		// remove flags
		match, _ := regexp.MatchString("^-", v)
		if !match || !isOwnFlag(v) {
			args = append(args, v)
		}
	}

	return args
}

func isOwnFlag(f string) bool {
	for _, v := range flags {
		if v == f {
			return true
		}
	}

	return false
}

func findComposeDIr(cwd, name string) (string, error) {
	// Assuming compose dir cannot be at root
	if cwd == "" {
		return "", fmt.Errorf("cannot find the compose directory")
	}

	p := path.Join(cwd, name)
	debug("Trying path %s", p)
	_, err := os.Stat(p)
	if err != nil {
		// Keep searching up
		up, _ := path.Split(cwd)
		return findComposeDIr(up, name)
	}

	return p, nil
}

func composeFile(env string) string {
	return "docker-compose." + env + ".yml"
}
