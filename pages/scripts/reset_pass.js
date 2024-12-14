function DoReset() {
  // 为确认按钮添加点击事件监听器
  confirmBtn.addEventListener('click', function () {
    // 收集表单数据
    const formData = {
        username: document.getElementById('username').value,
        phone: document.getElementById('phone').value,
        code: document.getElementById("verify_code").value,
        password: document.getElementById('password').value,
        passwordConfirm: document.getElementById('password-confirm').value
    };
    if (formData.password === formData.passwordConfirm) {
      errorMessageDiv.style.display = 'none';
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
      const errorMessageDiv = document.getElementById('errorPassMessage');
      errorMessageDiv.textContent = '两次密码不一致';
      errorMessageDiv.style.display = 'block';
    }
});

}

document.addEventListener('DOMContentLoaded', function () {
  const confirmBtn = document.querySelector('.confirm-btn');
  const phoneInput = document.getElementById("phone")
  const codeInput = document.getElementById('verify_code')
  const errorMessageDiv = document.getElementById('errorMessage');

  confirmButton.addEventListener('click', function() {
    BindVerifyCode(phoneInput, codeInput, errorMessageDiv, DoReset)
  });

  const smsButton = document.getElementById("sms_btn");
  BindSMS(smsButton, phoneInput)

});