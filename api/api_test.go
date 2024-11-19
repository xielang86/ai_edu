package api

import (
	"api_server/dao"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

var dsn string = "xielang:lang.xie86@(127.0.0.1:3306)/knowledge_edu"

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

func SendRegisterPostRequest(url string, info dao.StudentBaseInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return SendPostRequest(url, data)
}

func SendLoginPostRequest(url string, login_info LoginInfo) error {
	data, err := json.Marshal(login_info)
	if err != nil {
		return err
	}
	return SendPostRequest(url, data)
}

func TestRegisterAndLogin(t *testing.T) {
	// start a webserver
	var wg sync.WaitGroup
	wg.Add(1)

	// 在一个goroutine中启动Web服务器
	go func() {
		defer wg.Done()
		http.HandleFunc("/register", RegisterHandler)
		http.HandleFunc("/login", LoginHandler)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Web服务器启动失败:", err)
		}
	}()

	// 等待Web服务器启动完成
	go func() {
		for {
			_, err := http.Get("http://localhost:8080/")
			if err == nil {
				fmt.Println("Web服务器已成功启动，可以访问了")
				return
			}
			fmt.Println("正在等待Web服务器启动...")
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()
	// construct register data
	info := dao.StudentBaseInfo{Id: 1, Name: "xielang", Age: 30, LessonId: "lesson_1", LessonName: "fake_lesson", CreateTime: 12345678,
		ParentName: "lang.xie", ParentDegree: "master", ParentMajor: "cs", ParentCareer: "cto", Fee: 1, PassPhrase: "123456", Phone: "15110245219"}
	err := SendRegisterPostRequest("http://localhost:8080/register", info)
	if err != nil {
		fmt.Println("访问Web /register failed:", err)
	} else {
		fmt.Println("register succ:")
	}

	login_info := LoginInfo{"xielang", "123456", "15110245219", "123456"}
	err = SendLoginPostRequest("http://localhost:8080/login", login_info)
	if err != nil {
		fmt.Println("访问Web /loginfailed:", err)
	} else {
		fmt.Println("login succ:")
	}
}
