package fs

import (
    "io/ioutil"
)

func ListFilesInDir(filePath string) ([]string, error) {
	files, err := ioutil.ReadDir(filePath)
    if err != nil {
        return nil, err;
    }

	fileList := []string{};
    for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name());
		}
    }

	return fileList, nil
}

