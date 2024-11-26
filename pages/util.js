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

module.exports = {
    sendVerificationCode: sendVerificationCode,
};