// 获取确认按钮元素
let verificationCode; // 用于存储后端发送过来的验证码（这里只是模拟，实际需后端交互）

document.addEventListener('DOMContentLoaded', function () {
  const confirmBtn = document.querySelector('.confirm-btn');

  // 为确认按钮添加点击事件监听器
  confirmBtn.addEventListener('click', function () {
    // 收集表单数据
  const inputVerificationCode = document.getElementById('verify_code').value;
  verificationCode = inputVerificationCode
  if (inputVerificationCode === verificationCode) {
    const formData = {
        username: document.getElementById('username').value,
        phone: document.getElementById('phone').value,
        password: document.getElementById('password').value,
        passwordConfirm: document.getElementById('password-confirm').value
    };
    if (formData.password === formData.passwordConfirm) {
      try {
        const response = fetch('/reset_pass_post', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(formData)
        });
        response.then(response=>{return response.json()}).then(data=>{
          if (data.status ==='success') {
            alert('change pass success');
            window.location.href = "./index";
          } else {
            alert(data.message || 'change failed ，retry again');
          }
        })
      } catch (error) {
        console.error(error);
        alert('网络异常，请稍后再试');
      }
    } else {
      const error = document.createElement('p');
      error.className = 'error';
      error.textContent = '两次输入的密码不一致，请重新输入';
      form.appendChild(error);
    }

  } else {
    const error = document.createElement('p');
    error.className = 'error';
    error.textContent = '验证码错误，请重新输入';
    form.appendChild(error);
  }
});
});