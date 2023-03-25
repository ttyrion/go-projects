package main

import (
	"fmt"
	"os"
	"encoding/json"
	"strings"
	"github.com/parsemd/parsemdpkg"
	"github.com/parsemd/fs"
)

// 字段名必须是大写开头，才能被转换为json
type PostInfo struct {
	FileName string `json:"fileName"`;
	PreviewContent string `json:"previewContent"`;
}

type PostItemForVue struct {
	Post PostInfo `json:"post"`;
	Header parsemdpkg.MarkDownHeader `json:"header"`;
}

func findAPreviewContent(parsedMarkDownData parsemdpkg.ParsedMarkDownData) (previewContent string) {
	if parsedMarkDownData.Prologue != "" {
		return parsedMarkDownData.Prologue;
	}

	if len(parsedMarkDownData.MainData) <= 0 {
		return "";
	}

	for _, data := range parsedMarkDownData.MainData {
		if len(data.Data) <= 0 {
			break;
		}

		// 从data里面找出一段文字当预览
		for _, baseDta := range data.Data {
			if len(data.Data) <= 0 {
				break;
			}
	
			if baseDta.Text != "" {
				previewContent = baseDta.Text;
				break;
			}
		}

		// 已经找到一段文字
		if previewContent != "" {
			break;
		}
	}

	return previewContent;
}

func main() {
	postDir := "./post";

	fileList , err := fs.ListFilesInDir(postDir);
	if err != nil {
		fmt.Printf("could not read dir: %s\n", err)
		return
	}

	mainPostForVue := []PostItemForVue{};

	for _, fileName := range fileList {
		fmt.Println(fileName);
		parsedMarkDownData := parsemdpkg.MarkdownFile2ParsedMarkDownData(postDir + "/" + fileName)
		jsonData, err := json.Marshal(parsedMarkDownData)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		// fmt.Println(string(jsonData));

		pureFileName,_ := strings.CutSuffix(fileName, ".md");
		jsonFileName := pureFileName + ".json"

		err = os.WriteFile("./json/" + jsonFileName, jsonData, 0666);
		if err != nil {
			fmt.Println("could not WriteFile:" + jsonFileName, err);
		} else {
			postInfo := PostInfo{
				FileName: jsonFileName,
				PreviewContent: findAPreviewContent(parsedMarkDownData),
			}

			postItem := PostItemForVue {
				Post: postInfo,
				Header: parsedMarkDownData.Header,
			}

			mainPostForVue = append(mainPostForVue, postItem);
		}
	}

	// main.json
	n := len(mainPostForVue);
	if n > 0 {
		// 反序，日期最新的博客文章排在前面
		reversedMainPostForVue := []PostItemForVue{};
		for i := n-1; i >= 0; i-- {
			reversedMainPostForVue = append(reversedMainPostForVue, mainPostForVue[i])
		}

		jsonData, err := json.Marshal(reversedMainPostForVue)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		if err := os.WriteFile("./json/main.json", jsonData, 0666); err != nil {
			fmt.Println("could not WriteFile main.json", err);
		}
	}

	return 
}
