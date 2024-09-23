package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestModelApi(t *testing.T) {
	req_body := SiliconFlowReqBody{
		Model:    "alibaba/Qwen1.5-110B-Chat",
		Messages: [3]GPTMessage{},
	}
	req_body.Messages[0] = GPTMessage{"user", "如何走出情劫"}
	var resp_body SiliconFlowRespBody
	err := GetRawResp(req_body, &resp_body)
	if err != nil {
		t.Errorf("failed to test raw resp for %v", req_body)
	}
	fmt.Println(resp_body)
}

func TestSingleAnswer(t *testing.T) {
	var query = "如何走出情劫"
	var ans = GetSingleAnswer(query, "Qwen/Qwen2-7B-Instruct")
	if len(ans) == 0 {
		t.Errorf("empty result for %s", query)
	}
	fmt.Println(ans)
}

func TestSiliconFlow(t *testing.T) {
	// 目标URL
	url := "https://api.siliconflow.cn/v1/chat/completions"

	req_body := SiliconFlowReqBody{
		Model:    "alibaba/Qwen1.5-110B-Chat",
		Messages: [3]GPTMessage{},
	}
	req_body.Messages[0] = GPTMessage{"user", "如何走出情劫"}

	json_data, err := json.Marshal(req_body)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
		return
	}

	// data := []byte(`{"model":"alibaba/Qwen1.5-110B-Chat", "messages":[{"role": "user", "content": "抛砖引玉是什么意思呀"}]}`)
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))

	// 创建一个POST请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 设置请求头，例如设置Content-Type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", "Bearer sk-acgzjedfzicxmxzuprvduikxfaoenzdmrkxyyyimxqvesppj")
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}
	fmt.Printf("响应状态码: %d\n响应体: %s\n", resp.StatusCode, body)

	var resp_body SiliconFlowRespBody
	json.Unmarshal(body, &resp_body)
	if len(resp_body.Choices) > 0 {
		fmt.Printf("answer: %s\n", resp_body.Choices[0].Message.Content)
	}
}
