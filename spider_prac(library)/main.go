package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	"strings"
)

var DB *sql.DB

const (
	USERNAME = "root"
	PASSWORD = "shr040726cc"
	HOST     = "127.0.0.1"
	PORT     = "3306"
	DBNAME   = "library"
)

type BookData struct {
	Title   string
	Rating  int
	Price   string
	Picture string
	Status  string
}

func main() {
	initDB()

	for i := 0; i < 50; i++ {
		fmt.Printf("正在爬取第 %d 页", i)
		spider(strconv.Itoa(i * 20))
	}

}

func spider(page string) {
	client := http.Client{}
	//构造请求
	req, err := http.NewRequest("GET", "http://books.toscrape.com/?start="+page, nil)
	if err != nil {
		fmt.Println("构造错误: ", err)
	}
	//添加请求头
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败: ", err)
	}
	defer resp.Body.Close()

	//解析响应

	//转换可操作文档
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析失败：", err)
	}

	//获取节点信息
	//#default > div > div > div > div > section > div:nth-child(2) > ol > li:nth-child(1)
	//#default > div > div > div > div > section > div:nth-child(2) > ol > li:nth-child(1) > article > h3 > a
	//#default > div > div > div > div > section > div:nth-child(2) > ol > li:nth-child(1) > article > p
	//#default > div > div > div > div > section > div:nth-child(2) > ol > li:nth-child(1) > article > div.image_container > a > img
	//#default > div > div > div > div > section > div:nth-child(2) > ol > li:nth-child(1) > article > div.product_price > p.price_color
	//#default > div > div > div > div > section > div:nth-child(2) > ol > li:nth-child(1) > article > div.product_price > p.instock.availability
	docDetail.Find("#default > div > div > div > div > section > div:nth-child(2) > ol > li").
		Each(func(i int, s *goquery.Selection) {
			var bookdata BookData

			title := s.Find("article > h3 > a").Text()
			//score := s.Find("article > p").Text()
			img := s.Find("article > div.image_container > a > img")
			imgtem, ok := img.Attr("src")
			price := s.Find("article > div.product_price > p.price_color").Text()
			//status := s.Find("article > div.product_price > p.instock.availability").Text()
			rawStatus := s.Find("article > div.product_price > p.instock.availability").Text()
			status := strings.TrimSpace(rawStatus) // 去掉前后的空格和换行符

			scoreClass, exists := s.Find("article > p").Attr("class")
			score := "Unknown"
			if exists {
				// scoreClass 可能是 "star-rating Three" 或类似字符串
				scoreParts := strings.Split(scoreClass, " ")
				if len(scoreParts) > 1 {
					score = scoreParts[1] // 提取 "Three"
				}
			}

			// 将评分转换为数字存储
			rating := 0
			switch score {
			case "One":
				rating = 1
			case "Two":
				rating = 2
			case "Three":
				rating = 3
			case "Four":
				rating = 4
			case "Five":
				rating = 5
			default:
				rating = 0 // 如果没有匹配到
			}

			if ok {
				//fmt.Println("title:", title)
				//fmt.Println("rating:", rating)
				//fmt.Println("imgtem:", imgtem)
				//fmt.Println("price:", price)
				//fmt.Println("status:", status)
				//return
				bookdata.Title = title
				bookdata.Rating = rating
				bookdata.Price = price
				bookdata.Status = status
				bookdata.Picture = imgtem
				//fmt.Println("TEST6666")
				//fmt.Println(bookdata.Title)
				//fmt.Println(bookdata.Rating)
				//fmt.Println(bookdata.Price)
				//fmt.Println(bookdata.Status)
				if InsertData(bookdata) {

				} else {
					fmt.Println("插入失败")
				}
			}
		})
	fmt.Println("插入成功")
	return
}

func initDB() {
	path := USERNAME + ":" + PASSWORD + "@tcp(" + HOST + ":" + PORT + ")/" + DBNAME + "?charset=utf8"
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(10)
	DB.SetMaxOpenConns(5)

	if err := DB.Ping(); err != nil {
		fmt.Println("open database err:", err)
		return
	}
	fmt.Println("connect success")
}

func InsertData(bookdata BookData) bool {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("Begin Err:", err)
		return false
	}
	stmt, err := tx.Prepare("INSERT INTO library_data(`Title`,`Rating`,`Price`,`Picture`,`Status`)VALUES(?,?,?,?,?)")
	if err != nil {
		fmt.Println("Prepare Err:", err)
		return false
	}
	_, err = stmt.Exec(bookdata.Title, bookdata.Rating, bookdata.Price, bookdata.Picture, bookdata.Status)
	if err != nil {
		fmt.Println("Exec Err:", err)
		return false
	}
	_ = tx.Commit()
	return true
}
