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
    "io/ioutil"
    "github.com/ghodss/yaml"

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
    Size            int64       `json:"Size"`
    Name            string      `json:"Name"`
    Children        []File      `json:"Children"` 
}

var config *Config
var counter int = 0

func init() {
    const (
        recursiveDefault = true
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
        iterateFilesText(path, file)
    } else if strings.ToUpper(config.Output) == "JSON" {
        iterateFilesJSON(path, file)
    } else if strings.ToUpper(config.Output) == "YAML" {
        iterateFilesYAML(path, file)
    }
    return nil
}

//Cannot have '/' in filesystem (other than to separate directories). Or else this will get messed up
func iterateFilesText(path string, file os.FileInfo) {
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


func iterateFilesJSON(path string, file os.FileInfo) {
    // JSONFile := &File{ModifiedTime: file.ModTime(), 
    //                   IsDir: file.IsDir(),
    //                   IsLink: file.Mode()&os.ModeSymlink == os.ModeSymlink, 
    //                   Size: file.Size(), 
    //                   Name: file.Name()}

}

func iterateFilesYAML(path string, file os.FileInfo) {

}


func recurseFiles(path string, depth int) {
    var files []os.FileInfo
    files, depth = getDir(path, depth)
    if strings.ToUpper(config.Output) == "TEXT" {
        printText(files, depth, path)
    } else if strings.ToUpper(config.Output) == "JSON" {
        JSONFiles:= getJSON(files, path, depth)
        for i:=0; i<len(JSONFiles); i++ {
            jsonOutput, _ := json.MarshalIndent(JSONFiles[i], "", "     ")
            fmt.Println(string(jsonOutput))
        }     
    } else if strings.ToUpper(config.Output) == "YAML" {
        fmt.Print("YAML")
        YAMLFiles:= getJSON(files, path, depth)
        for i:=0; i<len(YAMLFiles); i++ {
            yamlOutput, _ := yaml.Marshal(YAMLFiles[i])
            fmt.Println(string(yamlOutput))
        } 
    }

}

func getDir(path string, depth int) ([]os.FileInfo, int) {
    depth++
    Files, _ := ioutil.ReadDir(path)
    return Files, depth
}

func printText(files []os.FileInfo, depth int, path string) {
    var dir bool
    for i:=0; i<len(files); i++ {
        dir = false
        for j:=0; j<depth; j++ {
            fmt.Print(" ");
        }
        fmt.Print(files[i].Name())
        if (files[i].IsDir()) {
            fmt.Print("/")
            dir = true
        }
        if files[i].Mode()&os.ModeSymlink == os.ModeSymlink {
            fmt.Print("* (symlink)")
        }
        fmt.Print("\n")
        if (dir == true) {
            files, depth := getDir(path, depth)
            printText(files,depth, path + "/" + files[i].Name())
        }
    }
}

func getJSON(files []os.FileInfo, path string, depth int) []File {
    var Children []File
    var JSONFiles []File
    for i:=0; i<len(files);i++ {
        if (files[i].IsDir()) {
            DirEntries, depth := getDir(path + "/" + files[i].Name(), depth)
            Children = getJSON(DirEntries, path + "/" + files[i].Name(), depth)
        }
        JSONFile := File{ModifiedTime: files[i].ModTime(), 
                      IsDir: files[i].IsDir(),
                      IsLink: files[i].Mode()&os.ModeSymlink == os.ModeSymlink, 
                      Size: files[i].Size(), 
                      Name: files[i].Name(),
                    Children: Children}
        JSONFiles = append(JSONFiles, JSONFile)
    }
    return JSONFiles
}

func recurseFilesJSON(path string) {}
//     JSONFile := &File{ModifiedTime: file.ModTime(), 
//                       IsDir: file.IsDir(),
//                       IsLink: file.Mode()&os.ModeSymlink == os.ModeSymlink, 
//                       Size: file.Size(), 
//                       Name: file.Name()}
// }

func recurseFilesYAML(path string){}




func main() {
    flag.Parse()

    fmt.Println("Recursive:", config.Recursive)
    fmt.Println("Root path:", config.Path)
    fmt.Println("Output type:", config.Output)

    root := config.Path
    fmt.Println(root)
    if config.Recursive == false {
        filepath.Walk(root, walkFiles)
    } else {
        recurseFiles(root, 0)
    }
    

}