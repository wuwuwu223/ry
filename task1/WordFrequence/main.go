package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// 线程数=CPU核心数
var treadNum = 20 * runtime.NumCPU()
var stopWordMap = make(map[string]bool)
var vWordMap = make(map[string]string)

func main() {
	menu()
}

func menu() {
	//wf -c 文件名 -n 10
	var fileName1 string
	var fileName2 string
	var filedir string
	var stopWordsFile string
	var vWordsFile string
	var n int
	flag.StringVar(&fileName1, "c", "default", "排序字符文件名")
	flag.StringVar(&fileName2, "f", "default", "排序单词文件名")
	flag.StringVar(&filedir, "d", "default", "排序单词文件夹")
	flag.StringVar(&stopWordsFile, "x", "default", "停用词文件名")
	flag.StringVar(&vWordsFile, "v", "default", "替换词文件名")
	flag.IntVar(&n, "n", 0, "输出前n个单词/字符")
	flag.Parse()

	if stopWordsFile != "default" {
		stopWordsStr := ReadFile(stopWordsFile)
		stopWords := strings.Split(stopWordsStr, "\n")
		for _, v := range stopWords {
			stopWordMap[v] = true
		}
	}

	if vWordsFile != "default" {
		vWordsStr := ReadFile(vWordsFile)
		arr1 := strings.Split(vWordsStr, "\n")
		for i := range arr1 {
			arr2 := strings.Split(arr1[i], " ")
			for j := range arr2 {
				vWordMap[arr2[j]] = arr2[0]
			}
		}
	}

	if fileName1 == "default" && fileName2 == "default" && filedir == "default" {
		fmt.Println("请输入文件名或文件夹名")
		return
	}
	if fileName1 != "default" && fileName2 == "default" && filedir == "default" {
		str := ReadFile(fileName1)
		CharSort(str, n)
		return
	}
	if fileName2 != "default" && filedir == "default" {
		str := ReadFile(fileName2)
		WordSort(str, n)
		return
	}
	if filedir != "default" {
		str := ReadDir(filedir)
		WordSort(str, n)
		return
	}

	fmt.Println("请输入文件名或文件夹名")
}

func WordSort(str string, n int) []WordFrequent {
	str = strings.ToLower(str)
	var wordMap sync.Map
	var mutex sync.Mutex
	var wordRank []WordFrequent
	allWord := int64(0)

	var wg sync.WaitGroup
	start := 0
	end := len(str) / treadNum
	//遍历开始时间
	startTime := time.Now()
	if len(str) < treadNum {
		treadNum = 1
	}
	for i := 0; i < treadNum; i++ {
		wg.Add(1)
		if start != 0 {
			start = end + 1
			end = start + len(str)/treadNum
			if end > len(str) {
				end = len(str)
			}
		}

		if i == treadNum-1 {
			end = len(str)
		} else {
			if str[end] != ' ' {
				for str[end] >= 'a' && str[end] <= 'z' || str[end] >= 'A' && str[end] <= 'Z' || str[end] >= '0' && str[end] <= '9' || str[end] == '\'' {
					//fmt.Println("end:", end)
					end++
					if end >= len(str) {
						break
					}
				}
			}
		}
		txt := str[start:end]
		//fmt.Println("start:", start, "end:", end)
		//i3 := i
		go func(start int, end int, str string) {
			//替换所有非字母数字的字符为空格
			for j := 0; j < len(str); j++ {
				if !(str[j] >= 'a' && str[j] <= 'z' || str[j] >= 'A' && str[j] <= 'Z' || str[j] >= '0' && str[j] <= '9' || str[j] == '\'') {
					str = str[:j] + " " + str[j+1:]
				}
			}
			//fmt.Println("str:", str)
			//fmt.Println(i3, "替换完成")
			arrSlice := strings.Split(str, " ")
			for _, v := range arrSlice {
				if v != "" {
					if strings.Contains(v, "'") {
						arr := strings.Split(v, "'")
						v = arr[0]
					}
					if _, ok := stopWordMap[v]; ok {
						continue
					}
					if _, ok := vWordMap[v]; ok {
						v = vWordMap[v]
					}
					//如果单词已经存在，频率+1
					mutex.Lock()
					allWord++
					if value, ok := wordMap.Load(v); ok {
						wordMap.Store(v, value.(int)+1)
					} else {
						wordMap.Store(v, 1)
					}
					mutex.Unlock()
				}
			}

			wg.Done()
		}(start, end, txt)
		start = end + 1
	}
	wg.Wait()
	//遍历结束时间
	endTime := time.Now()
	fmt.Println("遍历耗时：", endTime.Sub(startTime))
	//把map转换成切片
	wordMap.Range(func(key, value interface{}) bool {
		wordRank = append(wordRank, WordFrequent{key.(string), int64(value.(int))})
		return true
	})

	fmt.Println("排名\t单词\t出现次数\t出现频率")
	wordRank = sortWordFrequence(wordRank)
	for i := range wordRank {
		if n != 0 && i == n {
			break
		}
		fmt.Printf("%d\t%s\t%d\t%.2f%%\n", i+1, wordRank[i].Word, wordRank[i].Frequent, float64(wordRank[i].Frequent)/float64(allWord)*100)
	}
	return wordRank
}

