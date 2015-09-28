package main

import (
    "flag"
    "fmt"
)

type Config struct {
    Recursive 	bool
    Path     	string
    Output		string
}

var config *Config

func init() {
    const (
        recursiveDefault = false
        recursiveDescription   = "Trawl files recursively if true, iteratively if not. Default is iterative."

        pathDefault = "/Applications"
        pathDescription   = "Root directory to begin file listing. Default is /Applications."

        outputDefault = "text"
        outputDescription   = "Accepts 3 arguments, json|yaml|text. Default is text."
    )
    config = &Config{}
    flag.BoolVar(&config.Recursive, "recursive", recursiveDefault, recursiveDescription)

    flag.StringVar(&config.Path, "path", pathDefault, pathDescription)

    flag.StringVar(&config.Output, "output", outputDefault, outputDescription)
}

func main() {
    flag.Parse()

    fmt.Println("Recursive:", config.Recursive)
    fmt.Println("Root path:", config.Path)
    fmt.Println("Output type:", config.Output)
}