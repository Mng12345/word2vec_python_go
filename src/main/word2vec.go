package main

import (
	"bufio"
	"container/list"
	"fmt"
	"go_workspace/goPython/src/utils"
	"io"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"strings"
	"C"
)

type Word struct {
	content   *string
	count     int
	dimention int
	vector    []float32
	theta     []float32
	sigmod    float32
}

func NewWord(word *string, count int, dimention int) *Word {
	word_ := Word{}
	word_.content = word
	word_.count = count
	word_.dimention = dimention
	word_.vector = utils.RandomFloat32(dimention)
	word_.theta = utils.RandomFloat32(dimention)
	word_.sigmod = 0
	return &word_
}

type DisInfo struct {
	word string
	dis float32
}

type Word2vec struct {
	WordReader   *WordReader
	Window       int
	Threads      int
	Sentence     []*Word
	SentenceSize int
	wordCount    int
	dimention    int
	negative     int
	outputFile   string
	wordMap      *map[string]*Word
	wordList     []*string
	lr           float32
	EOF          bool
}

func NewWord2vec(trainFile string, outputFile string, window int, threads int,
	sentenceSize int, dimention int, negative int, lr float32) (*Word2vec) {
	word2vec := Word2vec{}
	word2vec.WordReader = NewWordReader(trainFile)
	word2vec.Window = window
	word2vec.Threads = threads
	word2vec.SentenceSize = sentenceSize
	word2vec.Sentence = make([]*Word, 0, sentenceSize)
	word2vec.wordCount = 0
	word2vec.dimention = dimention
	word2vec.outputFile = outputFile
	word2vec.negative = negative
	word2vec.lr = lr
	word2vec.EOF = false
	return &word2vec
}

func (word2vec *Word2vec) ReadWordMap() {
	fmt.Println("开始读取文件")
	wordMap := make(map[string]*Word)
	for {
		word := word2vec.WordReader.ReadWord()
		if *word == WORDREADER_STOPFLAG {
			break
		}
		if *word == " " {
			continue
		}
		if v, ok := wordMap[*word]; ok {
			v.count += 1
		} else {
			wordObj := NewWord(word, 1, word2vec.dimention)
			wordMap[*word] = wordObj
		}
		word2vec.wordCount += 1
		if word2vec.wordCount % 10000 == 0 {
			fmt.Println("单词数量：", word2vec.wordCount / 1000, "k")
		}
	}
	word2vec.wordMap = &wordMap
	fmt.Println("读取文件完毕")
}

func (word2vec *Word2vec) start() {
	// 读取wordMap
	word2vec.ReadWordMap()
	// 生成单词列表
	word2vec.generateWordList()
	// 训练
	word2vec.train()
}

func (word2vec *Word2vec) readSentence() {
	if len(word2vec.Sentence) == 0 {
		for i := 0; i < word2vec.SentenceSize; i++ {
			word := word2vec.WordReader.ReadWord()
			word2vec.Sentence = append(word2vec.Sentence, (*word2vec.wordMap)[*word])
		}
	} else {
		word := word2vec.WordReader.ReadWord()
		if *word == WORDREADER_STOPFLAG {
			word2vec.EOF = true
			return
		}
		word2vec.Sentence = append(word2vec.Sentence, (*word2vec.wordMap)[*word])
		word2vec.Sentence = word2vec.Sentence[1:]
	}
}

