package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"regexp"
	"strconv"
)

const ( //存储mysql数据库连接信息
	USERNAME = "root"
	PASSWORD = "shr040726cc"
	HOST     = "127.0.0.1"
	PORT     = "3306"
	DBNAME   = "douban_movie"
)

var DB *sql.DB //定义类型是*sql.DB的全局变量DB，表示数据库连接池，后续所有数据库操作都通过这个变量进行

type MovieData struct { //定义结构体
	Title    string `json:"title"`
	Year     string `json:"year"`
	Director string `json:"director"`
	Picture  string `json:"picture"`
	Actor    string `json:"actor"`
	Score    string `json:"score"`
	Quote    string `json:"quote"`
}

func main() {
	InitDB()

	for i := 0; i < 10; i++ {
		fmt.Printf("正在爬取第 %d 页", i)
		Spider(strconv.Itoa(i * 25))
	}

}

func Spider(page string) {
	//1.发送请求
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250?start="+page, nil)

	if err != nil { //如果出错
		fmt.Println("req err", err)
	}
	//防止浏览器检测爬虫访问，所有添加一些请求头伪造成浏览器访问
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败", err)
	}
	defer resp.Body.Close() //关闭，防止资源泄露

	//2.解析响应
	//转换为可操作的文档
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析失败", err)
	}

	//3.获取节点信息
	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.hd > a > span:nth-child(1)
	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.hd > a > span:nth-child(1)
	//#content > div > div.article > ol > li:nth-child(1)
	//#content > div > div.article > ol > li:nth-child(1) > div > div.pic > a > img
	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > p:nth-child(1)
	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > div > span.rating_num
	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > p.quote > span

	//find根据css选择器找到目标html元素，each遍历元素执行自定义操作func
	docDetail.Find("#content > div > div.article > ol > li").
		Each(func(i int, s *goquery.Selection) { //s是当前节点的索引

			//存储元素，为插入数据库做准备
			var data MovieData

			title := s.Find("div > div.info > div.hd > a > span:nth-child(1)").Text() //提取文本内容：标题

			img := s.Find("div > div.pic > a > img") //找到img元素
			imgtem, ok := img.Attr("src")            //提取src

			info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text() //简洁（导演，主演，年份）

			score := s.Find("div > div.info > div.bd > div > span.rating_num").Text() //评分
			quote := s.Find("div > div.info > div.bd > p.quote > span").Text()        //引言

			if ok {
				director, actor, year := InfoSpite(info) //构造InfoSpite函数解析字段提取主演，导演，年份

				data.Title = title
				data.Picture = imgtem
				data.Score = score
				data.Quote = quote

				data.Actor = actor
				data.Director = director
				data.Year = year

				if InsertData(data) { //插入数据库

				} else {
					fmt.Println("插入失败")
					return
				}

			}
		})
	fmt.Println("插入成功")
	return
}

func InfoSpite(info string) (director, actor, year string) {
	//用正则表达式解析

	directorRe, _ := regexp.Compile(`导演: (.*)主演:`)
	director = string(directorRe.Find([]byte(info)))

	actorRe, _ := regexp.Compile(`主演: (.*)`)
	actor = string(actorRe.Find([]byte(info)))

	yearRe, _ := regexp.Compile(`(\d+)`)
	year = string(yearRe.Find([]byte(info)))
	return
}

func InitDB() { //初始化连接数据库
	//path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	path := USERNAME + ":" + PASSWORD + "@tcp(" + HOST + ":" + PORT + ")/" + DBNAME + "?charset=utf8"

	DB, _ = sql.Open("mysql", path) //打开数据库连接
	//设置了连接的最大存活时间为 10（单位应该是与具体实现相关，通常是秒），意味着一个数据库连接在经过 10 个单位时间后，如果还没有被使用，将会被关闭回收，这样可以避免长时间闲置的连接占用过多资源
	DB.SetConnMaxLifetime(10)
	//设定了最大空闲连接数为 5，也就是数据库连接池中最多允许存在 5 个处于空闲状态的连接，当空闲连接数超过这个值时，多余的空闲连接会被关闭，有助于合理控制资源使用和提升数据库连接性能
	DB.SetMaxIdleConns(5)

	//通过 DB.Ping() 语句来测试与数据库的连接是否可用
	//它会向数据库发送一个简单的测试请求（比如 MySQL 中类似 SELECT 1 的操作）
	//如果连接正常且数据库可以响应，则返回 nil ，否则返回一个包含错误信息的 error 类型值。
	if err := DB.Ping(); err != nil {
		fmt.Println("opon database fail")
		return
	}

	fmt.Println("connect success")
}

func InsertData(movieData MovieData) bool {

	//使用全局变量 DB 来创建一个数据库事务 tx
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("Begin err ", err)
		return false
	}

	//准备sol语句
	stmt, err := tx.Prepare("INSERT INTO movie_data(`Title`,`Director`,`Picture`,`Actor`,`Year`,`Score`,`Quote`)VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println("Prepare err ", err)
		return false
	}

	//stmt.Exec:通过预处理语句执行插入操作，将 MovieData 结构体中的字段逐个绑定到 SQL 语句中的占位符 ?
	_, err = stmt.Exec(movieData.Title, movieData.Director, movieData.Picture, movieData.Actor, movieData.Year, movieData.Score, movieData.Quote)
	if err != nil {
		fmt.Println("Exec err ", err)
		return false
	}

	//提交事务，将sql语句写入数据库
	_ = tx.Commit()

	return true
}
