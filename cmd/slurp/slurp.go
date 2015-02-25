package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	gopaths = strings.Split(os.Getenv("GOPATH"), ":")
	cwd     string

	build   = flag.Bool("build", false, "build the current build as slurp-bin")
	install = flag.Bool("install", false, "install current slurp.Go as slurp.PKG.")

	keep = flag.Bool("keep", false, "keep the generated source under $GOPATH/src/slurp-run-*")
)

func main() {

	flag.Parse()

	if len(gopaths) == 0 || gopaths[0] == "" {
		log.Fatal("$GOPATH must be set.")
	}

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	path, err := generate()
	if err != nil {
		return err
	}

	//Don't forget to clean up.
	if !*keep {
		defer os.RemoveAll(path)
	}

	get := exec.Command("go", "get", "-tags=slurp", "-v")
	get.Dir = path
	get.Stdin = os.Stdin
	get.Stdout = os.Stdout
	get.Stderr = os.Stderr

	if *build || *install {
		err := get.Run()
		if err != nil {
			return err
		}
	}

	var args []string

	if *build {
		args = []string{"build", "-tags=slurp", "-o=slurp-bin", path}

	} else if *install {
		args = []string{"install", "-tags=slurp", path}

	} else {
		params := flag.Args()

		if len(params) > 0 && params[0] == "init" {
			err := get.Run()
			if err != nil {
				return err
			}
		}

		args = []string{"run", "-tags=slurp", filepath.Join(path, "main.go")}
		args = append(args, params...)
	}

	cmd := exec.Command("go", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func generate() (string, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// find the correct gopath
	var gopathsrc string
	var pkgpath string
	for _, gopath := range gopaths {
		gopathsrcTest := filepath.Join(gopath, "src")
		// the target package import path.
		pkgpath, err = filepath.Rel(gopathsrcTest, cwd)
		if err != nil {
			return "", err
		}
		if base := filepath.Base(pkgpath); base == "." || base == ".." {
			continue // cwd is outside this gopath
		}
		gopathsrc = gopathsrcTest
	}

	if gopathsrc == "" {
		return "", errors.New("forbidden path. Your CWD must be under $GOPATH/src.")
	}

	//build our package path.
	path := filepath.Join(gopathsrc, "slurp", filepath.Dir(pkgpath), "slurp-"+filepath.Base(pkgpath))

	//Clean it up.
	// os.RemoveAll(path)

	err = os.MkdirAll(path, 0700)
	if err != nil {
		return path, err
	}

	//log.Println("Generating the runner...")
	file, err := os.Create(filepath.Join(path, "main.go"))
	if err != nil {
		return path, err
	}

	pkg, err := filepath.Rel(gopathsrc, cwd)
	if err != nil {
		return path, err
	}

	err = runnerSrc.Execute(file, filepath.ToSlash(pkg))
	if err != nil {
		return path, err
	}

	err = file.Close()
	if err != nil {
		return path, err
	}
	return path, nil
}
