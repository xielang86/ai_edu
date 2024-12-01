package api

import (
	"api_server/dao"
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	// "github.com/tencentcloud/cos-go-sdk-v5"
)

var (
	local_file_root          = "./local_files" // just for testing
	tencent_cloud_cos_root   = ""
	tencent_cloud_cos_key    = ""
	tencent_cloud_cos_secret = ""
)

func FileCenterHandler(w http.ResponseWriter, r *http.Request) {
	PageHandler(w, "./pages/file_center.html")
}

func FilePageHandler(w http.ResponseWriter, r *http.Request) {
	PageHandler(w, "./pages/article.html")
}

func FileCenterPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}
	responseData := ResponseData{
		Status:  "success",
		Message: "操作成功",
	}

	PostResponse(w, responseData)
}

func SaveFileToLocal(fileData []byte, filePath string) error {
	err := os.WriteFile(filePath, fileData, 0644)
	if err != nil {
		log.Printf("保存文件到本地失败: %v", err)
		return err
	}
	return nil
}

// UploadFileToTencentCOS 将文件上传到腾讯云COS
// func UploadFileToTencentCOS(fileData []byte, bucket, objectKey string) error {
// 	// 填入你的腾讯云 SecretId 和 SecretKey
// 	secretId := "your_secret_id"
// 	secretKey := "your_secret_key"
// 	// 填入你的腾讯云COS所在的地域
// 	region := "your_region"
//
// 	u, _ := url.Parse(fmt.Sprintf("https://cos-%s.myqcloud.com", region))
// 	b := &cos.BaseURL{BucketURL: u}
// 	client := cos.NewClient(b, &http.Client{
// 		Transport: &cos.AuthorizationTransport{
// 			SecretID:  secretId,
// 			SecretKey: secretKey,
// 		},
// 	})
//
// 	fileData, err := ioutil.ReadFile(localFilePath)
// 	if err != nil {
// 		log.Printf("读取本地文件失败: %v", err)
// 		return err
// 	}
//
// 	_, err = client.Object.Put(context.Background(), objectKey, fileData, nil)
// 	if err != nil {
// 		log.Printf("上传文件到腾讯云COS失败: %v", err)
// 		return err
// 	}
//
// 	return nil
// }

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}
	file, header, err := r.FormFile("file")
	username := r.FormValue("username")
	lesson := r.FormValue("lesson")

	if err != nil {
		http.Error(w, "read file failed!", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 将图片文件内容读取到字节数组中
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "read file failed ", http.StatusInternalServerError)
		return
	}

	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	// cal md5
	hash := md5.Sum(fileBytes)
	md5_str := fmt.Sprintf("%x", hash)
	var file_info dao.FileInfo
	err = dao.GetFileByMd5(mydao, md5_str, &file_info)
	need_save := true

	responseData := ResponseData{
		Status:  "success",
		Message: fmt.Sprintf("upload %s succ", header.Filename),
	}

	if err != nil {
		need_save = false
		responseData.Message = fmt.Sprintf("dedup file for %s")
	}

	// save to local	path
	if need_save {
		filename := fmt.Sprintf("%s_%d", header.Filename, time.Now().Unix())
		local_file_path := fmt.Sprintf("%s/%s", local_file_root, filename)
		err = SaveFileToLocal(fileBytes, local_file_path)
		if err != nil {
			http.Error(w, "store file failed", http.StatusInternalServerError)
		}

		// query for user_id
		// TODO(*): just for test username, risk for attacking
		var info dao.UserInfo
		info.Name = username
		info.Id = 0
		if username != "default" {
			err = dao.QueryUser(mydao, username, "", &info)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "用户不存在或登录信息错误", http.StatusUnauthorized)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		} else {
			fmt.Println("default user for test")
		}
		// insert record to db
		err = dao.AddFile(mydao, info.Id, username, local_file_path, md5_str, lesson)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	PostResponse(w, responseData)
}

type FileData struct {
	Folder  interface{} `json:"folder"`
	AllFile interface{} `json:"file"`
}

func GetAllFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	mydao := dao.NewUserDAO(nil, kEduKnowledgeDB)
	dao.ConnectDB(mydao)
	defer dao.CloseDB(mydao)

	username := r.FormValue("username")
	var result []dao.FileInfo
	dao.QueryAllFileByUsername(mydao, username, &result)

	responseData := ResponseData{
		Status:  "success",
		Message: fmt.Sprintf("get files succ for user ", username),
	}

	var data FileData
	data.AllFile = result
	responseData.Data = data
	PostResponse(w, responseData)
}
