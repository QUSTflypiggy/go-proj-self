package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//
//https://api.bilibili.com/x/v2/reply/wbi/main?oid=251119469&type=1&mode=2&pagination_str=%7B%22offset%22:%22%22%7D&plat=1&seek_rpid=&web_location=1315875&w_rid=f01364da5ac363ac8ac9ff5fee0d327b&wts=1734609730

type KingRankResp struct {
	Code int64 `json:"code"`
	Data struct {
		Replies []struct {
			Content struct {
				Device  string        `json:"device"`
				JumpURL struct{}      `json:"jump_url"`
				MaxLine int64         `json:"max_line"`
				Members []interface{} `json:"members"`
				Message string        `json:"message"`
				Plat    int64         `json:"plat"`
			} `json:"content"`
			Count  int64 `json:"count"`
			Folder struct {
				HasFolded bool   `json:"has_folded"`
				IsFolded  bool   `json:"is_folded"`
				Rule      string `json:"rule"`
			} `json:"folder"`
			Like    int64 `json:"like"`
			Replies []struct {
				Action  int64 `json:"action"`
				Assist  int64 `json:"assist"`
				Attr    int64 `json:"attr"`
				Content struct {
					Device  string   `json:"device"`
					JumpURL struct{} `json:"jump_url"`
					MaxLine int64    `json:"max_line"`
					Message string   `json:"message"`
					Plat    int64    `json:"plat"`
				} `json:"content"`
				Rcount  int64       `json:"rcount"`
				Replies interface{} `json:"replies"`
			} `json:"replies"`
			Type int64 `json:"type"`
		} `json:"replies"`
	} `json:"data"`
	Message string `json:"message"`
}

func main() {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/v2/reply/wbi/main?oid=251119469&type=1&mode=2&pagination_str=%7B%22offset%22:%22%22%7D&plat=1&seek_rpid=&web_location=1315875&w_rid=f01364da5ac363ac8ac9ff5fee0d327b&wts=1734609730", nil)
	if err != nil {
		fmt.Println("req err ", err)
	}
	req.Header.Set("user-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	resq, err := client.Do(req)
	bodyText, err := ioutil.ReadAll(resq.Body)
	if err != nil {
		fmt.Println("io err ", err)
	}

	var resultList KingRankResp
	_ = json.Unmarshal(bodyText, &resultList)
	for _, v := range resultList.Data.Replies {
		fmt.Println("一级评论 ", v.Content.Message)
		for _, r := range v.Replies {
			fmt.Println("二级评论 ", r.Content.Message)
		}
	}
}
