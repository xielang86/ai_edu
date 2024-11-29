package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	// "github.com/tencentcloud/cos-go-sdk-v5"
)

var (
	local_file_root          = "../../" // just for testing
	tencent_cloud_cos_root   = ""
	tencent_cloud_cos_key    = ""
	tencent_cloud_cos_secret = ""
)

func FileCenterHandler(w http.ResponseWriter, r *http.Request) {
	PageHandler(w, "./pages/file_center.html")
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

func UploadHandler(w http.ResponseWriter, r http.Request) {
	// 检查请求方法是否为 POST
	if r.Method != "POST" {
		http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	// 读取上传的图片文件
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "读取图片文件失败", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 将图片文件内容读取到字节数组中
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "read file failed ", http.StatusInternalServerError)
		return
	}
	// save to local	path
	filename := fmt.Sprintf("%s_%d", header.Filename, time.Now().Unix())
	err = SaveFileToLocal(imageBytes, fmt.Sprintf("%s/%s", local_file_root, filename))
	if err != nil {
		http.Error(w, "store file failed", http.StatusInternalServerError)
	}

	responseData := ResponseData{
		Status:  "success",
		Message: fmt.Sprintf("upload %s succ", header.Filename),
	}

	PostResponse(w, responseData)
}
