package dao

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type LessonInfo struct {
	Id                  int64
	Name                string `json:"name"`
	TeacherId           int64
	TeacherName         string `json:"teacher_name"`
	InvolvedStudentId   string
	InvolvedStudentName string `json:"involved_student_name"`
	InitFileList        string `json:"init_file_list"`
	CreateTime          int64
	Desc                string `json:"desc"`
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
		text VARCHAR(1024),
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
		PRIMARY KEY (student_name, lesson_name)
	);
	`
	CreateTable(dao, sql_str)
}

func QueryLessonByName(dao *UserDAO, lesson_name string, info *LessonInfo) error {
	query := fmt.Sprintf("select id,name,teacher_id,teacher_name,involved_student_id,involved_student_name,init_file_list,text from lesson_info where name=\"%s\"", lesson_name)
	row := dao.db.QueryRow(query)
	err := row.Scan(&info.Id, &info.Name, &info.TeacherId, &info.TeacherName, &info.InvolvedStudentId, &info.InvolvedStudentName, &info.InitFileList, &info.Desc)
	if err != nil {
		fmt.Printf("query lesson failed: %s, query=%s\n", err, query)
		return err
	}
	return nil
}

func QueryLessonProcessByName(dao *UserDAO, student_name string, lesson_name string, info *StudentLessonProcess) error {
	query := fmt.Sprintf("select student_name,lesson_name,file_list,process_info from student_process_info where student_name=\"%s\" and lesson_name=\"%s\"",
		student_name, lesson_name)
	row := dao.db.QueryRow(query)
	err := row.Scan(&info.StudentName, &info.LessonName, &info.FileList, &info.ProcessInfo)
	if err != nil {
		return err
	}
	return nil
}

func QueryAllActiveLesson(dao *UserDAO, result *[]LessonInfo) error {
	query := "select id,name, teacher_name,involved_student_name,init_file_list,text from lesson_info where create_time>0"
	rows, err := dao.db.Query(query)
	if err != nil {
		fmt.Println("查询错误:", err)
		return err
	}
	defer rows.Close()
	var info LessonInfo
	for rows.Next() {
		err := rows.Scan(&info.Id, &info.Name, &info.TeacherName, &info.InvolvedStudentName, &info.InitFileList, &info.Desc)
		if err != nil {
			fmt.Println("query allactive less mapping error:", err)
			continue
		}
		// 将结构体添加到数组
		*result = append(*result, info)
	}
	return err
}

func QueryAllLessonNameByUsername(dao *UserDAO, username string, result *[]string) error {
	query := fmt.Sprintf("select lesson_name from student_lesson_process where student_name=\"%s\"", username)

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

func AddFileForLesson(dao *UserDAO, lesson_name string, filename_list []string) error {
	var info LessonInfo
	err := QueryLessonByName(dao, lesson_name, &info)
	if err != nil {
		return err
	}
	info.InitFileList = AddAndDedupString(info.InitFileList, filename_list)
	update_sql := `UPDATE lesson_info SET init_file_list=? WHERE id=?;`
	_, update_err := dao.db.Exec(update_sql, info.InitFileList, info.Id)
	return update_err
}

func AddFileForStudentLesson(dao *UserDAO, student_name string, lesson_name string, filename_list []string) error {
	var info StudentLessonProcess
	err := QueryLessonProcessByName(dao, student_name, lesson_name, &info)
	if err != nil {
		return err
	}

	info.FileList = AddAndDedupString(info.FileList, filename_list)
	update_sql := `UPDATE student_lesson_process SET file_list=? WHERE student_name=? and lesson_name=?;`
	_, update_err := dao.db.Exec(update_sql, info.FileList, info.StudentName, info.LessonName)
	return update_err
}

func AddLesson(dao *UserDAO, name string, teacher_name string, init_file_list string, text string) error {
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
	(name, teacher_id, teacher_name,involved_student_id,involved_student_name,init_file_list,text,create_time)
	VALUES (?,?,?,?,?,?,?,?)`
	_, insert_err := dao.db.Exec(insert_sql, name, user.Id, teacher_name, "", "", init_file_list, text, create_time)
	if insert_err != nil {
		fmt.Printf("insert into lesson_info failed for lesson %s", name)
		return insert_err
	}

	// for teacher, the lesson he teached
	user.LessonName = AddAndDedupString(user.LessonName, []string{name})

	update_sql := `UPDATE user_info SET lesson_name=? WHERE id=?;`
	_, update_err := dao.db.Exec(update_sql, user.LessonName, user.Id)
	return update_err
}

func AddAndDedupString(origin_str string, new_str_array []string) string {
	strMap := make(map[string]bool)
	keys := strings.Split(origin_str, ",")
	for _, key := range keys {
		if len(key) == 0 {
			continue
		}
		strMap[key] = true // 将分割后的每个元素作为键存入map，值这里简单设为true，可按需调整
	}
	buffer := bytes.NewBufferString("")
	for key := range strMap {
		buffer.WriteString(key)
		buffer.WriteString(",")
	}

	for _, new_str := range new_str_array {
		_, ok := strMap[new_str]
		if !ok {
			buffer.WriteString(new_str)
		}
	}
	return buffer.String()
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

	if insert_err != nil {
		fmt.Printf("insert lesson process failed for %s, lesson=%s", user.Name, lesson_name)
		return insert_err
	}
	// update userinfo's lesson_name and lesson_id
	user.LessonName = AddAndDedupString(user.LessonName, []string{lesson_name})

	update_sql := `UPDATE user_info SET lesson_name=? WHERE id=?;`
	_, update_err := dao.db.Exec(update_sql, user.LessonName, user.Id)
	if update_err != nil {
		fmt.Printf("update lesson failed for user_info user=%s, lesson=%s", user.Name, lesson_name)
		return update_err
	}
	return nil
}

func QueryAllStudentNameByTeacher(dao *UserDAO, teacher_name string, result *[]string) error {
	query := fmt.Sprintf("select InvolvedStudentName from lesson_info where teacher_name=\"%s\"", teacher_name)

	rows, err := dao.db.Query(query)
	if err != nil {
		fmt.Println("查询错误:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var str_list string
		err := rows.Scan(&str_list)
		if err != nil {
			fmt.Println("映射错误:", err)
			continue
		}
		parts := strings.Split(str_list, ",")
		// 将结构体添加到数组
		*result = append(*result, parts...)
	}
	return nil
}
