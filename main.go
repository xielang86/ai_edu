package main

import (
	"api_server/api"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/login", api.LoginHandler)
	http.HandleFunc("/register", api.RegisterHandler)
	http.HandleFunc("/register_post", api.RegisterPostHandler)
	http.HandleFunc("/login_post", api.LoginPostHandler)
	http.HandleFunc("/reset_pass", api.ResetpassHandler)
	http.HandleFunc("/reset_pass_post", api.ResetPassPostHandler)
	http.HandleFunc("/user_center", api.UserCenterHandler)
	http.HandleFunc("/file_center", api.FileCenterHandler)
	http.HandleFunc("/check_auth", api.CheckAuthHandler)
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("pages"))))
	fmt.Println("Server starting on port :8080...")
	http.ListenAndServe(":8080", nil)

	// router := gin.Default()

	// router.POST("/send_verification_code", api.SendVerificationCodeHandler)

	// // 验证验证码接口
	// router.POST("/verify_verification_code", api.VerificationCodeHandler)
	// router.Run(":8080")
}
