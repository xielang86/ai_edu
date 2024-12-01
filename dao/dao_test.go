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

	var query string = "select * from knowledge_edu.student_lesson_process"
	result := &StudentLessonProcess{}
	QueryStudentLessonProcess(mydao, query, result)

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
	// insert a fake teacher
	info := UserInfo{1, "xielang_teacher", "teacher", 30, "", "", 12345678,
		"parent_name", "phd", "cs", "professor", "thu", 1, "12345678", "15110245219", "just a teacher"}

	err := QueryUser(mydao, info.Name, "", &info)
	if err != nil {
		InsertUserInfo(mydao, info)
	}

	// insert two lesson for the teacher
	err = AddLesson(mydao, "c++", "xielang_teacher", "")
	err = AddLesson(mydao, "concrete_math", "xielang_teacher", "")
	if err != nil {
		fmt.Printf("fail! to add lesson for teacher %s", err)
	}

	// add lesson for two student
	err = AddLessonForStudent(mydao, "xielang", "c++")
	err = AddLessonForStudent(mydao, "xielang", "concrete math")
	err = AddLessonForStudent(mydao, "fuyun", "concrete_math")
	if err != nil {
		fmt.Printf("fail! to add lesson for student %s", err)
	}

	// test for query all lesson for student
	students := []string{"xielang", "fuyun"}
	for _, student := range students {
		var result []string
		err = QueryAllLessonNameByUsername(mydao, student, &result)
		if err != nil {
			fmt.Printf("fail! query all less for user %s, %s", student, err)
		} else {
			fmt.Printf("find %d lesson for user %s", len(result), student)
		}
	}
}
