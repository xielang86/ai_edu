package api

import (
	"net/http"
	// "github.com/tencentcloud/cos-go-sdk-v5"
)

func ProjectHandler(w http.ResponseWriter, r *http.Request) {
	PageHandler(w, "./pages/project.html")
}

func ProjectPostHandler(w http.ResponseWriter, r *http.Request) {
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
