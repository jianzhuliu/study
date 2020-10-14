/*
每日金句
1、数据来源 v1.hitokoto.cn

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"path/filepath"
	"html/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jianzhuliu/tools/browser"
)

//数据库配置
const (
	DB_HOST string = "root"
	DB_PASSWD string =""
	DB_NAME string ="yiyan"
	DB_TABLE string ="yiyan"
)

var (
	port = flag.Int("port", 1919, "set port")
	host = flag.String("host", "localhost", "set host")

	rootPath string
	dataPath string

	wg     sync.WaitGroup
	tokens chan struct{} = make(chan struct{}, 100)

	db *sql.DB
	dbMaxId int 
)

var tpl *template.Template 


//初始化操作
func init() {
	log.SetFlags(0)
	flag.Parse()

	//以当前执行脚本为根目录
	curPath, err := os.Getwd()
	if err != nil {
		log.Println("fail to do os.Getwd()")
		return
	}

	rootPath = curPath
	dataPath = filepath.Join(rootPath, "data","hitokoto")

	dns := fmt.Sprintf("%s:%s@/%s",DB_HOST,DB_PASSWD,DB_NAME)
	db, err = sql.Open("mysql", dns)
	if err != nil {
		log.Println("connect to db fail")
		return
	}
	
	//加载模板
	tpl = template.Must(template.New("yiyan").Parse(homeTmp))
}

//主入口文件
func main() {
	switch flag.Arg(0) {
	case "init":
		initDb()
	default:
		openWeb()
	}
}

//初始化db
func initDb() {
	defer func(beginTime time.Time) {
		log.Printf("done -- spend %v \n", time.Since(beginTime))
	}(time.Now())

	_, err := db.Exec(fmt.Sprintf("drop table if exists `%s`;", DB_TABLE))
	if err != nil {
		log.Printf("drop table(%s) fail -- %v \n", DB_TABLE, err)
		return
	}
	_, err = db.Exec(fmt.Sprintf("create table if not exists `%s` (`id`  int(10) unsigned not null primary key auto_increment,`type` char(1) not null,`word` varchar(100) not null,`from` varchar(50) default null,`from_who` varchar(30) default null,`length`  int(10) unsigned not null,`created` timestamp not null default CURRENT_TIMESTAMP)engine=Innodb default charset=utf8mb4 collate utf8mb4_unicode_ci;", DB_TABLE))

	if err != nil {
		log.Printf("create table(%s) fail -- %v \n", DB_TABLE, err)
		return
	}

	log.Println("begin to init db")

	//匹配目录文件
	matches, err := filepath.Glob(filepath.Join(dataPath, "id_*.json"))
	if err != nil {
		log.Println("read dir fail", err)
		return
	}

	for i, filename := range matches {
		wg.Add(1)

		//开启新 goroutine 处理
		go toDb(filename)

		if i == 2 {
			//break
		}
	}

	wg.Wait()
}

type Yiyan struct {
	Id       int    `json:"id"`
	Uuid     string `json:"uuid"`
	Hitokoto string `json:"hitokoto"`
	Type     string `json:"type"`
	From     string `json:"from"`
	From_who string `json:"from_who,omitempty"`
	Length   int    `json:"length"`
}

//保存至数据库
func toDb(filename string) {
	defer wg.Done()
	tokens <- struct{}{}
	defer func() {
		<-tokens
	}()

	f, err := os.Open(filename)
	if err != nil {
		log.Printf("%s -- open fail, err= %v \n", filename, err)
		return
	}

	obj := Yiyan{}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&obj)
	if err != nil {
		log.Printf("%s -- decode fail, err= %v \n", filename, err)
		return
	}

	result, err := db.Exec("insert into yiyan(`type`,`word`,`from`,`from_who`,`length`) values(?,?,?,?,?)", obj.Type, obj.Hitokoto, obj.From, obj.From_who, obj.Length)

	if err != nil {
		log.Printf("%s -- insert fail, err= %v \n", filename, err)
		return
	}

	lastId, _ := result.LastInsertId()
	log.Println("succ", obj, lastId)
}

//////////////////////////////////////////////////////////////web
//web 服务
func openWeb() {
	//获取数据库中最大的值
	err := db.QueryRow(fmt.Sprintf("select count(1) as c from %s limit 1", DB_TABLE)).Scan(&dbMaxId)
	if err != nil {
		log.Printf("query db max id fail, %d, -- %v \n", dbMaxId, err)
		return
	}
	
	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("begin to open web at %s \n", addr)
	go browser.OpenWithNotice(addr)

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/yiyan", handleYiyan)

	//静态文件
	http.Handle("/static/", http.FileServer(http.Dir(rootPath)))
	http.Handle("/favicon.ico", http.FileServer(http.Dir(rootPath)))

	log.Println(http.ListenAndServe(addr, nil))
}

func JSON(w http.ResponseWriter, r *http.Request, statusCode int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

//首页
func handleIndex(w http.ResponseWriter, r *http.Request) {
	
	obj, err := getRandIdJson()
	if err != nil {
		obj.Hitokoto = "欢迎光临"
		fmt.Println("getRandIdJson -- fail ",err)
	}
	
	tpl.Execute(w, obj)
}

//一言
func handleYiyan(w http.ResponseWriter, r *http.Request) {
	obj, err := getRandIdJson()
	status := http.StatusOK
	var data []byte
	if err != nil {
		fmt.Println("getRandIdJson -- fail ",err)
		status = http.StatusInternalServerError
	} else {
		data, err = json.Marshal(obj)
		if err != nil {
			fmt.Println("json.Marshal -- fail ",err)
			status = http.StatusInternalServerError
		}
	}
	
	JSON(w, r, status, data)
}

//获取随机
func getRandIdJson() (Yiyan, error) {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(dbMaxId) + 1
	obj := Yiyan{}
	err := db.QueryRow(fmt.Sprintf("select word,`from`,`from_who` from %s where id= %d limit 1", DB_TABLE, num)).Scan(&obj.Hitokoto,&obj.From,&obj.From_who)
	return obj, err 
}

var homeTmp string = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<title>一言</title>
<link rel="stylesheet" href="/static/css/jquery.fullPage.css">
<style>
.section { text-align: left; font: 50px "Microsoft Yahei"; color: #fff;}

</style>
<script src="/static/js/jquery-1.8.3.min.js"></script>
<script src="/static/js/jquery.fullPage.min.js"></script>
<script>
$(function(){
	$('#dowebok').fullpage({
		sectionsColor: ['#1bbc9b'],
	});
	
	function loadData(){
		$.ajax({
			contentType:"application/json;charset=UTF-8",
			url:"/yiyan",
			success:function(result){
				console.log(result);
				if (result.hitokoto !=="" ) {
					$("#yiyan").html(result.hitokoto)
					$("#from").html(result.from)
					$("#from_who").html(result.from_who)
				}
			},
			error:function(e){
				console.log(e);
			}
		});
	}
	
	window.setTimeout(loadData, 10000);
});
</script>

</head>

<body>

<div id="dowebok">
	<div class="section section1">
		<h3 style="text-align: center;">一言</h3>
		<p id="yiyan" style="text-align: center;">『{{ .Hitokoto }}』</p>
		<div style="text-align:right">—— <span id="from">{{ .From }}</span>「<span id="from_who"></span>」{{ .From_who }}</div>
		<div style="text-align:right;font-size: 20px;margin-top: 15px;">数据收集来源: <a href="https://hitokoto.cn/" target="_blank">https://hitokoto.cn/</a></div>
	</div>
</div>

</body>
</html>
`