func (word2vec *Word2vec) train() {
	count := 0
	currLr := word2vec.lr
	for {
		if word2vec.EOF {
			break
		}
		currLr = word2vec.lr * (1.0 -  float32(count) / float32(word2vec.wordCount))
		count += 1
		if currLr < word2vec.lr * 0.001 {
			currLr = word2vec.lr * 0.001
		}
		// 读取句子
		word2vec.readSentence()
		centerWord := word2vec.Sentence[word2vec.SentenceSize/2]
		// 计算xw
		xw := word2vec.calXw(centerWord.content)
		// 随机取出negative个负样本
		negativeWords := word2vec.sampleNegative(centerWord.content)
		pBeforeUpdate := calP(xw, centerWord, negativeWords)
		// 计算正例和负例的sigmod值
		centerWord.sigmod = utils.Sigmod(xw, centerWord.theta)
		for _, negativeWord := range negativeWords {
			negativeWord.sigmod = utils.Sigmod(xw, negativeWord.theta)
		}
		deltaXw := utils.ArrayMultiply1((1 - centerWord.sigmod), centerWord.theta)
		deltaTheta0 := utils.ArrayMultiply1(1-centerWord.sigmod, xw)
		centerWord.theta = utils.ArrayAdd(centerWord.theta, utils.ArrayMultiply1(currLr, deltaTheta0))
		for _, negativeWord := range negativeWords {
			currDeltaXw := utils.ArrayMultiply1(-1*negativeWord.sigmod, negativeWord.theta)
			deltaXw = utils.ArrayAdd(deltaXw, currDeltaXw)
			deltaThetai := utils.ArrayMultiply1(-1*negativeWord.sigmod, xw)
			negativeWord.theta = utils.ArrayAdd(negativeWord.theta, utils.ArrayMultiply1(currLr, deltaThetai))
		}
		// 更新上下文的vector
		for _, word := range word2vec.Sentence {
			if *word.content == *centerWord.content {continue}
			word.vector = utils.ArrayAdd(word.vector, utils.ArrayMultiply1(currLr, deltaXw))
		}
		newXw := word2vec.calXw(centerWord.content)
		pAfterUpdate := calP(newXw, centerWord, negativeWords)
		if (count+1) % 10000 == 0 {
			fmt.Println("学习率：", currLr, "已训练：", count*100/word2vec.wordCount, "%, 概率=更新前->更新后: ",
				pBeforeUpdate, "->", pAfterUpdate)
		}
	}
}

// 随机取出negative个负样本
func (word2vec *Word2vec) sampleNegative(centerWord *string) []*Word {
	res := make([]*Word, 0, word2vec.negative)
	for i := 0; i < word2vec.negative; i++ {
		var word *string
		for {
			word = word2vec.wordList[rand.Intn(word2vec.wordCount)]
			if *word != *centerWord {
				break
			}
		}
		res = append(res, (*word2vec.wordMap)[*word])
	}
	return res
}

// 计算正例和负例出现的概率
func calP(xw []float32, centerWord *Word, negativeWords []*Word) float32 {
	p := float32(1);
	p *= utils.Sigmod(xw, centerWord.theta)
	for _, word := range negativeWords {
		p *= 1-utils.Sigmod(xw, word.theta)
	}
	return p
}

func (word2vec *Word2vec) generateWordList() {
	word2vec.wordList = make([]*string, 0, word2vec.wordCount)
	for k, v := range *word2vec.wordMap {
		word := k
		for i := 0; i < v.count; i++ {
			word2vec.wordList = append(word2vec.wordList, &word)
		}
	}
}

func (word2vec *Word2vec) calXw(word *string) []float32 {
	res := make([]float32, word2vec.dimention, word2vec.dimention)
	for i := 0; i < word2vec.dimention; i++ {
		res[i] = 0
		for j := 0; j < word2vec.SentenceSize; j++ {
			if *word2vec.Sentence[j].content == *word {
				continue
			}
			res[i] += word2vec.Sentence[j].vector[i]
		}
		res[i] /= float32(word2vec.SentenceSize)
	}
	return res
}

type WordReader struct {
	Train_file string
	file       *os.File
	Reader     *bufio.Reader
	wordBuf    *list.List // 存储每次读取的行分割后的单词
}

var (
	WORDREADER_STOPFLAG = "<=stop=>"
)

func NewWordReader(train_file string) (*WordReader) {
	file, err := os.Open(train_file)
	if err != nil {
		fmt.Println("打开文件", train_file, "失败")
		os.Exit(-1)
	}
	reader := bufio.NewReader(file)
	wordBuf := list.New()
	return &WordReader{train_file, file, reader, wordBuf}
}

func (reader *WordReader) Close() {
	if reader.file != nil {
		reader.file.Close()
		reader.file = nil
	}
}

