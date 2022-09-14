package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 线程数=CPU核心数
var treadNum = runtime.NumCPU()

func main() {
	menu()
}

func menu() {
	fmt.Println("1.统计单词出现次数")
	fmt.Println("2.统计字符出现次数")
	fmt.Println("3.退出")
	id := 0
	fmt.Scan(&id)
	switch id {
	case 1:
		//读取文件
		str := ReadFile()
		//开始时间
		startTime := time.Now()
		WordSort(str)
		//结束时间
		endTime := time.Now()
		fmt.Println("总耗时：", endTime.Sub(startTime))
	case 2:
		fmt.Println("请输入字符串：")
		var str string
		fmt.Scan(&str)
		CharSort(str)
	case 3:
		os.Exit(0)
	default:
		fmt.Println("输入错误")
	}
}

func WordSort(str string) []WordFrequent {
	str = strings.ToLower(str)
	var wordRank []WordFrequent
	allWord := int64(0)

	var wg sync.WaitGroup
	start := 0
	end := len(str) / treadNum
	//遍历开始时间
	startTime := time.Now()
	for i := 0; i < treadNum; i++ {
		wg.Add(1)
		if str[end] != ' ' {
			for str[end] >= 'a' && str[end] <= 'z' || str[end] >= 'A' && str[end] <= 'Z' || str[end] >= '0' && str[end] <= '9' {
				//fmt.Println("end:", end)
				end++
				if end >= len(str) {
					break
				}
			}
		}
		//fmt.Println("start:", start, "end:", end)
		//fmt.Println(111)
		if start != 0 {
			start = end + 1
			end = start + len(str)/treadNum
			if end > len(str) {
				end = len(str)
			}
		}
		if i == treadNum-1 {
			end = len(str)
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
			//fmt.Println(i3, "替换完成")
			arrSlice := strings.Split(str, " ")
			for _, v := range arrSlice {
				if v != "" {
					if strings.Contains(v, "'") {
						arr := strings.Split(v, "'")
						v = arr[0]
					}
					//如果单词已经存在，频率+1
					if len(wordRank) > 0 {
						for i2 := 0; i2 < len(wordRank); i2++ {
							if wordRank[i2].Word == v {
								atomic.AddInt64(&wordRank[i2].Frequent, 1)
								atomic.AddInt64(&allWord, 1)
								break
							}
							if i2 == len(wordRank)-1 {
								wordRank = append(wordRank, WordFrequent{v, 1})
								atomic.AddInt64(&allWord, 1)
								break
							}
						}
					} else {
						wordRank = append(wordRank, WordFrequent{v, 1})
						atomic.AddInt64(&allWord, 1)
					}
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

	fmt.Println("排名\t单词\t出现次数\t出现频率")
	wordRank = sortWordFrequence(wordRank)
	for i := 0; i < 1000; i++ {
		fmt.Printf("%d\t%s\t%d\t%.2f%%\n", i+1, wordRank[i].Word, wordRank[i].Frequent, float64(wordRank[i].Frequent)/float64(allWord)*100)
	}
	return wordRank
}

func CharSort(str string) []WordFrequent {
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
	for _, v := range charRank {
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

func ReadFile() string {
	startTime := time.Now()
	file, err := os.OpenFile("/Users/wuwuwu/soft-project/task1/WordFrequence/gone-with-the-wind.txt", os.O_RDONLY, 0666)
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
