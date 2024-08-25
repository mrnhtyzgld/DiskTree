package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type TreeFile struct {
	Parent   *TreeFile
	Name     string
	TreeSize int64
	FileSize int64
	isDir    bool
	Childs   []*TreeFile
}

func main() {
	file, oldOut := RedirectTo("./", "logs", "txt")

	root := make([]*TreeFile, 0)
	err := StartOffRoot(&root, "./")
	root, err = root, err

	file, _ = RedirectTo("./", "logs1", "txt")
	err = StartOffRoot(&root, "./../")

	ResetOutput(oldOut, file)

	fmt.Println("scanning complete")

}

func RedirectTo(path string, name string, fileType string) (*os.File, *os.File) {

	oldOut := os.Stdout
	file, err := os.Create(path + name + "." + fileType)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil, nil
	}

	os.Stdout = file

	return file, oldOut
}

func ResetOutput(out *os.File, file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Println("Error closing file:", err)
	}
	os.Stdout = out
}

func Iterate(Childs *[]*TreeFile, path string, Parent *TreeFile) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	*Childs = append(*Childs, &TreeFile{
		Name:     info.Name(),
		FileSize: info.Size(),
		isDir:    info.IsDir(),
		Childs:   make([]*TreeFile, 0),
		Parent:   Parent,
	})

	index := len(*Childs) - 1
	me := (*Childs)[index]

	files, err := os.ReadDir(path)
	for _, file := range files {
		Iterate(&me.Childs, filepath.Join(path, file.Name()), me)
	}

	fmt.Println("succesfull: \t" + path)
	return nil
}

func StartOffRoot(rootParam *[]*TreeFile, path string) error {
	return Iterate(rootParam, path, nil)
}

func findTreeSize(file *TreeFile) int64 {
	return 0
}

func findFullPath(file *TreeFile) string {
	return ""
}

// FIXME i need to learn more about callback funcs
// FIXME these aint gonna work

func fromRootToTreeFileIter(folder *TreeFile, path string, fn func(data ...interface{})) error {
	lastIndex := len(path) - len(folder.Name) - 1
	err := fromRootToTreeFileIter(folder.Parent, path[0:lastIndex], fn)
	fn()
	if err != nil {
		return err
	}
	return nil
}

func postOrderIter(folder *TreeFile, path string, fn func(...interface{})) error {

	err := *new(error)
	for _, file := range (*folder).Childs {
		err = postOrderIter(file, filepath.Join(path, file.Name), fn)
	}
	fn()
	if err != nil {
		return err
	}
	return nil
}
func preOrderIter(folder *TreeFile, path string, fn func(...interface{})) error {

	fn()

	err := *new(error)
	for _, file := range (*folder).Childs {
		err = preOrderIter(file, filepath.Join(path, file.Name), fn)
	}
	if err != nil {
		return err
	}
	return nil
}
