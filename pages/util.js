function sendVerificationCode(phone) {
  // const phone = document.getElementById('phone').value;
  if (!phone.match(/^1[3-9]\d{9}$/)) {
    alert('请输入正确的手机号');
    return;
  }
  // 使用fetch发送AJAX请求到后端获取验证码
  fetch('/send_verification_code', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ phone: phone })
  })
.then(response => response.json())
.then(data => {
      if (response.status == 200 && data.status ==='success') {
        verificationCode = data.code; // 保存后端返回的验证码
        alert('验证码已发送');
      } else {
        alert('验证码发送失败，请稍后再试');
      }
    })
.catch(error => {
      console.error(error);
      alert('发生错误，请稍后再试');
    });
}

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
