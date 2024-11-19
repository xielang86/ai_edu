package api

import (
	"api_server/dao"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

var kEduKnowledgeDB string = "xielang:lang.xie86@(127.0.0.1:3306)/knowledge_edu"

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	// get knowledge
	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	var query string = "select level4, level5 from knowledge_edu.en_knowledge_point where level4=?"
	result := &dao.EduKnowledge{}
	dao.CreateKnowledgeTable(mydao)
	dao.QueryEduKnowledge(mydao, query, result)
	dao.CloseDB(mydao)

	fmt.Println("Hello, World!")
}

// RegisterHandler处理用户注册的请求
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}
	var user dao.StudentBaseInfo
	// 解析请求体中的JSON数据到user结构体
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "解析请求数据失败", http.StatusBadRequest)
		return
	}
	// connect to db
	// TODO(xl): use dao pool to opt the performance
	var dsn string = "xielang:lang.xie86@(127.0.0.1:3306)/knowledge_edu"
	mydao := dao.NewUserDAO(nil, dsn)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	insert_err := dao.InsertStudentBasicInfo(mydao, user)
	// 插入用户数据到数据库的SQL语句
	if insert_err != nil {
		http.Error(w, "插入用户数据失败", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "注册成功")
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

// LoginHandler处理登录请求
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}
	var login_info LoginInfo
	err := json.NewDecoder(r.Body).Decode(&login_info)
	if err != nil {
		http.Error(w, "解析登录请求数据失败", http.StatusBadRequest)
		return
	}

	var dsn string = "xielang:lang.xie86@(127.0.0.1:3306)/knowledge_edu"
	mydao := dao.NewUserDAO(nil, dsn)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	// 判断是用户名密码登录还是手机号验证码登录
	if login_info.Phone != "" && login_info.Code != "" {
		// 先验证验证码是否正确（这里只是简单对比示例，实际可能需更复杂逻辑和验证有效期等）
		sentCode := SendVerificationCode(login_info.Phone)
		if login_info.Code != sentCode {
			http.Error(w, "验证码错误", http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(w, "登录信息不完整", http.StatusBadRequest)
		return
	}
	var info dao.StudentBaseInfo
	err = dao.QueryStudent(mydao, login_info.Username, login_info.Password, login_info.Phone, &info)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "用户不存在或登录信息错误", http.StatusUnauthorized)
		} else {
			http.Error(w, "数据库查询错误", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "登录成功")
}
