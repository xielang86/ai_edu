package api

import (
	"api_server/dao"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "github.com/tencentcloud/cos-go-sdk-v5"
)

type ProjectFilter struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	NeedAll  int    `json:"need_all"`
}

type LessonInfo struct {
	Name   string `json:"name"`
	Joined int    `json:"joined"`
}

func GetAllProjectHandler(w http.ResponseWriter, r *http.Request) {
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
	var filter ProjectFilter
	err = json.Unmarshal(body, &filter)
	if err != nil {
		fmt.Printf("err filter info %s , would get all project", err)
	}

	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	responseData := ResponseData{
		Status:  "success",
		Message: "操作成功",
	}

	var data []LessonInfo
	var all_lessons []dao.LessonInfo
	if filter.NeedAll > 0 {
		err = dao.QueryAllActiveLesson(mydao, &all_lessons)
		if err != nil {
			fmt.Printf("get all active lesson failed!")
		} else {
			fmt.Printf("get all active lesson succ! %d", len(all_lessons))
		}
	}

	var result []string
	err = dao.QueryAllLessonNameByUsername(mydao, filter.Username, &result)
	less_dict := make(map[string]int)
	for _, r := range result {
		data = append(data, LessonInfo{Name: r, Joined: 1})
		less_dict[r] = 1
	}
	for _, lesson := range all_lessons {
		// 查询元素
		_, ok := less_dict[lesson.Name]
		if ok {
			fmt.Printf("dup lesson%s\n", lesson.Name)
		} else {
			data = append(data, LessonInfo{Name: lesson.Name, Joined: 0})
		}
	}

	if err != nil {
		responseData.Status = "failed"
		responseData.Message = "get projecdts failed"
	} else {
		responseData.Data = data
	}

	PostResponse(w, responseData)
}

func GetAllStudentHandler(w http.ResponseWriter, r *http.Request) {
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
	var filter ProjectFilter
	err = json.Unmarshal(body, &filter)
	if err != nil {
		fmt.Printf("err filter info %s , would get all project", err)
		fmt.Println(body)
	}

	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	responseData := ResponseData{
		Status:  "success",
		Message: "操作成功",
	}

	var result []string
	err = dao.QueryAllStudentNameByTeacher(mydao, filter.Username, &result)
	var data []dao.UserInfo
	for _, r := range result {
		data = append(data, dao.UserInfo{Name: r})
	}

	if err != nil {
		responseData.Status = "failed"
		responseData.Message = "get students failed"
	} else {
		responseData.Data = data
	}

	PostResponse(w, responseData)
}
