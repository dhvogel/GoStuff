//Dan Vogel
//Must include link handler
package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "encoding/json"
    "time"
    //"errors"
)

type Config struct {
    Recursive 	bool
    Path     	string
    Output		string
}

type File struct {
    ModifiedTime    time.Time   `json:"ModifiedTime"`
    IsLink          bool        `json:"IsLink"`
    IsDir           bool        `json:"IsDir"`
    //LinksTo         bool        `json:"LinksTo"`
    Size            int64         `json:"Size"`
    Name            string      `json:"Name"`
    //Children        []File      `json:"File"` 
}

var config *Config
var counter int = 0

func init() {
    const (
        recursiveDefault = false
        recursiveDescription   = "Trawl files recursively if true, iteratively if not. Default is iterative."

        pathDefault = "/Users/dhvogel/Documents/CS341"
        pathDescription   = "Root directory to begin file listing. Default is /Users/dhvogel/Documents/CS341."

        outputDefault = "json"
        outputDescription   = "Accepts 3 arguments, json|yaml|text. Default is text."
    )
    config = &Config{}
    flag.BoolVar(&config.Recursive, "recursive", recursiveDefault, recursiveDescription)

    flag.StringVar(&config.Path, "path", pathDefault, pathDescription)

    flag.StringVar(&config.Output, "output", outputDefault, outputDescription)
}

func walkFiles(path string, file os.FileInfo, err error) error {
    if strings.ToUpper(config.Output) == "TEXT" {
        printTextFile(path, file)
    } else if strings.ToUpper(config.Output) == "JSON" {
        printJSONFile(path, file)
    } else if strings.ToUpper(config.Output) == "YAML" {
        printYAMLFile(path, file)
    }
    return nil
}

//Cannot have '/' in filesystem (other than to separate directories). Or else this will get messed up
func printTextFile(path string, file os.FileInfo) {
    for i:=0; i<len(strings.Split(path, "/")); i++ {
        fmt.Print(" ")
    }
    // if (file.Name() == nil) {
    //     err := errors.New("Please enter valid root directory")
    //     fmt.Println(err)
    // }
    fmt.Printf("%s", file.Name())
    if file.IsDir() {
        fmt.Print("/")
    }
    if file.Mode()&os.ModeSymlink == os.ModeSymlink {
        fmt.Print("* (symlink)")
    }
    fmt.Print("\n")
}

func printJSONFile(path string, file os.FileInfo) {
    JSONFile := &File{ModifiedTime: file.ModTime(), 
                      IsDir: file.IsDir(),
                      IsLink: file.Mode()&os.ModeSymlink == os.ModeSymlink, 
                      Size: file.Size(), 
                      Name: file.Name()}
    jsonOutput, _ := json.Marshal(JSONFile)
    fmt.Println(string(jsonOutput))
}

func printYAMLFile(path string, file os.FileInfo) {

}



func main() {
    flag.Parse()

    fmt.Println("Recursive:", config.Recursive)
    fmt.Println("Root path:", config.Path)
    fmt.Println("Output type:", config.Output)

    root := config.Path
    fmt.Println(root)
    filepath.Walk(root, walkFiles)

}