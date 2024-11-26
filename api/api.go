package api

import (
	"api_server/dao"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"
)

var kEduKnowledgeDB string = "xielang:lang.xie86@(127.0.0.1:3306)/knowledge_edu"

type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func IsUserExist(user dao.StudentBaseInfo, mydao *dao.UserDAO) int {
	var result dao.StudentBaseInfo
	dao.QueryStudent(mydao, user.Name, user.Phone, &result)
	fmt.Printf("tow phone %s, %s", user.Phone, result.Phone)
	if result.Name == user.Name {
		return 1
	}
	if result.Phone == user.Phone {
		return 2
	}
	return 0
}

func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}
	var user dao.StudentBaseInfo
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 关闭请求体，释放资源
	defer r.Body.Close()

	// 解析JSON数据到User结构体
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO(xl): use dao pool to opt the performance
	// connect to db
	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	responseData := ResponseData{
		Status:  "success",
		Message: "操作成功",
	}

	st := IsUserExist(user, mydao)
	fmt.Printf("use exit code: %d", st)
	if st == 1 {
		responseData.Status = "failed"
		responseData.Message = "username exist"
	} else if st == 2 {
		responseData.Status = "failed"
		responseData.Message = "phone exist"
	} else {
		// before insert , set some value
		user.CreateTime = time.Now().Unix()
		user.Fee = 0
		user.Age = 0
		fmt.Printf("insert user name=%s and pass=%s, creat_time=%d, caree=%s, phone=%s",
			user.Name, user.PassPhrase, user.CreateTime, user.ParentCareer, user.Phone)
		insert_err := dao.InsertStudentBasicInfo(mydao, user)
		// 插入用户数据到数据库的SQL语句
		if insert_err != nil {
			responseData.Status = "failed"
			responseData.Message = "注册失败"
		}
	}

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

// LoginInfo 结构体用于接收客户端传来的登录请求信息
type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Code     string `json:"code"`
}

// SendVerificationCode模拟发送验证码，这里简单返回固定验证码示例，实际需调用短信服务等
func SendVerificationCode(phone string) string {
	// 实际场景可生成随机验证码并通过短信接口发送，这里返回固定的 "123456" 作为示例
	return "123456"
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	PageHandler(w, "./pages/login.html")
}

// RegisterHandler处理用户注册的请求
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	PageHandler(w, "./pages/register.html")
}

func ResetpassHandler(w http.ResponseWriter, r *http.Request) {
	PageHandler(w, "./pages/reset_pass.html")
}

func UserCenterHandler(w http.ResponseWriter, r *http.Request) {
	PageHandler(w, "./pages/user_center.html")
}

func ValidLoginInfo(w http.ResponseWriter, login_info LoginInfo) bool {
	if login_info.Phone != "" && login_info.Code != "" {
		// 先验证验证码是否正确（这里只是简单对比示例，实际可能需更复杂逻辑和验证有效期等）
		sentCode := SendVerificationCode(login_info.Phone)
		if login_info.Code != sentCode {
			http.Error(w, "验证码错误", http.StatusUnauthorized)
			return false
		}
	} else if len(login_info.Username) < 4 || len(login_info.Password) < 6 {
		http.Error(w, "username or passphrase not suitable", http.StatusBadRequest)
		return false
	}
	return true
}

// LoginPostHandler处理登录请求
func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	var login_info LoginInfo
	r.ParseForm()
	login_info.Username = r.FormValue("username")
	login_info.Password = r.FormValue("password")
	fmt.Printf("Username: %s, Password: %s", login_info.Username, login_info.Password)

	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	// 判断是用户名密码登录还是手机号验证码登录
	if !ValidLoginInfo(w, login_info) {
		return
	}
	var info dao.StudentBaseInfo
	err := dao.QueryStudent(mydao, login_info.Username, login_info.Phone, &info)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "用户不存在或登录信息错误", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	responseData := ResponseData{
		Status:  "success",
		Message: "login success",
	}

	if login_info.Password != info.PassPhrase {
		http.Error(w, "password falied", http.StatusUnauthorized)
		return
	}

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
	}
}

func ResetPassPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}
	var info LoginInfo
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 关闭请求体，释放资源
	defer r.Body.Close()

	// 解析JSON数据到User结构体
	err = json.Unmarshal(body, &info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO(xl): use dao pool to opt the performance
	// connect to db
	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	responseData := ResponseData{
		Status:  "success",
		Message: "操作成功",
	}

	// before insert , set some value
	err = dao.ModifyPassphrase(mydao, info.Username, info.Phone, info.Password)
	// 插入用户数据到数据库的SQL语句
	if err != nil {
		responseData.Status = "failed"
		responseData.Message = "reset failed"
	}

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
