package main

import (
	"fmt"
	"testing"
)

func TestReadSence(t *testing.T) {
	dimention := 100
	outputFile := "../../data/vectors"
	word2vec := NewWord2vec("../../data/text8split", outputFile, 10,
		8, 10, dimention, 6, 0.05)
	defer word2vec.WordReader.Close()
	word2vec.ReadWordMap()
	word2vec.generateWordList()
	for i:=0; i<10; i++ {
		word2vec.readSentence()
		for _, v := range word2vec.Sentence {
			fmt.Print(*v.content, " ")
		}
		fmt.Println()
		if len(word2vec.Sentence) != word2vec.SentenceSize {
			t.Errorf("readSentence读取错误，读取得到的数组长度为：%d与SentenceSize不一致",
				len(word2vec.Sentence))
		}
	}
}