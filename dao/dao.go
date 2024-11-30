package dao

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// u can change dsn for different db

type StudentBaseInfo struct {
	Id           int64
	Name         string `json:"username"`
	Age          int8
	LessonId     string
	LessonName   string
	CreateTime   int64
	ParentName   string
	ParentDegree string `json:"degree"`
	ParentMajor  string `json:"major"`
	ParentCareer string `json:"jobDirection"`
	ParentSchool string `json:"graduateSchool"`
	Fee          int32
	PassPhrase   string `json:"password"`
	Phone        string `json:"phone"`
}

type LessonBaseInfo struct {
	id                  int64
	name                string
	teacher_id          int64
	teacher_name        string
	involved_student_id string
	init_file_list      string
	create_time         int64
}

type TeacherBaseInfo struct {
	id          int64
	name        string
	lesson_list string
	create_time int64
}

type StudentLessonProcess struct {
	student_id   int64
	student_name string
	lesson_id    int64
	lesson_name  string
	fee_info     string
	file_list    string
}

func QueryStudent(dao *UserDAO, name string, phone string, result *StudentBaseInfo) error {
	select_str := "select name,age,lesson_id,lesson_name,parent_name,parent_degree,parent_major,parent_career,fee,pass_phrase,phone from"
	name_query := ""
	phone_query := ""
	if len(name) > 1 {
		name_query = fmt.Sprintf("%s knowledge_edu.student_basic_info where name=\"%s\"", select_str, name)
	} else if len(phone) > 7 {
		phone_query = fmt.Sprintf("%s knowledge_edu.student_basic_info where phone=\"%s\"", select_str, phone)
	}

	var row *sql.Row
	if len(name_query) > 10 {
		row = dao.db.QueryRow(name_query)
	} else if len(phone_query) > 10 {
		row = dao.db.QueryRow(phone_query)
	} else {
		return fmt.Errorf("name or phone is too short name=%s and phone=%s", name, phone)
	}
	err := row.Scan(&result.Name, &result.Age, &result.LessonId, &result.LessonName, &result.ParentName, &result.ParentDegree,
		&result.ParentMajor, &result.ParentCareer, &result.Fee, &result.PassPhrase, &result.Phone)
	if err == sql.ErrNoRows && len(name_query) > 10 {
		row = dao.db.QueryRow(phone_query)
		err = row.Scan(&result.Name, &result.Age, &result.LessonId, &result.LessonName, &result.ParentName, &result.ParentDegree,
			&result.ParentMajor, &result.ParentCareer, &result.Fee, &result.PassPhrase, &result.Phone)
	}
	return err
}

func ModifyPassphrase(dao *UserDAO, name string, phone string, new_pass string) error {
	sql_str := fmt.Sprintf("UPDATE student_basic_info SET pass_phrase='%s' WHERE name=\"%s\" and phone=\"%s\"",
		new_pass, name, phone)
	_, err := dao.db.Exec(sql_str)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func QueryStudentLessonProcess(dao *UserDAO, query string, result *StudentLessonProcess) error {
	row := dao.db.QueryRow(query)

	err := row.Scan(&result.student_id, &result.student_name, &result.lesson_id, &result.lesson_name)
	if err != nil {
		return err
	}
	return nil
}

func CreateStudentTable(dao *UserDAO) {
	var sql_str string = `
	CREATE TABLE IF NOT EXISTS student_basic_info (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
		role VARCHAR(32) NOT NULL,
    age TINYINT,
    lesson_id VARCHAR(255),
    lesson_name VARCHAR(255),
    create_time BIGINT,
    parent_name VARCHAR(255),
    parent_degree VARCHAR(255),
    parent_major VARCHAR(255),
    parent_career VARCHAR(255),
    fee INT,
    pass_phrase VARCHAR(255),
    phone VARCHAR(255)
	);
	`
	CreateTable(dao, sql_str)
}
func CreateLessonTable(dao *UserDAO) {
	var sql_str string = `
	CREATE TABLE IF NOT EXISTS lesson_base_info (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    teacher_id BIGINT,
    teacher_name VARCHAR(255),
    involved_student_id VARCHAR(255),
    init_file_list VARCHAR(255),
    create_time BIGINT);
	`
	CreateTable(dao, sql_str)
}

func CreateTeacherTable(dao *UserDAO) {
	var sql_str string = `
	CREATE TABLE IF NOT EXISTS lesson_base_info (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    lesson_list VARCHAR(255),
    create_time BIGINT);`
	CreateTable(dao, sql_str)
}

func CreateStudentLessonProcessTable(dao *UserDAO) {
	var sql_str string = `
	CREATE TABLE IF NOT EXISTS student_lesson_process (
    student_id BIGINT,
    student_name VARCHAR(255),
    lesson_id BIGINT,
    lesson_name VARCHAR(255),
    fee_info	VARCHAR(255),
		file_list	VARCHAR(255),
		 PRIMARY KEY (student_id, lesson_id)
	);
	`
	CreateTable(dao, sql_str)
}

func InsertStudentBasicInfo(dao *UserDAO, info StudentBaseInfo) error {
	insert_sql := `INSERT INTO student_basic_info
	(name,age, lesson_id, lesson_name,create_time,parent_name,parent_degree,parent_major,parent_career,
	fee,pass_phrase,phone) 
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`
	_, insert_err := dao.db.Exec(insert_sql, info.Name, info.Age, info.LessonId, info.LessonName,
		info.CreateTime, info.ParentName, info.ParentDegree, info.ParentMajor, info.ParentCareer,
		info.Fee, info.PassPhrase, info.Phone)
	return insert_err
}
