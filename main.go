package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

var impls = [...]struct {
	bin  string // binary to execute
	name string // nice name for the implementation
}{
	{"python2", "python2"},
	{"python3", "python3"},
	{"starlark", "starlark-go"},
	{"starlark-repl", "starlark-rust"},
}

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatal("usage: starbug [flags] <prog>")
	}

	// Attempt to auto-assign flags to the appropriate interpreter.
	// For now, simply recognize and accept starlark-go flags.
	flags := os.Args[1 : len(os.Args)-1]
	validgoflag := map[string]bool{
		"-bitwise":        true,
		"-float":          true,
		"-globalreassign": true,
		"-lambda":         true,
		"-nesteddef":      true,
		"-recursion":      true,
		"-set":            true,
		"-showenv":        true,
	}
	allflags := make(map[string][]string)
	for _, s := range flags {
		if validgoflag[s] {
			allflags["starlark-go"] = append(allflags["starlark-go"], s)
		} else {
			log.Fatal("unrecognized flag", s)
		}
	}

	prog := os.Args[len(os.Args)-1]

	for _, impl := range impls {
		path, err := exec.LookPath(impl.bin)
		if err != nil {
			log.Printf("%s not found: %v", err)
		}
		iflags := allflags[impl.name]
		iflags = iflags[:len(iflags):len(iflags)]
		args := append(iflags, "-c", prog)
		cmd := exec.Command(path, args...)
		out, err := cmd.CombinedOutput()
		rc := -1
		if err == nil {
			rc = 0
		} else if ee, ok := err.(*exec.ExitError); ok {
			rc = ee.ExitCode()
		}
		quote := `'`
		if strings.Contains(prog, `'`) {
			quote = `"`
			if strings.Contains(prog, `"`) {
				// give up, rely on user to understand
				quote = ""
			}
		}
		flagstr := ""
		for _, f := range iflags {
			flagstr += " " + f
		}
		log.Printf("$ %s%s -c %s%s%s", impl.name, flagstr, quote, prog, quote)
		log.Print(string(out))
		if rc != -1 {
			log.Printf("(exit %v)", rc)
		}
		log.Println()
	}
}
