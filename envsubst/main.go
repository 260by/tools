package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/drone/envsubst"
)

func main() {
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprint(os.Stderr, "no input found\n")
		os.Exit(-1)
	}

	s := string(in)

	out, err := envsubst.EvalEnv(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "replace error: %s\n", err.Error())
		os.Exit(-1)
	}

	n, err := fmt.Fprint(os.Stdout, out)
	if n != len(out) {
		fmt.Fprint(os.Stderr, "write stdout crashed\n")
		os.Exit(-1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "write stdout error: %s\n", err.Error())
		os.Exit(-1)
	}

	os.Exit(0)
}
