function goToResetPassword() {
  window.location.href = "/reset_pass";
}

function validateForm() {
   var username = document.getElementById("username").value;
   var password = document.getElementById("password").value;
   var illegalChars = /[^\w]/;
   if (illegalChars.test(username)) {
     alert("用户名只能包含字母、数字和下划线");
     return false;
   }
   // 验证密码长度是否符合要求
   if (password.length < 8) {
     alert("密码长度至少为8位");
     return false;
   }
   return true;
}
document.addEventListener('DOMContentLoaded', function () {
  var form = document.getElementById('myForm');

  form.addEventListener('submit', async function (event) {
    // 阻止表单默认提交行为
    event.preventDefault();

    if (!validateForm()) {
      return;
    }
    // 使用fetch发送POST请求到后端的/register_post接口
    try {
      const response = await fetch('/login_post', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
      });
      const result = await response.json();
      if (response.status == 200 && result.status ==='success') {
        alert('success login');
        window.location.href = "./user_center";
      } else {
        alert(result.message || 'login failed，try again');
      }
    } catch (error) {
      console.error(error);
      alert('网络异常，请稍后再试');
    }
  });
});