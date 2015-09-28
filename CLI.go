package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

type Config struct {
    Recursive 	bool
    Path     	string
    Output		string
}

var config *Config
var counter int = 0

func init() {
    const (
        recursiveDefault = false
        recursiveDescription   = "Trawl files recursively if true, iteratively if not. Default is iterative."

        pathDefault = "/Users/dhvogel/Documents/CS341"
        pathDescription   = "Root directory to begin file listing. Default is /Users/dhvogel/Documents/CS341."

        outputDefault = "text"
        outputDescription   = "Accepts 3 arguments, json|yaml|text. Default is text."
    )
    config = &Config{}
    flag.BoolVar(&config.Recursive, "recursive", recursiveDefault, recursiveDescription)

    flag.StringVar(&config.Path, "path", pathDefault, pathDescription)

    flag.StringVar(&config.Output, "output", outputDefault, outputDescription)
}

func walkFiles(path string, file os.FileInfo, err error) error {
    printFile(path, file, config.Output)
    return nil
}

//Cannot have '/' in filesystem (other than to separate directories). Or else this will get messed up
func printFile(path string, file os.FileInfo, mode string) {
    for i:=0; i<len(strings.Split(path, "/")); i++ {
        fmt.Print(" ")
    }
    fmt.Printf("%s", file.Name())
    if file.IsDir() {
        fmt.Print("/")
    }
    fmt.Print("\n")
}



func main() {
    flag.Parse()

    fmt.Println("Recursive:", config.Recursive)
    fmt.Println("Root path:", config.Path)
    fmt.Println("Output type:", config.Output)

    root := config.Path
    err := filepath.Walk(root, walkFiles)
    fmt.Printf("filepath.Walk() returned %v\n", err)

}