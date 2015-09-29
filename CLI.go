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
    "gopkg.in/yaml.v2"

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
    LinksTo         string        `json:"LinksTo"`
    Size            int64       `json:"Size"`
    Name            string      `json:"Name"`
    Children        []File      `json:"Children"` 
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

}

func iterateFilesYAML(path string, file os.FileInfo) {

}


func recurseFiles(path string) {
    var files []os.FileInfo
    files, _ = ioutil.ReadDir(path)
    if strings.ToUpper(config.Output) == "TEXT" {
        recurseText(path,0)
    } else if strings.ToUpper(config.Output) == "JSON" {
        JSONFiles:= recurseJSON(files, path)
        for i:=0; i<len(JSONFiles); i++ {
            jsonOutput, _ := json.MarshalIndent(JSONFiles[i], "", "     ")
            fmt.Println(string(jsonOutput))
        }     
    } else if strings.ToUpper(config.Output) == "YAML" {
        YAMLFiles := recurseJSON(files, path)
        for i:=0; i<len(YAMLFiles); i++ {
            yamlOutput, _ := yaml.Marshal(YAMLFiles[i])
            fmt.Println(string(yamlOutput))
        } 
    }

}



func recurseText(path string, depth int) {
    var dir bool
    depth++
    Files, _ := ioutil.ReadDir(path)
    for i:=0; i<len(Files); i++ {
         dir = false
         for j:=0; j<depth; j++ {
             fmt.Print(" ");
         }
         fmt.Print(Files[i].Name())
         if (Files[i].IsDir()) {
             fmt.Print("/")
             dir = true
         }
         if Files[i].Mode()&os.ModeSymlink == os.ModeSymlink {
            target, _ := filepath.EvalSymlinks(path + "/" + Files[i].Name())
            fmt.Print("* (symlink)    target: '" + target + "'")
         }
         fmt.Print("\n")
         if (dir == true) {  
             recurseText(path + "/" + Files[i].Name(), depth)
         }
     }
}

func recurseJSON(files []os.FileInfo, path string) []File {
    var Children []File
    var JSONFiles []File
    var isLink bool
    var linksTo string
    for i:=0; i<len(files);i++ {
        if (files[i].IsDir()) {
            DirEntries, _ := ioutil.ReadDir(path + "/" + files[i].Name())
            Children = recurseJSON(DirEntries, path + "/" + files[i].Name())
        }
        if (files[i].Mode()&os.ModeSymlink == os.ModeSymlink) {
            isLink = true
            linksTo, _ = filepath.EvalSymlinks(path + "/" + files[i].Name())
        } else {
            isLink = false
            linksTo = ""
        }
        JSONFile := File{ModifiedTime: files[i].ModTime(), 
                      IsDir: files[i].IsDir(),
                      IsLink: isLink,
                      LinksTo: linksTo,
                      Size: files[i].Size(), 
                      Name: files[i].Name(),
                    Children: Children}
        JSONFiles = append(JSONFiles, JSONFile)
    }
    return JSONFiles
}




func main() {
    flag.Parse()

    fmt.Println("Recursive:", config.Recursive)
    fmt.Println("Root path:", config.Path)
    fmt.Println("Output type:", config.Output)

    root := config.Path
    fmt.Println("\nStarting file tree output:\n")
    if config.Recursive == false {
        filepath.Walk(root, walkFiles)
    } else {
        recurseFiles(root)
    }
    

}