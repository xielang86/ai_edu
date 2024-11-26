package dao

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// u can change dsn for different db

type UserDAO struct {
	db  *sql.DB
	dsn string
}

func NewUserDAO(db *sql.DB, dsn string) *UserDAO {
	return &UserDAO{
		db:  db,
		dsn: dsn,
	}
}

type EduKnowledge struct {
	level1        string
	level2        string
	level3        string
	level4        string
	level5_prompt string // maybe it's json
}

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

type FileInfo struct {
	id                  string
	name                string
	cloud_path          string
	md5                 string
	related_lesson_list string
	student_id_list     string
	teacher_id_list     string
}

func CloseDB(dao *UserDAO) {
	defer dao.db.Close()
	fmt.Printf("Close MySQL db: %s\n", dao.dsn)
}

func ConnectDB(dao *UserDAO) error {
	// open db connection
	var err error
	dao.db, err = sql.Open("mysql", dao.dsn)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// check the connection
	err = dao.db.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Connected to the database successfully!")

	// check the version, just for a query test
	var version string
	err = dao.db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("MySQL version: %s\n", version)
	return nil
}

func QueryEduKnowledge(dao *UserDAO, query string, result *EduKnowledge) error {
	row := dao.db.QueryRow(query)

	err := row.Scan(&result.level4, &result.level5_prompt)
	if err != nil {
		return err
	}
	return nil
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

func CreateTable(dao *UserDAO, sql_str string) {
	_, err := dao.db.Exec(sql_str)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("succ create")
}

// func CreateKnowledgeTable(dao *UserDAO) {
// 	var knowledge_edu_create_sql string = `
//   CREATE TABLE IF NOT EXISTS  knowledge_edu.en_knowledge_point (
//     id INT AUTO_INCREMENT PRIMARY KEY,
//     level1 VARCHAR(255),
//     level2 VARCHAR(255),
//     level3 VARCHAR(255),
//     level4 VARCHAR(255),
//     level5 VARCHAR(255),
//     level5_prompt VARCHAR(1024)
//   );
// 	`
// 	CreateTable(dao, knowledge_edu_create_sql)
// }

func CreateStudentTable(dao *UserDAO) {
	var sql_str string = `
	CREATE TABLE IF NOT EXISTS student_basic_info (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
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

func CreateFileTable(dao *UserDAO) {
	var sql_str string = `
	CREATE TABLE IF NOT EXISTS  student_lesson_process (
    id	BIGINT PRIMARY KEY AUTO_INCREMENT,
  	name VARCHAR(255),
  	cloud_path	VARCHAR(255),
		md5                 VARCHAR(255),
		related_lesson_list VARCHAR(255),
		student_id_list     VARCHAR(255),
		teacher_id_list     VARCHAR(255) );
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

// func main() {
// 	dao := &UserDAO{nil, kEduKnowledgeDB}
// 	ConnectDB(dao)
// 	var query string = "select level4, level5 from knowledge_edu.en_knowledge_point where level4=?"
// 	result := &EduKnowledge{}
// 	CreateKnowledgeTable(dao)
// 	QueryEduKnowledge(dao, query, result)
// 	CloseDB(dao)
// }
