package dao

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type LessonInfo struct {
	Id                  int64
	Name                string `json:"name"`
	TeacherId           int64
	TeacherName         string `json:"teacher_name"`
	InvolvedStudentId   string
	InvolvedStudentName string `json:"involved_student_ name"`
	InitFileList        string `json:"init_file_list"`
	CreateTime          int64
}

type StudentLessonProcess struct {
	StudentId   int64
	StudentName string `json:"student_name"`
	LessonId    int64
	LessonName  string `json:"lesson_name"`
	FeeInfo     string `json:"fee_info"`
	FileList    string `json:"file_list"`
	ProcessInfo string `json:"process_info"`
	CreateTime  int64
}

func CreateLessonTable(dao *UserDAO) {
	var sql_str string = `
	CREATE TABLE IF NOT EXISTS lesson_info (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    teacher_id BIGINT,
    teacher_name VARCHAR(255),
    involved_student_id VARCHAR(1024),
    involved_student_name VARCHAR(1024),
    init_file_list VARCHAR(255),
    create_time BIGINT);
	`
	CreateTable(dao, sql_str)
}

func CreateStudentLessonProcessTable(dao *UserDAO) {
	var sql_str string = `
	CREATE TABLE IF NOT EXISTS student_lesson_process (
    student_id BIGINT,
    student_name VARCHAR(64),
    lesson_id BIGINT,
    lesson_name VARCHAR(128),
    fee_info	VARCHAR(64),
		file_list	VARCHAR(255),
		process_info VARCHAR(255),
    create_time BIGINT,
		PRIMARY KEY (student_id, lesson_id)
	);
	`
	CreateTable(dao, sql_str)
}

func QueryLessonByName(dao *UserDAO, lesson_name string, info *LessonInfo) error {
	query := fmt.Sprintf("select id,name,teacher_id,teacher_name,involved_student_id,involved_student_name,init_file_list from lesson_info where name=\"%s\"", lesson_name)
	row := dao.db.QueryRow(query)
	err := row.Scan(&info.Id, &info.Name, &info.TeacherId, &info.TeacherName, &info.InvolvedStudentId, &info.InvolvedStudentName, &info.InitFileList)
	if err != nil {
		fmt.Printf("query lesson failed: %s, query=%s\n", err, query)
		return err
	}
	return nil
}

func QueryAllLessonNameByUsername(dao *UserDAO, username string, result *[]string) error {
	query := fmt.Sprintf("select lesson_name from student_lesson_process where student_name=%s", username)
	rows, err := dao.db.Query(query)
	if err != nil {
		fmt.Println("查询错误:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var lesson_name string
		err := rows.Scan(&lesson_name)
		if err != nil {
			fmt.Println("映射错误:", err)
			continue
		}
		// 将结构体添加到数组
		*result = append(*result, lesson_name)
	}

	return nil
}

func QueryAllLessonByUserId(dao *UserDAO, student_id int64, result *[]string) error {
	// query := fmt.Sprintf("select student_id, student_name, lesson_id, lesson_name, file_list,process_info from student_lesson_process where student_id=%d",
	query := fmt.Sprintf("select lesson_name from student_lesson_process where student_id=%d",
		student_id)
	rows, err := dao.db.Query(query)
	if err != nil {
		fmt.Println("查询错误:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var lesson_name string
		err := rows.Scan(&lesson_name)
		if err != nil {
			fmt.Println("映射错误:", err)
			continue
		}
		// 将结构体添加到数组
		*result = append(*result, lesson_name)
	}
	return nil
}

func AddLesson(dao *UserDAO, name string, teacher_name string, init_file_list string) error {
	var info LessonInfo
	err := QueryLessonByName(dao, name, &info)
	if err == nil {
		return fmt.Errorf("lesson %s exits", name)
	}
	var user UserInfo
	err = QueryUser(dao, teacher_name, "", &user)
	if err != nil {
		fmt.Printf("add lesson %s failed for user %s", name, teacher_name)
		return err
	}
	create_time := time.Now().Unix()
	insert_sql := `INSERT INTO lesson_info
	(name, teacher_id, teacher_name,involved_student_id,involved_student_name,init_file_list,create_time)
	VALUES (?,?,?,?,?,?,?)`
	_, insert_err := dao.db.Exec(insert_sql, name, user.Id, teacher_name, "", "", init_file_list, create_time)
	return insert_err
}

func AddLessonForStudent(dao *UserDAO, student_name string, lesson_name string) error {
	var user UserInfo
	err := QueryUser(dao, student_name, "", &user)
	if err != nil {
		fmt.Printf("add lesson %s failed for student%s, because of query user", lesson_name, student_name)
		return err
	}
	var lesson_info LessonInfo
	err = QueryLessonByName(dao, lesson_name, &lesson_info)
	if err != nil {
		fmt.Printf("add lesson %s failed for student%s, because of query lesson", lesson_name, student_name)
		return err
	}

	// add process
	create_time := time.Now().Unix()
	insert_sql := `INSERT INTO student_lesson_process
	(student_id,student_name,lesson_id,lesson_name,fee_info,file_list,process_info,create_time)
	VALUES (?,?,?,?,?,?,?,?)`
	_, insert_err := dao.db.Exec(insert_sql, user.Id, user.Name, lesson_info.Id, lesson_info.Name, "", "", "", create_time)
	return insert_err
}
