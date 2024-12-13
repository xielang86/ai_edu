package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type VerifyCodeReqData struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

// VerificationCodeInfo 存储验证码相关信息
type VerificationCodeInfo struct {
	Code     string
	SendTime time.Time
}

var (
	verificationCodes = make(map[string]*VerificationCodeInfo)
	verify_mutex      sync.Mutex
	expirationMinutes = 2
)

// generateVerificationCode 生成6位数字验证码
func generateVerificationCode() string {
	min := 100000
	max := 999999
	return strconv.Itoa(min + rand.Intn(max-min))
}

// 调用阿里云短信服务发送验证码（此处需替换为真实有效的配置）
// templateID := "your_template_id"         // 替换为你在腾讯云短信服务配置的短信模板ID
// templateParam := []string{`{"code":"123456"}`}      // 替换为真实的模板参数，格式要与短信模板定义的参数格式一致
func AliSendSms(phoneNumber string, templateCode string, templateParam string) error {
	accessKeyId := "LTAI5tLUFcjzBxkSrKtrMGTz"
	accessKeySecret := "rxp45PovC6lq0rl09k1dGmyLFWnIva"
	// 模板CODE：SMS_254130904
	// accessKeyId := os.Getenv("ALIYUN_ACCESS_KEY_ID")
	// 你的阿里云AccessKey Secret
	// accessKeySecret := os.Getenv("ALIYUN_ACCESS_KEY_SECRET")

	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	fmt.Println("create client")
	if err != nil {
		fmt.Println(err)
		return err
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	// 设置接收短信的手机号码，多个号码以逗号分隔
	request.PhoneNumbers = phoneNumber
	// 设置短信签名名称，在阿里云短信服务控制台配置好的
	request.SignName = "清大开创"
	// 设置短信模板编码，在阿里云短信服务控制台配置好的
	request.TemplateCode = templateCode
	// 设置短信模板变量的JSON格式字符串，如果模板没有变量则传空字符串
	request.TemplateParam = templateParam
	response, err := client.SendSms(request)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if response.Code != "OK" {
		return fmt.Errorf("短信发送失败，错误码: %s，错误信息: %s", response.Code, response.Message)
	}

	return nil
}

func VerifyCode(phoneNumber, inputCode string, expirationMinutes int) bool {
	verify_mutex.Lock()
	codeInfo, ok := verificationCodes[phoneNumber]
	verify_mutex.Unlock()

	if ok {
		// 计算时间差
		elapsed := time.Since(codeInfo.SendTime).Minutes()
		if elapsed <= float64(expirationMinutes) && inputCode == codeInfo.Code {
			return true
		}
	}
	return false
}

func VerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 关闭请求体，释放资源
	defer r.Body.Close()

	// 解析JSON数据到User结构体
	var phone_code VerifyCodeReqData
	err = json.Unmarshal(body, &phone_code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("phone: %s, code: %s", phone_code.Phone, phone_code.Code)

	responseData := ResponseData{
		Status:  "success",
		Message: "verify code succ",
	}

	if VerifyCode(phone_code.Phone, phone_code.Code, expirationMinutes) {
	} else {
		fmt.Printf("verifycode erro or expored")
		responseData.Status = "failed"
		responseData.Message = "verify code failed"
	}
	PostResponse(w, responseData)
}

func SendSMSHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 关闭请求体，释放资源
	defer r.Body.Close()

	// 解析JSON数据到User结构体
	var info VerifyCodeReqData
	err = json.Unmarshal(body, &info)
	fmt.Printf("phone: %s", info.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if info.Phone == "" {
		http.Error(w, "手机号码不能为空", http.StatusBadRequest)
		return
	}

	// 生成6位验证码
	verificationCode := generateVerificationCode()

	// 根据你的短信模板变量名来设置JSON格式的参数
	templateParam := fmt.Sprintf("{\"code\":\"%s\"}", verificationCode)
	templateCode := "SMS_254130904"

	responseData := ResponseData{
		Status:  "success",
		Message: fmt.Sprintf("%d", expirationMinutes),
	}

	err = AliSendSms(info.Phone, templateCode, templateParam)
	if err != nil {
		responseData.Status = "failed"
		responseData.Message = "message failed to send"
	} else {
		verify_mutex.Lock()
		verificationCodes[info.Phone] = &VerificationCodeInfo{
			Code:     verificationCode,
			SendTime: time.Now(),
		}
		verify_mutex.Unlock()
	}
	PostResponse(w, responseData)
}
