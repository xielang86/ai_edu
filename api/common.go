package api

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

var kEduKnowledgeDB string = "xielang:lang.xie86@(127.0.0.1:3306)/knowledge_edu"

type ResponseData struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Token   string      `json:"token"`
	Data    interface{} `json:"data"`
}

func PageHandler(w http.ResponseWriter, page_path string) {
	t, err := template.ParseFiles(page_path)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func PostResponse(w http.ResponseWriter, responseData ResponseData) {
	// 将结构体转换为JSON格式的字节切片
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		// 如果转换出错，返回500 Internal Server Error状态码及错误信息
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		// 如果转换出错，返回500 Internal Server Error状态码及错误信息
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
