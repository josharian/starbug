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
	if len(os.Args) != 2 {
		log.Fatal("usage: starbug <prog>")
	}
	prog := os.Args[1]
	for _, impl := range impls {
		path, err := exec.LookPath(impl.bin)
		if err != nil {
			log.Printf("%s not found: %v", err)
		}
		cmd := exec.Command(path, "-c", prog)
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
		log.Printf("$ %s -c %s%s%s", impl.name, quote, prog, quote)
		log.Print(string(out))
		if rc != -1 {
			log.Printf("(exit %v)", rc)
		}
		log.Println()
	}
}
