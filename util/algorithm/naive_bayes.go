package algorithm

import (
	"bufio"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/go-ego/gse"
	"io"
	"math"
	"os"
	"strings"
)

var (
	seg gse.Segmenter
)

type FileOperate struct {
	FilePath string
}

type Bayes struct {
	words      mapset.Set // 分词后的词集合
	WordsCount int        // 分词后的词集合的数量

	CountAd int // 广告短信的数量
	CountNd int // 非广告短信的数量

	WordCountAd int // 广告词的数量
	WordCountNd int // 非害广告词的数量

	WordAdMap map[string]int // 广告中出现的词的数量
	WordNdMap map[string]int // 非广告中出现的词的数量

}

// ReadFile 读取文件
func (f *FileOperate) ReadFile() ([][]string, []string) {
	file, err := os.OpenFile(f.FilePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	data := make([]string, 0)
	// 读取文件内容
	for {
		str, err := reader.ReadString('\n')
		data = append(data, strings.Replace(str, "\n", "", -1))
		if err == io.EOF {
			break
		}
	}
	fmt.Printf("数据总条目：%d\n", len(data))

	X := make([][]string, 0)
	Y := make([]string, 0)
	for _, v := range data {
		str := strings.TrimSpace(v)
		line := strings.Split(str, "\t")
		if len(line) != 2 {
			continue
		}
		X = append(X, f.splitChinese(line[1]))
		Y = append(Y, line[0])
		//if i%10000 == 0 {
		//	fmt.Println(i)
		//}
	}
	return X, Y
}

// 分词
func (f *FileOperate) splitChinese(data string) []string {
	return seg.Slice(data)
}

func (b *Bayes) Init() {
	b.words = mapset.NewSet()
	b.WordAdMap = make(map[string]int)
	b.WordNdMap = make(map[string]int)
}

// GetWordSet 计算词数量
func (b *Bayes) GetWordSet(X [][]string, Y []string) {
	for i, v := range X {
		if Y[i] == "1" {
			b.CountAd++
		} else {
			b.CountNd++
		}
		for _, word := range v {
			b.words.Add(word)
			if Y[i] == "1" {
				// 统计词频
				if _, ok := b.WordAdMap[word]; ok {
					b.WordAdMap[word]++
				} else {
					b.WordAdMap[word] = 1
				}

				b.WordCountAd++
				//b.wordsAd = append(b.wordsAd, word)
			} else {
				// 统计词频
				if _, ok := b.WordNdMap[word]; ok {
					b.WordNdMap[word]++
				} else {
					b.WordNdMap[word] = 1
				}

				b.WordCountNd++
				//b.wordsNd = append(b.wordsNd, word)
			}
		}
	}
	b.WordsCount = b.words.Cardinality()
}

// Predict 预测
func (b *Bayes) Predict(words []string) bool {
	adProb := math.Log(float64(b.CountAd) / float64(b.CountNd+b.CountAd))
	ndProb := math.Log(float64(b.CountNd) / float64(b.CountNd+b.CountAd))

	for _, word := range words {
		adProb += math.Log(float64(b.WordAdMap[word]+1) / float64(b.WordsCount))
		ndProb += math.Log(float64(b.WordNdMap[word]+1) / float64(b.WordsCount))
	}

	return adProb > ndProb
}
