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
	expirationMinutes = 1
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
func AliSendSMS(phoneNumbers string, templateID string, templateParam string) error {
	// accessKeyId("LTAI5tLUFcjzBxkSrKtrMGTz")
	// accessKeySecret("rxp45PovC6lq0rl09k1dGmyLFWnIva")
	// 模板CODE：SMS_254130904
	// err := SendSMS({phone}, "", {code})
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "验证码发送失败"})
	// 	return
	// }
	fmt.Printf("send code %s for %s", templateParam, phoneNumbers)
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
		fmt.Fprintf(w, "verifycode pass")
	} else {
		fmt.Fprintf(w, "verifycode erro or expored")
		responseData.Status = "failed"
		responseData.Message = "verify code failed"
	}
	PostResponse(w, responseData)
}

func SendSMSHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "sms handler parse form failed!", http.StatusBadRequest)
		return
	}
	phoneNumber := r.Form.Get("phone_number")
	if phoneNumber == "" {
		http.Error(w, "手机号码不能为空", http.StatusBadRequest)
		return
	}

	// 生成6位验证码
	verificationCode := generateVerificationCode()

	// 根据你的短信模板变量名来设置JSON格式的参数
	templateParam := fmt.Sprintf("{\"jinxiao project verify code\":\"%s\"}", verificationCode)
	templateCode := "SMS_254130904"

	responseData := ResponseData{
		Status:  "success",
		Message: fmt.Sprintf("%d", expirationMinutes),
	}

	err = AliSendSMS(phoneNumber, templateCode, templateParam)
	if err != nil {
		responseData.Status = "failed"
		responseData.Message = "message failed to send"
	} else {
		verify_mutex.Lock()
		verificationCodes[phoneNumber] = &VerificationCodeInfo{
			Code:     verificationCode,
			SendTime: time.Now(),
		}
		verify_mutex.Unlock()
		fmt.Fprintf(w, "验证码已发送，请注意查收")
	}
	PostResponse(w, responseData)
}
