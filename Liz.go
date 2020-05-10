package main

import (
	"flag"
	"fmt"
	"os"

	"Liz/builder"
	"Liz/kernel/container"
	"Liz/kernel/services"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	newName := newCmd.String("name", "App", "type new project name")
	newPath := newCmd.String("path", "./", "path to new project")

	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildPath := buildCmd.String("path", "./", "path to project to build")

	newCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [new|build] [flags...]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n-----------------------------------------\n\nflags for new\n")
		newCmd.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n-----------------------------------------\n\nflags for build\n")
		buildCmd.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n-----------------------------------------\n")
	}
	buildCmd.Usage = newCmd.Usage
	flag.Usage = newCmd.Usage
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	services.Build()

	switch os.Args[1] {
	case "new":
		check(newCmd.Parse(os.Args[2:]))
		makeDir(*newPath)
		check(os.Chdir(*newPath))
		container.Get("project_starter_builder").(*builder.ProjectStarter).Build(*newName)
	case "build":
		check(buildCmd.Parse(os.Args[2:]))
		makeDir(*newPath)
		check(os.Chdir(*buildPath))
		container.Get("container_builder").(*builder.Container).Build()
		break
	default:
		flag.Usage()
		os.Exit(2)
	}
}

func makeDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		check(os.MkdirAll(path, os.ModePerm))
	}
}

// func absPath(filename string) (string, error) {
// 	eval, err := filepath.EvalSymlinks(filename)
// 	if err != nil {
// 		return "", err
// 	}
// 	return filepath.Abs(eval)
// }

// func fill() {
// 	path, err := absPath("./")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	cfg := &packages.Config{
// 		Mode:  packages.LoadAllSyntax,
// 		Tests: true,
// 		Dir:   filepath.Dir(path),
// 		Fset:  token.NewFileSet(),
// 		Env:   os.Environ(),
// 	}

// 	pkgs, err := packages.Load(cfg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, pkg := range pkgs {
// 		fmt.Printf("%v\n", pkg)
// 	}
// }
