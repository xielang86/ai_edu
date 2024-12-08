package dao

import (
	"fmt"
	"testing"
)

var dsn string = "xielang:lang.xie86@(127.0.0.1:3306)/knowledge_edu"

func TestDAOCreate(t *testing.T) {
	mydao := NewUserDAO(nil, dsn)
	ConnectDB(mydao)
	defer CloseDB(mydao)

	CreateStudentLessonProcessTable(mydao)
	CreateUserInfoTable(mydao)
	CreateFileTable(mydao)
	CreateLessonTable(mydao)

	fmt.Println("test table!")
}

func TestQueryStudent(t *testing.T) {
	// insert a fake
	mydao := NewUserDAO(nil, dsn)
	ConnectDB(mydao)
	defer CloseDB(mydao)

	info := UserInfo{1, "xielang", "student", 30, "lesson_1", "fake_lesson", 12345678,
		"lang.xie", "master", "cs", "cto", "seu", 1, "12345678", "15110245219", "just a geek"}

	err := QueryUser(mydao, info.Name, "", &info)
	if err != nil {
		err = InsertUserInfo(mydao, info)
	}
	if err != nil {
		fmt.Printf("fail! to insert student %s", err)
	}

	err = QueryUser(mydao, "xielang", "15110245219", &info)
	if err != nil {
		fmt.Printf("fail! to query student %s", err)
	}
	fmt.Printf("query succ for user %s", info.Name)

	// change name for another studeng
	info.Name = "fuyun"
	info.Phone = "12345678910"
	info.ParentName = "fujian"
	err = QueryUser(mydao, info.Name, "", &info)
	if err != nil {
		err = InsertUserInfo(mydao, info)
	}
	if err != nil {
		fmt.Printf("fail! to insert student %s", err)
	}

	err = QueryUser(mydao, "fuyun", "", &info)
	if err != nil {
		fmt.Printf("fail! to query student %s", err)
	}
	fmt.Printf("query succ for user %s", info.Name)
}

func TestQueryLesson(t *testing.T) {
	mydao := NewUserDAO(nil, dsn)
	ConnectDB(mydao)
	defer CloseDB(mydao)

	teachers := [3]UserInfo{
		{Name: "黄向东", Role: "teacher", Age: 38, Phone: "12345678910"},
		{Name: "刘英博", Role: "teacher", Age: 58, Phone: "22345678910"},
		{Name: "何伟", Role: "teacher", Age: 48, Phone: "32345678910"},
	}
	for _, teacher := range teachers {
		var user UserInfo
		err := QueryUser(mydao, teacher.Name, "", &user)
		if err != nil {
			InsertUserInfo(mydao, teacher)
		}
	}

	// insert two lesson for the teacher
	lessons := []LessonInfo{
		{Name: "人工智能在供应链管理中的应用", TeacherName: "黄向东"},
		{Name: "人工智能在股票市场预测中的应用", TeacherName: "刘英博"},
		{Name: "人工智能在人才招聘和管理中的应用", TeacherName: "刘英博"},
		{Name: "人工智能在社交媒体营销中的应用", TeacherName: "何伟"},
		{Name: "AI驱动的项目管理工具的效率研究", TeacherName: "何伟"},
	}
	for _, lesson := range lessons {
		err := AddLesson(mydao, lesson.Name, lesson.TeacherName, "")
		if err != nil {
			fmt.Printf("fail! to add lesson for teacher %s", err)
		}
	}

	type StudentLesson struct {
		Name       string
		LessonName string
	}

	selections := []StudentLesson{
		{Name: "xielang", LessonName: "人工智能在供应链管理中的应用"},
		{Name: "xielang", LessonName: "人工智能在股票市场预测中的应用"},
		{Name: "fuyun", LessonName: "人工智能在社交媒体营销中的应用"},
		{Name: "fuyun", LessonName: "基于人工智能结构健康检测系统"},
	}
	// add lesson for two student
	for _, a := range selections {
		err := AddLessonForStudent(mydao, a.Name, a.LessonName)
		if err != nil {
			fmt.Printf("fail! to add lesson for student %s", err)
		}
	}

	// test for query all lesson for student
	students := []string{"xielang", "fuyun"}
	for _, student := range students {
		var result []string
		err := QueryAllLessonNameByUsername(mydao, student, &result)
		if err != nil {
			fmt.Printf("fail! query all less for user %s, %s", student, err)
		} else {
			fmt.Printf("find %d lesson for user %s", len(result), student)
		}
	}
}