func CharSort(str string, n int) []WordFrequent {
	freqmap := make(map[string]int)
	allChar := 0
	//遍历开始时间
	startTime := time.Now()
	for i := 0; i < len(str); i++ {
		if str[i] >= 'a' && str[i] <= 'z' || str[i] >= 'A' && str[i] <= 'Z' || str[i] >= '0' && str[i] <= '9' {
			freqmap[string(str[i])]++
			allChar++
		}
	}
	//遍历结束时间
	endTime := time.Now()
	fmt.Println("遍历耗时：", endTime.Sub(startTime))
	var charRank []WordFrequent
	for k, v := range freqmap {
		charRank = append(charRank, WordFrequent{k, int64(v)})
	}

	fmt.Println("字符\t出现次数\t出现频率")
	charRank = sortWordFrequence(charRank)
	for i, v := range charRank {
		if n != 0 && i >= n {
			break
		}
		fmt.Printf("%s\t%d\t%.2f%%\n", v.Word, v.Frequent, float64(v.Frequent)/float64(allChar)*100)
	}
	return charRank
}

// 排序，如果频率相同，按字母顺序排序
func sortWordFrequence(wf []WordFrequent) []WordFrequent {
	for i := 0; i < len(wf); i++ {
		for j := i + 1; j < len(wf); j++ {
			if wf[i].Frequent < wf[j].Frequent {
				wf[i], wf[j] = wf[j], wf[i]
			} else if wf[i].Frequent == wf[j].Frequent {
				if wf[i].Word > wf[j].Word {
					wf[i], wf[j] = wf[j], wf[i]
				}
			}
		}
	}
	return wf
}

func ReadFile(filename string) string {
	startTime := time.Now()
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("文件打开失败")
		fmt.Println(err)
		return ""
	}
	//读取整个文件
	buf := make([]byte, 1024*4)
	var str string
	for {
		n, _ := file.Read(buf)
		if n == 0 {
			break
		}
		str += string(buf[:n])
	}
	endTime := time.Now()
	fmt.Println("读取文件耗时：", endTime.Sub(startTime))
	return str
}

// readdir return string
func ReadDir(dir string) string {
	var str string
	startTime := time.Now()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("文件打开失败")
		fmt.Println(err)
		return ""
	}
	for _, file := range files {
		if file.IsDir() {
			str += ReadDir(dir + "/" + file.Name())
		} else {
			str += " " + ReadFile(dir+"/"+file.Name())
		}
	}
	endTime := time.Now()
	fmt.Println("读取文件耗时：", endTime.Sub(startTime))
	return str
}
