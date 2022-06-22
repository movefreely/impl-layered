package service

import (
	"encoding/gob"
	"github.com/go-ego/gse"
	"in-server/util/algorithm"
	"os"
)

type Tools struct{}

var (
	bayes = &algorithm.Bayes{}
	seg   gse.Segmenter
)

func InitBayes() {
	seg.LoadDict()
	file, err := os.Open("wordGob")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&bayes)
	if err != nil {
		panic(err.Error())
	}
}

func (t *Tools) JudgeIsSpam(message string) bool {
	return bayes.Predict(seg.Slice(message))
}
