package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ---------------------------------------------------------------
func main() {
	fmt.Printf("*** 開始 ***\n")
	url_target := "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk"
	args := url.Values{}
	args.Add("apikey", "DZZktPm2sWWq5NIjkAHYqnfXPjGujK4Q")
	args.Add("query", "こんにちは！")
	res, err := http.PostForm(url_target, args)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}

	str_json := string(body)
	fmt.Println(str_json)

	fmt.Printf("*** 終了 ***\n")
}
