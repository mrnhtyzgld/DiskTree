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
	struct1, err := StartOffRoot("./")
	root = append(root, &struct1)
	root, err = root, err

	file, _ = RedirectTo("./", "logs1", "txt")
	struct2, err := StartOffRoot("./../")
	root = append(root, &struct2)

	ResetOutput(oldOut, file)

	fmt.Println("scanning complete")
	fmt.Println(findFullPath(root[1].Childs[0].Childs[1]))

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
		TreeSize: info.Size(), // önce kendi sizeına set edilir sonra size bulma metodu çağırılacak
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

func StartOffRoot(path string) (TreeFile, error) {
	info, err := os.Stat(path)
	if err != nil {
		return TreeFile{}, err
	}

	startingOff := TreeFile{
		Name:     info.Name(),
		FileSize: info.Size(),
		isDir:    info.IsDir(),
		Childs:   make([]*TreeFile, 0),
		Parent:   nil,
		TreeSize: info.Size(), // önce kendi sizeına set edilir sonra size bulma metodu çağırılacak
	}

	files, err := os.ReadDir(path)
	for _, file := range files {
		err = Iterate(&startingOff.Childs, filepath.Join(path, file.Name()), &startingOff)
	}
	if err != nil {
		return TreeFile{}, err
	}

	return startingOff, nil
}

func findTreeSize(file *TreeFile) {
	for _, ffile := range file.Childs {
		findTreeSize(ffile)
	}
	for _, ffile := range file.Childs {
		file.TreeSize += ffile.TreeSize
	}
}

func findFullPath(file *TreeFile) string {
	if file.Parent == nil {
		return ""
	}
	return "" + findFullPath(file.Parent) + "/" + file.Name
}

// FIXME i need to learn more about callback funcs
// FIXME these aint gonna work
// FIXME also not tested

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
