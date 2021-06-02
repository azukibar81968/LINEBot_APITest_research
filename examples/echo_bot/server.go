// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	bot, err := linebot.New(
		"bb1a39e0b86d7a0f1048281cabe57e0d",
		"KvZfSUoSY9+eAutIRaIH1LrdF8TrGaZU8DRhux4l+B2NdvsRGGhTcXs71VueDNoZU3eVvm3gV76vsDrcEdMBE0QIMAdUp2n9o6pO10OJIxbiypLc652newypoogxIYrn3AJWiYYrJkcE98cBM9Xa5wdB04t89/1O/w1cDnyilFU=",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					var replyMassage = tapTalkAPI(message.Text)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMassage)).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func tapTalkAPI(query string) string {
	fmt.Printf("*** Tapping API ***\n")
	url_target := "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk"
	args := url.Values{}
	args.Add("apikey", "DZZktPm2sWWq5NIjkAHYqnfXPjGujK4Q")
	args.Add("query", query)
	res, err := http.PostForm(url_target, args)
	if err != nil {
		fmt.Println("Request error:", err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Request error:", err)
		return ""
	}

	str_json := string(body)
	fmt.Println(str_json)
	var response = decodeTalkApiRes([]byte(str_json))

	fmt.Printf("*** finish ***\n")
	return response

}

type talkApiRes struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Results []talkApiResResult `json:"results"`
}
type talkApiResResult struct {
	Perplexity float32 `json:"perplexity"`
	Reply      string  `json:"reply"`
}

func decodeTalkApiRes(res []byte) string {
	// JSONデコード
	var decodedResponse talkApiRes
	if err := json.Unmarshal(res, &decodedResponse); err != nil {
		log.Fatal(err)
	}
	// デコードしたデータを表示
	fmt.Printf("reply : %s\n", decodedResponse.Results[0].Reply)
	return decodedResponse.Results[0].Reply
}
