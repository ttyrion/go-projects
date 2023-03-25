package parsemdpkg

import (
	"bufio"
    "fmt"
    "os"
	"strings"
	"regexp"
	"strconv"
)

type MarkDownHeader struct {
	Title string `json:"title"`;
	Date string `json:"date"`;
	Author string `json:"author"`;
	HeaderImage string `json:"headerImage"`;
}

type CodeBlock []string;

type BaseData struct {
	Text string `json:"text"`;
	ImageUrl string `json:"imageUrl"`;
	Code string `json:"code"`;
	Textlist []string `json:"textList"`;
}

type MarkDownTopic struct {
	H int `json:"h"`;  // h3, h4, h5
	Title string `json:"title"`;
	Data []BaseData `json:"data"`;
}

type ParsedMarkDownData struct {
	Header MarkDownHeader  `json:"header"`;
	Prologue string `json:"prologue"`;
	MainData []MarkDownTopic `json:"mainData"`;
}

func MarkdownFile2ParsedMarkDownData(filePath string) ParsedMarkDownData {
	result := ParsedMarkDownData{}

	file, err := os.Open(filePath)
    if err != nil {
		fmt.Println("open " + filePath + " failed.");
        return result;
    }

	scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    var textLines []string
    for scanner.Scan() {
        textLines = append(textLines, scanner.Text())
    }
    file.Close()
  

	// 0 未开始
	// 1 正在解析内部内容
	// 2 内部内容解析结束
	mdHeaderBegin := 0;
	codeBlockBegin := 0;
	
	// textListBegin := 0;

	markDownTopicSlice := []MarkDownTopic{};
	markDownTopic := MarkDownTopic{};
	codeBlock := CodeBlock{};
	header := MarkDownHeader{};
	textListData := BaseData{};
    for _, line := range textLines {
		lineType, content := parseLine(line);

		fmt.Println("MarkdownFile2ParsedMarkDownData: type=" + strconv.Itoa(lineType) + " codeBlockBegin=" + strconv.Itoa(codeBlockBegin) + " " + content);

		if lineType == 1 {
			// textListBegin = 0;
			if len(textListData.Textlist) > 0 {
				markDownTopic.Data = append(markDownTopic.Data, textListData);
				textListData = BaseData{};
			}

			if mdHeaderBegin == 0 {
				mdHeaderBegin = 1;
			} else {
				mdHeaderBegin = 2
			}
		} else if lineType == 2 {
			// textListBegin = 0;
			if len(textListData.Textlist) > 0 {
				markDownTopic.Data = append(markDownTopic.Data, textListData);
				textListData = BaseData{};
			}

			if codeBlockBegin == 0 || codeBlockBegin == 2 {
				codeBlockBegin = 1;
			} else if codeBlockBegin == 1 {
				codeBlockBegin = 2;
				code := "";
				if len(codeBlock) > 0 {
					for _, codeLine := range codeBlock {
						code += codeLine;
					}
				}

				codeData := BaseData{
					Text: "",
					ImageUrl: content,
					Code: code,
				}
				markDownTopic.Data = append(markDownTopic.Data, codeData);
				codeBlock = CodeBlock{};
			}
		}  else if lineType == 3 || lineType == 4 || lineType == 5 {
			// textListBegin = 0;
			if len(textListData.Textlist) > 0 {
				markDownTopic.Data = append(markDownTopic.Data, textListData);
				textListData = BaseData{};
			}

			if markDownTopic.H != 0 {
				markDownTopicSlice = append(markDownTopicSlice, markDownTopic);
				markDownTopic = MarkDownTopic{}
			}
			markDownTopic.H = lineType;
			markDownTopic.Title = content;
		} else if lineType == 6 {
			// textListBegin = 0;
			if len(textListData.Textlist) > 0 {
				markDownTopic.Data = append(markDownTopic.Data, textListData);
				textListData = BaseData{};
			}

			result.Prologue = content;
		} else if lineType == 7 {
			// textListBegin = 0;
			if len(textListData.Textlist) > 0 {
				markDownTopic.Data = append(markDownTopic.Data, textListData);
				textListData = BaseData{};
			}

			// 可能是头部字段，代码块，也可能是主体内容
			if mdHeaderBegin == 1 {
				// 头部字段
				if strings.HasPrefix(content, "title:") {
					header.Title = strings.TrimSpace(strings.Replace(content, "title:", "", 1));
				} else if strings.HasPrefix(content, "date:") {
					header.Date = strings.TrimSpace(strings.Replace(content, "date:", "", 1));
				} else if strings.HasPrefix(content, "author:") {
					header.Author = strings.TrimSpace(strings.Replace(content, "author:", "", 1));
				} else if strings.HasPrefix(content, "header-img:") {
					header.HeaderImage = strings.TrimSpace(strings.Replace(content, "header-img:", "", 1));
				} else {
					// 忽略
				} 
			} else if codeBlockBegin == 1 {
				// 代码块
				codeBlock = append(codeBlock, content + "\n");
			} else if mdHeaderBegin == 2 {
				// 主体内容
				data := BaseData{
					Text: content,
					ImageUrl: "",
				};
				markDownTopic.Data = append(markDownTopic.Data, data);
			}
		} else if lineType == 8 {
			// textListBegin = 0;
			if len(textListData.Textlist) > 0 {
				markDownTopic.Data = append(markDownTopic.Data, textListData);
				textListData = BaseData{};
			}

			// 图片
			data := BaseData{
				Text: "",
				ImageUrl: content,
			};
			markDownTopic.Data = append(markDownTopic.Data, data);
		} else if lineType == 9 {
			// textListBegin = 1;
			// 子文本列表
			textListData.Textlist = append(textListData.Textlist, content);
		}
    }

	markDownTopicSlice = append(markDownTopicSlice, markDownTopic);

	result.Header = header;
	result.MainData = markDownTopicSlice;

	return result;
}

/*
@return int lineType:
0(none), 不是特殊开头
1(mdheader flag)
// 1(mdheader begin)
// 2(mdheader end) 
2(code block flag)
3(header3)
4(header4)
5(header5)
6(prologue)
7(text)
8(image)
9(子文本列表，以"%d "开头的文本行)
*/
func parseLine(line string) (lineType int, content string) {
	fmt.Println("parseLine:" + line);

	if line == "" {
		return 0, "";
	}

	if strings.HasPrefix(line, "---") {
		return 1, "";
	} else if strings.HasPrefix(line, "```") {
		return 2, "";
	} else if strings.HasPrefix(line, ">") {
		return 6, strings.Replace(line, "> ", "", 1);
	} else if strings.HasPrefix(line, "#####") {
		return 5, strings.Replace(line, "##### ", "", 1);
	} else if strings.HasPrefix(line, "####") {
		return 4, strings.Replace(line, "#### ", "", 1);
	} else if strings.HasPrefix(line, "###") {
		return 3, strings.Replace(line, "### ", "", 1);
	} else if strings.HasPrefix(line, "![") {
		// ![docker-engine](http://106.14.155.231/image/k8s/docker-engine.png)
		// 解析图片地址出来
		re := regexp.MustCompile(`.*\((.*?)\)`)
		parts := re.FindStringSubmatch(line)
		if len(parts) == 2 {
			return 8, parts[1];
		}
		return 0, line;
	} else {
		// 判断是不是文本列表，以 "%d. "开头的行
		re := regexp.MustCompile(`^[0-9]*\.\s(.*)$`)
		parts := re.FindStringSubmatch(line)
		if len(parts) == 2 {
			return 9, parts[1];
		}

		return 7, line;
	}
}