package api

import (
	"api_server/dao"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// LoginInfo 结构体用于接收客户端传来的登录请求信息
type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Code     string `json:"code"`
}

func IsUserExist(user dao.UserInfo, mydao *dao.UserDAO) int {
	var result dao.UserInfo
	dao.QueryUser(mydao, user.Name, user.Phone, &result)
	fmt.Printf("tow phone %s, %s", user.Phone, result.Phone)
	if result.Name == user.Name {
		return 1
	}
	if result.Phone == user.Phone {
		return 2
	}
	return 0
}

// SendVerificationCode模拟发送验证码，这里简单返回固定验证码示例，实际需调用短信服务等
func SendVerificationCode(phone string) string {
	// 实际场景可生成随机验证码并通过短信接口发送，这里返回固定的 "123456" 作为示例
	return "123456"
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

func PersonalHandler(w http.ResponseWriter, r *http.Request) {
	PageHandler(w, "./pages/personal_desc.html")
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 关闭请求体，释放资源
	defer r.Body.Close()

	// 解析JSON数据到User结构体
	var login_info LoginInfo
	err = json.Unmarshal(body, &login_info)
	fmt.Printf("Username: %s, Password: %s", login_info.Username, login_info.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	// 判断是用户名密码登录还是手机号验证码登录
	if !ValidLoginInfo(w, login_info) {
		return
	}
	var info dao.UserInfo
	err = dao.QueryUser(mydao, login_info.Username, login_info.Phone, &info)
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
		Token:   "",
	}

	w.Header().Set("Content-Type", "application/json")
	if login_info.Password != info.PassPhrase {
		responseData.Status = "failed"
		responseData.Message = "password error"
	} else {
		token, err := generateJWT(login_info.Username, info.Role)
		if err != nil {
			http.Error(w, "生成token失败", http.StatusInternalServerError)
			return
		}
		responseData.Token = token
		fmt.Printf("gen token %s\n", token)
	}

	PostResponse(w, responseData)
}

func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}
	var user dao.UserInfo
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
		insert_err := dao.InsertUserInfo(mydao, user)
		// 插入用户数据到数据库的SQL语句
		if insert_err != nil {
			responseData.Status = "failed"
			responseData.Message = "注册失败"
		}
	}

	PostResponse(w, responseData)
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

	PostResponse(w, responseData)
}

// 生成JWT的密钥，应该妥善保管，实际应用中可从配置文件等获取
var jwtSecret = []byte("your_secret_key")

func generateJWT(username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_str, err := token.SignedString([]byte(jwtSecret))
	fmt.Printf("gen jwt for user %s, %s\n", username, token_str)
	if err != nil {
		fmt.Printf("gen jwt failed user %s, %s", username, err.Error())
	}
	return token_str, err
}

func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	fmt.Printf("auth for token: %s\n", tokenString)
	if tokenString == "" {
		http.Error(w, "未提供认证token", http.StatusUnauthorized)
		return
	}
	if len(tokenString) >= 6 && tokenString[0:6] == "Bearer" {
		tokenString = tokenString[6:]
	}
	fmt.Printf("final token: %s\n", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		fmt.Printf("token err: %s", err.Error())
		http.Error(w, "token验证失败", http.StatusUnauthorized)
		return
	}

	responseData := ResponseData{
		Status:  "failed!",
		Message: "need to login",
		Token:   "",
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		role := claims["role"].(string)
		// NOTE(*): trick
		responseData.Message = fmt.Sprintf("%s,%s", username, role)
		responseData.Status = "success"
	} else {
		http.Error(w, "token无效", http.StatusUnauthorized)
	}
	PostResponse(w, responseData)
}

func GetAllTeacher(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	// TODO(xl): use dao pool to opt the performance
	// connect to db
	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	responseData := ResponseData{
		Status:  "success",
		Message: "",
	}

	var result []dao.UserInfo
	err := dao.QueryAllUser(mydao, "teacher", &result)
	if err != nil {
		responseData.Status = "failed"
		responseData.Message = "failed to get all teachers"
	} else {
		responseData.Data = result
	}
	PostResponse(w, responseData)
}
