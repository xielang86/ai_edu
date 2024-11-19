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
	CreateStudentTable(mydao)
	CreateFileTable(mydao)
	CreateLessonTable(mydao)
	CreateTeacherTable(mydao)

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

	info := StudentBaseInfo{1, "xielang", 30, "lesson_1", "fake_lesson", 12345678,
		"lang.xie", "master", "cs", "cto", 1, "123456", "15110245219"}

	err := InsertStudentBasicInfo(mydao, info)
	if err != nil {
		fmt.Printf("fail! to insert student %s", err)
	}

	err = QueryStudent(mydao, "xielang", "123456", "15110245219", &info)
	if err != nil {
		fmt.Printf("fail! to query student %s", err)
	}
	fmt.Printf("query succ for user %s", info.Name)
}
