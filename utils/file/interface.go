package file

import (
	"io/ioutil"
	"os"
	"path"
)

func CreateAndWrite(filename string, content []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(content)
	return nil
}

func ReadAll(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

func FullName(filename string) string {
	if path.IsAbs(filename) {
		return filename
	}

	dir, _ := os.Getwd()
	return dir + "/" + filename
}

/*
func ListDir(dir string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dirs {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dir+PthSep+fi.Name())
		}
	}
	return files, nil
}
*/
