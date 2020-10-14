package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

/*
a	动画
b	漫画
c	游戏
d	文学
e	原创
f	来自网络
g	其他
h	影视
i	诗词
j	网易云
k	哲学
l	抖机灵 	//no data yet
*/

var apiUrl string = "https://v1.hitokoto.cn/?c="
var paramCMaxNum int = 11

//监听系统信号
var signalChan chan os.Signal = make(chan os.Signal, 1)

var stopped = make(chan struct{}) //只能由中间层 close

var cacheMaxNum int = 1000 // id 通道最大缓存数
var idStrChan chan string = make(chan string, cacheMaxNum)

//参数类型是否达到上限记录
var paramCMap = make(map[string]bool, paramCMaxNum)

//开始信号
var startChan chan struct{} = make(chan struct{})

//保存目录
var dataPath string
var perm os.FileMode = 0644 //默认保存文件权限

//重复多少次之后停止
var repeatNum int = 50

func init() {
	//使用大部分cpu核心数
	numCpu := runtime.NumCPU() * 3 / 4
	if numCpu < 1 {
		numCpu = 1
	}

	runtime.GOMAXPROCS(numCpu)

	//获取并创建数据目录
	log.SetFlags(0)

	curPath := getwd()
	dataPath = filepath.Join(curPath, "data", "hitokoto")
	err := os.MkdirAll(dataPath, perm)
	if err != nil {
		log.Println("mkdir data direcotry fail ,", err)
		panic(err)
	}
}

//获取当前目录
func getwd() string {
	cur, err := os.Getwd()
	if err != nil {
		log.Println("get current direcotry fail ,", err)
		panic(err)
	}

	return cur
}

//中间层，判断 id 重复数，决定是否停止
func middlePart() {
	startChan <- struct{}{}

	defer close(stopped)

	//记录历史获得的id列表
	idStrMap := make(map[string]int, 10000)

	for baseName := range idStrChan {
		tmp := strings.Split(baseName, "_")
		paramC := tmp[0]
		idStr := tmp[1]

		if _, exists := paramCMap[paramC]; exists {
			if len(paramCMap) >= paramCMaxNum {
				return
			}

			continue
		}

		v, ok := idStrMap[idStr]
		if ok && v > repeatNum {
			paramCMap[paramC] = true
			if len(paramCMap) >= paramCMaxNum {
				return
			}
			continue
		}

		idStrMap[idStr] = v + 1
	}
}

func main() {
	//记录总耗时
	defer func(beginTime time.Time) {
		log.Printf("done  ----- cost %v \n", time.Since(beginTime))
	}(time.Now())

	//1、监听中断信号
	signal.Notify(signalChan, os.Interrupt)

	//2、中间层来处理条件终止信号，比如多次获取到重复数据
	go middlePart()

	//开始
	<-startChan
	log.Println("begin to collect")

	//3、无线循环开启新 goroutine
	i := 0
	for {
		//进入此处正好被 stopped
		select {
		case <-stopped:
			return
		default:
		}

		select {
		case <-signalChan:
			log.Printf("the system is downing, total open goroutine num %d\n", i)
			return
		case <-stopped:
			return
		default:
			for c := 'a'; c < 'a'+int32(paramCMaxNum); c++ {
				paramC := string(c)
				if _, exists := paramCMap[paramC]; exists {
					continue
				}

				i++
				go download(i, c)
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

//连接3个字符串
func joinAB(a, b, c string) string {
	return (" " + a + b + c)[1:]
}

//开启第几个 goroutine 下载
func download(num int, c rune) {
	paramC := string(c)
	targetUrl := joinAB(apiUrl, "", paramC)
	resp, err := http.Get(targetUrl)
	if err != nil {
		log.Printf("goroutine[%d] get url[%s] content fail, err:%v \n", num, targetUrl, err)
		return
	}

	defer resp.Body.Close() //释放连接

	//http回包状态值判断
	if resp.StatusCode != http.StatusOK {
		log.Printf("goroutine[%d] url[%s] response StatusCode is %d\n", num, targetUrl, resp.StatusCode)
		return
	}

	//读取数据
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("goroutine[%d] url[%s] body read fail, err:%v\n", num, targetUrl, err)
		return
	}

	//保存数据
	idByte := content[bytes.IndexByte(content, ':')+1 : bytes.IndexByte(content, ',')]
	id := string(idByte)

	log.Printf("goroutine[%d] from url[%s] get id %v \n", num, targetUrl, id)

	baseName := joinAB(paramC, "_", id)

	//保存文件
	//paramCPath := filepath.Join(dataPath, paramC)
	//os.MkdirAll(paramCPath, perm)
	//file := filepath.Join(paramCPath, joinAB(baseName, "", ".json"))

	file := filepath.Join(dataPath, "id_"+id+".json")
	ioutil.WriteFile(file, content, perm)

	//传递 id
	idStrChan <- baseName
}
