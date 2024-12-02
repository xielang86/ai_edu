package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// just for recommmend model names
var kAllModelNames = []string{"Qwen/Qwen2.5-72B-Instruct", "Qwen/Qwen2.5-32B-Instruct"}

// in the future, build our model server
type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SiliconFlowReqBody struct {
	Model    string        `json:"model"`
	Messages [3]GPTMessage `json:"messages"`
}

type SiliconFlowRespChoice struct {
	Index         int        `json:"index"`
	Message       GPTMessage `json:"message"`
	Finish_reason string     `json:"finish_reason"`
}

type SiliconFlowRespBody struct {
	Id      string                  `json:"id"`
	Object  string                  `json:"object"`
	Created int64                   `json:"created"`
	Model   string                  `json:"model"`
	Choices []SiliconFlowRespChoice `json:"choices"`
}

func GetRawResp(req_body SiliconFlowReqBody, resp_body *SiliconFlowRespBody) error {
	url := "https://api.siliconflow.cn/v1/chat/completions"
	json_data, err := json.Marshal(req_body)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
		return err
	}

	// 创建一个POST请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return err
	}

	// 设置请求头，例如设置Content-Type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	// TODO(xl): mv to config
	req.Header.Set("authorization", "Bearer sk-acgzjedfzicxmxzuprvduikxfaoenzdmrkxyyyimxqvesppj")
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read resp failed:%v", err)
		return err
	}
	fmt.Printf("响应状态码: %d\n响应体: %s\n", resp.StatusCode, body)

	err = json.Unmarshal(body, resp_body)
	if err != nil {
		fmt.Println("unmarshal failed!:", err)
		return err
	}
	return nil
}

func GetSingleAnswer(query string, model_name string) string {
	req_body := SiliconFlowReqBody{
		Model:    model_name,
		Messages: [3]GPTMessage{},
	}
	req_body.Messages[0] = GPTMessage{"user", query}
	var resp_body SiliconFlowRespBody
	err := GetRawResp(req_body, &resp_body)
	if err != nil {
		fmt.Printf("failed to get raw resp for %v", req_body)
		return ""
	}
	if len(resp_body.Choices) == 0 {
		fmt.Printf("empty resp body for %v", req_body)
		return ""
	}
	return resp_body.Choices[0].Message.Content
}
