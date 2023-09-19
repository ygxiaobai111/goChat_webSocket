package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ChatCompletion struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int64  `json:"created"`
	Result           string `json:"result"`
	IsTruncated      bool   `json:"is_truncated"`
	NeedClearHistory bool   `json:"need_clear_history"`
	Usage            struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func chat(me string) (ans string, err error) {

	url := "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/eb-instant?access_token=24.3eb2307a96c01396ed1f933a4bfe89d5.2592000.1697374262.282335-39389195/RTqD6bWGy/uWXhE/7c7VcXQmnybhEbso14zkJw=="
	payload := strings.NewReader(`{"messages":[]}`) // 初始化为空的messages数组
	client := &http.Client{}

	// 将新的对话信息添加到payload中，并包含过去的对话信息
	payload.Reset(`{"messages":[{"role":"user","content":"` + me + `"}]}`)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(string(body))
	var chatA ChatCompletion
	json.Unmarshal(body, &chatA)
	return chatA.Result, nil

}
