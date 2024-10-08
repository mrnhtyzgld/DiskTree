package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// FIXME hata çözme-paniik
// TODO add goroutine support

type TreeFile struct {
	Parent   *TreeFile
	Name     string
	TreeSize int64
	FileSize int64
	isDir    bool
	Childs   []*TreeFile
}

var logState = false
var errState = true

var infoState = true

func main() {
	file, oldOut := RedirectTo("./", "logs", "txt")

	root := make([]*TreeFile, 0)
	struct1, err := StartOffRoot("\"C:\\Users\\NihatEmreYüzügüldü\\Desktop\\100MEDİA\"")
	root = append(root, &struct1)
	if err != nil {
		printErr("DiskTree okuması başarısız oldu")
	}

	/*
		file, _ = RedirectTo("./", "logs1", "txt")
		struct2, err := StartOffRoot("C:\\Users\\NihatEmreYüzügüldü\\GolandProjects\\awesomeProject")
		root = append(root, &struct2)
		if err != nil {
			printErr("DiskTree okuması başarısız oldu")
		}
	*/

	ResetOutput(oldOut, file)

	printInfo("scanning complete")
	printInfo(struct1.TreeSize, struct1.Childs[0].TreeSize)
}

func RedirectTo(path string, name string, fileType string) (*os.File, *os.File) {

	oldOut := os.Stdout
	file, err := os.Create(path + name + "." + fileType)
	if err != nil {
		printErr("Error creating file:", err)
		return nil, nil
	}

	os.Stdout = file

	return file, oldOut
}

func ResetOutput(out *os.File, file *os.File) {
	if err := file.Close(); err != nil {
		printErr("Error closing file:", err)
	}
	os.Stdout = out
}

func Iterate(Childs *[]*TreeFile, path string, Parent *TreeFile) error {
	info, err := os.Stat(path)
	if err != nil {
		printErr("başarısız: " + path + " okunamadı")
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
		err = Iterate(&me.Childs, filepath.Join(path, file.Name()), me)
	}

	if err != nil {
		return err
	}
	printLog("succesfull: \t" + path)
	return nil
}

func StartOffRoot(path string) (TreeFile, error) {
	info, err := os.Stat(path)
	if err != nil {
		printErr("root dizini okumada problem çıktı")
		return TreeFile{}, err
	}

	startingOff := TreeFile{
		Name:     info.Name(),
		FileSize: info.Size(),
		isDir:    info.IsDir(),
		Childs:   make([]*TreeFile, 0),
		Parent:   nil,
		TreeSize: info.Size(), // önce kendi sizeına set edilir sonra recursive size bulma metodu çağırılacak
	}

	files, err := os.ReadDir(path)
	for _, file := range files {
		err = Iterate(&startingOff.Childs, filepath.Join(path, file.Name()), &startingOff)
	}
	findTreeSize(&startingOff)
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

func printLog(a ...any) {
	if logState {
		fmt.Println(a...)
	}
}
func printErr(a ...any) {
	if errState {
		fmt.Println(a...)
	}
}
func printInfo(a ...any) {
	if infoState {
		fmt.Println(a...)
	}
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
