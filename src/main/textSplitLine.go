package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// 文本分行
func splitLine(file string, fileOutput string) {
	f, err := os.Open(file)
	//os.Create(fileOutput)
	fout, err1 := os.OpenFile(fileOutput, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("打开文件", file, "失败")
		os.Exit(-1)
	}
	if err1 != nil {
		fmt.Println("打开文件", fileOutput, "失败")
		os.Exit(-1)
	}
	reader := bufio.NewReader(f)
	writer := bufio.NewWriter(fout)
	defer func() {
		f.Close()
		fout.Close()
	} ()
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("读取文件", file, "失败")
			os.Exit(-1)
		}

		lineStr := string(line)
		line1 := lineStr[: len(lineStr)/2]
		line2 := lineStr[len(lineStr) / 2:]
		_, err = writer.WriteString(line1 + "\n")
		if err != nil {
			fmt.Println(err)
			fmt.Println("写文件", fileOutput, "失败")
		}
		_, err = writer.WriteString(line2 + "\n")
		if err != nil {
			fmt.Println(err)
			fmt.Println("写文件", fileOutput, "失败")
		}
		writer.Flush()
	}
}

//func main() {
//	splitLine("data/text8", "data/text8split")
//}
