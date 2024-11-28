package api

import "net/http"

func FileCenterHandler(w http.ResponseWriter, r *http.Request) {
	PageHandler(w, "./pages/file_center.html")
}

func FileCenterPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}
	responseData := ResponseData{
		Status:  "success",
		Message: "操作成功",
	}

	PostResponse(w, responseData)
}
