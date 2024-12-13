package main

import (
	"api_server/api"
	"fmt"
	"html/template"
	"net/http"
)

func HandleAllPage(dir_path string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 获取请求路径，去除开头的"/"
		path := r.URL.Path[1:]

		if path == "" {
			path = "index" // 如果请求路径为空，默认返回index.html页面（可根据需求调整）
		}
		// 拼接完整的文件路径，假设HTML文件都放在pages目录下
		page_path := fmt.Sprintf("%s/%s.html", dir_path, path)

		t, err := template.ParseFiles(page_path)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
	// for js,css,image
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("pages"))))
	http.Handle("/local_files/", http.StripPrefix("/local_files/", http.FileServer(http.Dir("local_files"))))
}

func main() {
	HandleAllPage("./pages")

	http.HandleFunc("/register_post", api.RegisterPostHandler)
	http.HandleFunc("/login_post", api.LoginPostHandler)
	http.HandleFunc("/reset_pass_post", api.ResetPassPostHandler)

	http.HandleFunc("/upload", api.UploadHandler)
	// http.HandleFunc("/upload_ocr", api.UploadAndOcrHandler)
	http.HandleFunc("/get_all_file", api.GetAllFileHandler)

	http.HandleFunc("/get_all_project", api.GetAllProjectHandler)

	http.HandleFunc("/get_all_student", api.GetAllStudentHandler)
	http.HandleFunc("/get_all_teacher", api.GetAllTeacher)
	http.HandleFunc("/check_auth", api.CheckAuthHandler)

	http.HandleFunc("/verify_code", api.VerifyCodeHandler)
	http.HandleFunc("/send_verify_code", api.SendSMSHandler)

	fmt.Println("Server starting on port :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
