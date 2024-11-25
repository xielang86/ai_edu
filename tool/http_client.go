package main

import (
	"api_server/api"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SendPostRequest(url string, data []byte) error {
	// 将UserInfo结构体转换为JSON格式的字节切片

	// 创建一个新的HTTP POST请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// 设置请求头，指定Content-Type为application/json
	req.Header.Set("Content-Type", "application/json")

	// 创建HTTP客户端
	client := &http.Client{}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应体内容

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body failed")
		return err
	}
	fmt.Println("访问Web服务器得到的响应内容:", string(body))
	return err
}

func SendLoginPostRequest(url string, login_info api.LoginInfo) error {
	data, err := json.Marshal(login_info)
	if err != nil {
		return err
	}
	return SendPostRequest(url, data)
}

func main() {
	login_info := api.LoginInfo{Username: "xielang", Password: "123456"}
	err := SendLoginPostRequest("http://localhost:8080/login_post", login_info)
	if err != nil {
		fmt.Printf("err: %s", string(err.Error()))
	} else {
	}
}
