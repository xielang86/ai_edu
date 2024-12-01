package dao

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// u can change dsn for different db

type UserInfo struct {
	Id           int64
	Name         string `json:"username"`
	Role         string `json:"role"`
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
	Desc         string `json:"desc"`
}

func QueryUser(dao *UserDAO, name string, phone string, result *UserInfo) error {
	select_str := "select name,role,age,lesson_id,lesson_name,parent_name,parent_degree,parent_major,parent_career,fee,pass_phrase,phone,text from"
	name_query := ""
	phone_query := ""
	if len(name) > 1 {
		name_query = fmt.Sprintf("%s knowledge_edu.user_info where name=\"%s\"", select_str, name)
	} else if len(phone) > 7 {
		phone_query = fmt.Sprintf("%s knowledge_edu.user_info where phone=\"%s\"", select_str, phone)
	}

	var row *sql.Row
	if len(name_query) > 10 {
		row = dao.db.QueryRow(name_query)
	} else if len(phone_query) > 10 {
		row = dao.db.QueryRow(phone_query)
	} else {
		return fmt.Errorf("name or phone is too short name=%s and phone=%s", name, phone)
	}
	err := row.Scan(&result.Name, &result.Role, &result.Age, &result.LessonId, &result.LessonName, &result.ParentName, &result.ParentDegree,
		&result.ParentMajor, &result.ParentCareer, &result.Fee, &result.PassPhrase, &result.Phone, &result.Desc)
	if err == sql.ErrNoRows && len(name_query) > 10 {
		row = dao.db.QueryRow(phone_query)
		err = row.Scan(&result.Name, &result.Role, &result.Age, &result.LessonId, &result.LessonName, &result.ParentName, &result.ParentDegree,
			&result.ParentMajor, &result.ParentCareer, &result.Fee, &result.PassPhrase, &result.Phone, &result.Desc)
	}
	return err
}

func ModifyPassphrase(dao *UserDAO, name string, phone string, new_pass string) error {
	sql_str := fmt.Sprintf("UPDATE user_info SET pass_phrase='%s' WHERE name=\"%s\" and phone=\"%s\"",
		new_pass, name, phone)
	_, err := dao.db.Exec(sql_str)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func QueryStudentLessonProcess(dao *UserDAO, query string, result *StudentLessonProcess) error {
	row := dao.db.QueryRow(query)

	err := row.Scan(&result.ProcessInfo, &result.StudentName, &result.LessonId, &result.LessonName)
	if err != nil {
		return err
	}
	return nil
}

func CreateUserInfoTable(dao *UserDAO) {
	var sql_str string = `
	CREATE TABLE IF NOT EXISTS user_info (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
		role VARCHAR(32) NOT NULL,
    age TINYINT,
    lesson_id VARCHAR(128),
    lesson_name VARCHAR(128),
    create_time BIGINT,
    parent_name VARCHAR(128),
    parent_degree VARCHAR(128),
    parent_major VARCHAR(128),
    parent_career VARCHAR(128),
		parent_school VARCHAR(128),
    fee INT,
    pass_phrase VARCHAR(64),
    phone VARCHAR(64),
		text VARCHAR(1024)
	);
	`
	CreateTable(dao, sql_str)
}
func InsertUserInfo(dao *UserDAO, info UserInfo) error {
	insert_sql := `INSERT INTO user_info
	(name,role,age, lesson_id, lesson_name,create_time,parent_name,parent_degree,parent_major,parent_career,
	parent_school,fee,pass_phrase,phone,text) 
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, insert_err := dao.db.Exec(insert_sql, info.Name, info.Role, info.Age, info.LessonId, info.LessonName,
		info.CreateTime, info.ParentName, info.ParentDegree, info.ParentMajor, info.ParentCareer, info.ParentSchool,
		info.Fee, info.PassPhrase, info.Phone, info.Desc)
	return insert_err
}
