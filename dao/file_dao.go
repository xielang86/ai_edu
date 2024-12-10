package dao

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type FileInfo struct {
	Id                string
	Name              string `json:"name"`
	CloudPath         string `json:"cloud_path"`
	Md5               string
	RelatedLessonList string `json:"related_lessson_list"`
	Username          string `json:"username"`
	Role              string `json:"role"`
	AuthUserList      string
	CreateTime        int64
}

func CreateFileTable(dao *UserDAO) {
	var sql_str string = `
	CREATE TABLE IF NOT EXISTS  user_file (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
  	name VARCHAR(255),
  	cloud_path	VARCHAR(255),
		md5                 VARCHAR(255),
		related_lesson_list VARCHAR(255),
		username VARCHAR(128),
		role VARCHAR(32),
		auth_user_list VARCHAR(1024),
		create_time BIGINT);`
	CreateTable(dao, sql_str)
}

func QueryAllFileByUsername(dao *UserDAO, username string, result *[]FileInfo) error {
	query := fmt.Sprintf("select name, cloud_path, md5, related_lesson_list, username,role,auth_user_list from user_file where username=\"%s\"", username)
	rows, err := dao.db.Query(query)
	if err != nil {
		fmt.Println("查询错误:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var file_info FileInfo
		err := rows.Scan(&file_info.Name, &file_info.CloudPath, &file_info.Md5, &file_info.RelatedLessonList,
			&file_info.Username, &file_info.Role, &file_info.AuthUserList)
		if err != nil {
			fmt.Println("映射错误:", err)
			continue
		}

		// 将结构体添加到数组
		*result = append(*result, file_info)
	}

	return nil
}

func GetFileByMd5(dao *UserDAO, md5 string, info *FileInfo) error {
	query := fmt.Sprintf("select name, cloud_path, md5, related_lesson_list, username, role, auth_user_list from user_file where md5=\"%s\"", md5)
	row := dao.db.QueryRow(query)
	err := row.Scan(&info.Name, &info.CloudPath, &info.Md5, &info.RelatedLessonList,
		&info.Username, &info.Role, &info.AuthUserList)
	if err == sql.ErrNoRows {
		fmt.Println("没有找到匹配的行")
	} else if err != nil {
		fmt.Println("发生其他错误:", err)
	}
	return err
}

func AddFile(dao *UserDAO, username string, role string, name string, file_path string, md5 string, related_lesson string) error {
	// check by md5
	var info FileInfo
	err := GetFileByMd5(dao, md5, &info)
	if err != sql.ErrNoRows {
		fmt.Println("add failed for", name, err)
		return err
	}

	info.Name = name
	info.Username = username
	info.Role = role
	info.CloudPath = file_path
	info.Md5 = md5
	info.RelatedLessonList = related_lesson
	info.AuthUserList = username
	info.CreateTime = time.Now().Unix()
	// TODO(*): add auth to lesson's teacher

	insert_sql := `INSERT INTO user_file 
	(name,cloud_path, md5, related_lesson_list,username, role, auth_user_list,create_time)
	VALUES (?,?,?,?,?,?,?,?)`
	_, insert_err := dao.db.Exec(insert_sql, info.Name, info.CloudPath, info.Md5, info.RelatedLessonList,
		info.Username, info.Role, info.AuthUserList, info.CreateTime)
	return insert_err
}
