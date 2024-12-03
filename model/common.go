package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// TODO(*): use http pool
var http_client *http.Client

func init() {
	transport := &http.Transport{
		// 设置最大空闲连接数，比如设置为10
		MaxIdleConns: 16,
		// 设置每个主机的最大连接数，例如设置为5
		MaxConnsPerHost: 8,
		// 设置空闲连接超时时间，这里设置为30秒
		IdleConnTimeout: 30 * time.Second,
	}

	// 使用自定义的Transport创建http.Client实例
	http_client = &http.Client{
		Transport: transport,
	}
}

func ApplicationJsonRequest(your_url string, access_key_id string, access_key_secret string, body_params interface{}, url_params interface{}) (interface{}, error) {
	// 将 Body Params 转换为 JSON 格式
	body, err := json.Marshal(body_params)
	if err != nil {
		return "", fmt.Errorf("无法将 Body Params 转换为 JSON: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest(http.MethodPost, your_url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("无法创建 HTTP 请求: %v", err)
	}

	// 设置请求头
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Access-Key-ID", access_key_id)
	req.Header.Add("Access-Key-Secret", access_key_secret)

	// 添加 URL 参数
	if url_params != nil {
		query := url.Values{}
		for key, value := range url_params.(map[string]string) {
			query.Add(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}

	// 发送请求
	resp, err := http_client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送 HTTP 请求时出错: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read file failed %s", err)
	}
	return responseBody, err
}
