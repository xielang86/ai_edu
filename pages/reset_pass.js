let verificationCode; // 用于存储后端发送过来的验证码（这里只是模拟，实际需后端交互）

document.addEventListener('DOMContentLoaded', function () {
const form = document.getElementById('resetPasswordForm');
form.addEventListener('submit', async function (e) {
  e.preventDefault();
  const inputVerificationCode = document.getElementById('verificationCode').value;
  const newPassword = document.getElementById('newPassword').value;
  const confirmPassword = document.getElementById('confirmPassword').value;

  if (inputVerificationCode === verificationCode) {
    if (newPassword === confirmPassword) {
      try {
        const response = await fetch('/reset_pass_post', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(formData)
        });
        const result = await response.json();
        if (response.status == 200 && result.status ==='success') {
          alert('change pass success');
          window.location.href = "./login";
        } else {
          alert(result.message || 'change failed ，retry again');
        }
      } catch (error) {
        console.error(error);
        alert('网络异常，请稍后再试');
      }

      alert('密码重置成功');
      window.location.href = "./login";
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