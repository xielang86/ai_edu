const JWT_KEY = 'xxoo_jwt_token';
async function CheckAuth() {
  const token = localStorage.getItem(JWT_KEY);
  if (token) {
    const response = await fetch('./check_auth', {
      method: 'GET',
      headers: {
        'Authorization': 'Bearer' + token
      }
    });
    if (response.ok) {
      const data = await response.json()
      const parts = data.message.split(",")
      return JSON.stringify({username: parts[0], role:parts[1]})
    } else {
      alert('back to login');
      window.location.href = '/index';
      return null
    }
  } else {
    window.location.href = '/index';
    return null
  }
}

function getIconPath(fileName) {
    const fileExtension = fileName.split('.').pop().toLowerCase();
    switch (fileExtension) {
        case 'doc':
        case 'docx':
            return './pages/images/doc_icon.jpeg';
        case 'pdf':
            return './pages/images/pdf_icon.jpeg';
        case 'jpg':
        case 'png':
        case 'jpeg':
        case 'gif':
            return './pages/images/icon_image.png';
        default:
            return './pages/images/icon_default.png';
    }
}

function startCountdown(sendCodeBtn, countdown) {
  sendCodeBtn.disabled = true; // 禁用按钮，防止重复点击
  sendCodeBtn.textContent = `${countdown}秒后重发`;
  timer = setInterval(() => {
    countdown--;
    sendCodeBtn.textContent = `${countdown}秒后重发`;
    if (countdown === 0) {
      clearInterval(timer); // 倒计时结束，清除定时器
      timer = null; // 将定时器变量置空
      sendCodeBtn.disabled = false; // 恢复按钮可点击状态
      sendCodeBtn.textContent = '发送验证码';
      countdown = 60; // 重置倒计时初始值
    }
  }, 1000);
}

function BindSMS(smsButton, phoneInput) {
  smsButton.addEventListener("click", function() {
    const formData = {
    phone: phoneInput.value,
  }
   // check phone number first
  if (!phoneInput.value.match(/^1[3-9]\d{9}$/)) {
    alert('请输入正确的手机号');
    return;
  }

  const response = fetch('/send_verify_code', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(formData)
  }).then(response=>{
    if (response.status == 200) {
      return response.json()
    }
    return null
  }).then(result=>{
    if (result != null && result.status ==='success') {
      startCountdown(smsButton, parseInt(result.message) * 60)
    } else {
      alert(result.message);
    }
  });
});
}
function BindVerifyCode(phoneInput, codeInput, errorMessageDiv, succ_call_back) {
  const formData = {
    phone: phoneInput.value,
    code: codeInput.value,
  }
   // check phone number first
  if (!phoneInput.value.match(/^1[3-9]\d{9}$/)) {
    alert('请输入正确的手机号');
    return;
  }
  if (codeInput.value.length != 6) {
    alert('verify code must be 6 number');
    return;
  }

  const response = fetch('/verify_code', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(formData)
  }).then(response=>{
    if (response.status == 200) {
      return response.json()
    }
    return null
  }).then(result=>{
    if (result != null && result.status ==='success') {
      errorMessageDiv.style.display = 'none';
      succ_call_back()
    } else {
      errorMessageDiv.textContent = '验证码输入错误';
      errorMessageDiv.style.display = 'block';
    }
  });

}
