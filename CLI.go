//Dan Vogel
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//Holds command line aruments
type Config struct {
	Recursive bool
	Path      string
	Output    string
}

//File fields to be converted to json/yaml
type File struct {
	ModifiedTime time.Time `json:"ModifiedTime"`
	IsLink       bool      `json:"IsLink"`
	IsDir        bool      `json:"IsDir"`
	LinksTo      string    `json:"LinksTo"`
	Size         int64     `json:"Size"`
	Name         string    `json:"Name"`
	Path         string    `json:"Path"`
	Children     []*File   `json:"Children"`
}

//A "skinnier" version of File struct for when all we have to do is print file names for the text output.
//Only used in iterateText to create stack that simulates recursive calls.
type StackInfo struct {
	File os.FileInfo
	Path string
}

var config *Config

//--define command line arguments--

//init function
func init() {
	const (
		recursiveDefault     = false
		recursiveDescription = "Trawl files recursively if true, iteratively if not."

		pathDefault     = "/"
		pathDescription = "Root directory to begin file listing"

		outputDefault     = "text"
		outputDescription = "Accepts 3 arguments, json|yaml|text."
	)
	config = &Config{}
	flag.BoolVar(&config.Recursive, "recursive", recursiveDefault, recursiveDescription)

	flag.StringVar(&config.Path, "path", pathDefault, pathDescription)

	flag.StringVar(&config.Output, "output", outputDefault, outputDescription)
}

//--iterative--

//Iteratively prints out directory structure in plain text format
func iterateText(path string) {
	rootFileInfo, _ := os.Stat(path)
	rootFile := StackInfo{rootFileInfo, path}
	stack := []StackInfo{rootFile}
	basedepth := len(strings.Split(path, "/"))

	for len(stack) > 0 {
		file := stack[len(stack)-1].File
		path = stack[len(stack)-1].Path
		stack = stack[:len(stack)-1]
		for i := 0; i < len(strings.Split(path, "/"))-basedepth; i++ {
			fmt.Print(" ")
		}
		fmt.Print(file.Name())
		if file.IsDir() {
			fmt.Print("/")
			children, _ := ioutil.ReadDir(path)
			for i := 0; i < len(children); i++ {
				stack = append(stack, StackInfo{children[i], filepath.Join(path, children[i].Name())})
			}
		}
		if file.Mode()&os.ModeSymlink == os.ModeSymlink {
			target, _ := filepath.EvalSymlinks(path)
			fmt.Print("* (symlink)    target: '" + target + "'")
		}
		fmt.Println()
	}

}

//Iteratively constructs a tree of File structs, preserving the directory structure.
func iterateJSON(path string) *File {
	rootOSFile, _ := os.Stat(path)
	rootFile := toFile(rootOSFile, path, []*File{})
	stack := []*File{rootFile}

	for len(stack) > 0 {
		file := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		children, _ := ioutil.ReadDir(file.Path)
		for _, child := range children {
			curChild := toFile(child, filepath.Join(file.Path, child.Name()), []*File{})
			file.Children = append(file.Children, curChild)
			stack = append(stack, curChild)
		}
	}
	return rootFile
}

//--recursion--

//Recursively prints out directory in plain text
func recurseText(path string, depth int) {
	var dir bool
	depth++
	Files, _ := ioutil.ReadDir(path)
	for i := 0; i < len(Files); i++ {
		dir = false
		for j := 0; j < depth; j++ {
			fmt.Print(" ")
		}
		fmt.Print(Files[i].Name())
		if Files[i].IsDir() {
			fmt.Print("/")
			dir = true
		}
		if Files[i].Mode()&os.ModeSymlink == os.ModeSymlink {
			target, _ := filepath.EvalSymlinks(path + "/" + Files[i].Name())
			fmt.Print("* (symlink)    target: '" + target + "'")
		}
		fmt.Print("\n")
		if dir == true {
			recurseText(path+"/"+Files[i].Name(), depth)
		}
	}
}

//Recusively creates a tree of File structs, preserving directory structure.
func recurseJSON(path string) []*File {
	var children []*File
	var JSONFiles []*File
	files, _ := ioutil.ReadDir(path)
	for i := 0; i < len(files); i++ {
		if files[i].IsDir() {
			children = recurseJSON(filepath.Join(path, files[i].Name()))
		}
		JSONFile := toFile(files[i], filepath.Join(path, files[i].Name()), children)
		JSONFiles = append(JSONFiles, JSONFile)
	}

	return JSONFiles
}

//--helper function--
//Uses the path and os.FileInfo of a file to construct a File object from it
func toFile(file os.FileInfo, path string, children []*File) *File {

	JSONFile := File{ModifiedTime: file.ModTime(),
		IsDir:    file.IsDir(),
		Size:     file.Size(),
		Name:     file.Name(),
		Path:     path,
		Children: children}
	if file.Mode()&os.ModeSymlink == os.ModeSymlink {
		JSONFile.IsLink = true
		JSONFile.LinksTo, _ = filepath.EvalSymlinks(filepath.Join(path, file.Name()))
	}
	return &JSONFile
}

//--handler--
//Handles the arguments passed from command line.
func handler(path string) {
	var JSONFile *File
	var output []uint8
	if strings.ToUpper(config.Output) == "TEXT" {
		if config.Recursive == true {
			recurseText(path, 0)
		} else if config.Recursive == false {
			iterateText(path)
		}
	} else if strings.ToUpper(config.Output) == "JSON" || strings.ToUpper(config.Output) == "YAML" {
		if config.Recursive == true {
			rootInfo, _ := os.Stat(path)
			JSONFile = toFile(rootInfo, path, recurseJSON(path))
		} else if config.Recursive == false {
			JSONFile = iterateJSON(path)
		}
		if strings.ToUpper(config.Output) == "JSON" {
			output, _ = json.MarshalIndent(JSONFile, "", "     ")
			fmt.Println(string(output))

		} else if strings.ToUpper(config.Output) == "YAML" {
			output, _ = yaml.Marshal(JSONFile)
			fmt.Println(string(output))
		}

	}
}

//main
func main() {
	flag.Parse()

	fmt.Println("Recursive:", config.Recursive)
	fmt.Println("Root path:", config.Path)
	fmt.Println("Output type:", config.Output)

	root := config.Path
	fmt.Println("\nStarting file tree output:\n")
	handler(root)

}