// 读取一个单词
func (reader *WordReader) ReadWord() (*string) {
	if reader.file == nil {
		file, err := os.Open(reader.Train_file)
		if err != nil {
			fmt.Println("打开文件", reader.Train_file, "失败")
			os.Exit(-1)
		}
		reader.file = file
		reader.Reader = bufio.NewReader(reader.file)
		return &WORDREADER_STOPFLAG
	} else {
		if reader.wordBuf.Len() == 0 {
			// 重新读取一行
			byteLine, _, err := reader.Reader.ReadLine()
			if err == io.EOF {
				file, err := os.Open(reader.Train_file)
				if err != nil {
					fmt.Println("打开文件", reader.Train_file, "失败")
					os.Exit(-1)
				}
				reader.file = file
				reader.Reader = bufio.NewReader(reader.file)
				return &WORDREADER_STOPFLAG
			}
			if err != nil {
				fmt.Println("读取文件", reader.Train_file, "失败")
				os.Exit(-1)
			}
			line := string(byteLine)
			words := strings.Split(line, " ")
			for i := 0; i < len(words); i++ {
				reader.wordBuf.PushBack(words[i])
			}
			word := reader.wordBuf.Front()
			reader.wordBuf.Remove(word)
			wordStr := word.Value.(string)
			return &wordStr
		} else {
			word := reader.wordBuf.Front()
			reader.wordBuf.Remove(word)
			wordStr := word.Value.(string)
			return &wordStr
		}
	}
}

func RecongnizeWord(dimention int, outputFile string) {
	file, err := os.Open(outputFile)
	if err != nil {
		fmt.Println("打开文件", outputFile, "失败")
		os.Exit(-1)
	}
	reader := bufio.NewReader(file)
	defer file.Close()
	wordMap := make(map[string][]float32)
	for {
		byteLine, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("读取文件", outputFile, "失败")
			os.Exit(-1)
		}
		line := string(byteLine)
		word, vecotr := utils.String2Vector(line)
		if len(vecotr) != dimention {
			continue
		}
		wordMap[word] = vecotr
	}
	for {
		var word string
		fmt.Println("请输入单词：(:q->退出)")
		fmt.Scanln(&word)
		if word == ":q" {
			fmt.Println("退出")
			os.Exit(0)
		}
		if _, ok := wordMap[word]; !ok {
			fmt.Println("单词：", word, "不在词库中")
			continue
		}
		nearNum := 20
		nearWords := make([]DisInfo, nearNum, nearNum)
		for w, v := range wordMap {
			if w == word {continue}
			dis := utils.CosDistance(wordMap[word], v)
			if len(nearWords) < nearNum {
				item := DisInfo{w, dis}
				nearWords = append(nearWords, item)
			} else {
				for i, item := range nearWords {
					if item.dis < dis {
						item := DisInfo{w, dis}
						nearWords[i] = item
						break
					}
				}
			}
		}
		// 打印最匹配单词
		fmt.Println("======================匹配结果======================")
		for _, item := range nearWords {
			fmt.Println(item.word, " ", item.dis)
		}
		fmt.Println("======================匹配结果======================")
	}
}

// 保存单词
func (word2vec *Word2vec) SaveWord() {
	os.Create(word2vec.outputFile)
	file, err := os.OpenFile(word2vec.outputFile, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("打开文件", word2vec.outputFile, "失败")
		os.Exit(-1)
	}
	writer := bufio.NewWriter(file)
	defer file.Close()
	for _, word := range *word2vec.wordMap {
		line := *word.content + *utils.Vector2String(word.vector) + "\n"
		writer.WriteString(line)
		writer.Flush()
	}
}

//export Run
func Run(trainFile_ *C.char, vectorFile_ *C.char, windowSize int, dimention int,
	negative int, lr_ C.float, show int) {
	trainFile := C.GoString(trainFile_)
	vectorFile := C.GoString(vectorFile_)
	lr := float32(lr_)
	word2vec := NewWord2vec(trainFile, vectorFile, windowSize,
		1, windowSize, dimention, negative, lr)
	defer word2vec.WordReader.Close()
	if show == 1 {
		RecongnizeWord(dimention, vectorFile)
		return
	}
	word2vec.start()
	word2vec.SaveWord()
	RecongnizeWord(dimention, vectorFile)
}

func main() {

}
