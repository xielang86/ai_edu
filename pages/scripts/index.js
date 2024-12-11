// 获取用户协议链接元素并添加点击事件处理
const userAgreementLink = document.querySelector('.user-agreement');
userAgreementLink.addEventListener('click', function () {
    // 这里简单模拟弹出一个包含协议内容的对话框，实际中需要替换为真实的协议文本
    const agreementText = "这是标准的用户协议内容示例，你需要替换为真实详细的内容......";
    const dialog = window.alert(agreementText);
    // 如果需要更复杂的对话框，比如有确定按钮等，可以使用自定义的模态框组件或者HTML的dialog元素等进行扩展实现
});

// 获取隐私政策链接元素并添加点击事件处理
const userPrivacyLink = document.querySelector('.user-privacy');
userPrivacyLink.addEventListener('click', function () {
    // 同样模拟弹出隐私政策内容对话框
    const privacyText = "这是标准的隐私政策内容示例，你需要替换为真实详细的内容......";
    const dialog = window.alert(privacyText);
    // 后续可按实际需求优化对话框交互
});

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
  const button = document.getElementById('login-btn');
  button.addEventListener('click', async function (event) {
    // 阻止表单默认提交行为
    event.preventDefault();

    if (!validateForm()) {
      return;
    }
    const formData = {
        username: document.getElementById('username').value,
        password: document.getElementById('password').value,
    };
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
      if (result.status ==='success') {
        alert('success login');
        localStorage.setItem(JWT_KEY, result.token);
        setTimeout(() => {localStorage.removeItem(JWT_KEY);}, 24 * 60 * 60 * 1000);
        CheckAuth();
        username = result.data
        role = 
        window.location.href = "./user_project_center";
      } else {
        alert(result.message || 'login failed，try again');
      }
    } catch (error) {
      console.error(error);
      alert('网络异常，请稍后再试');
    }
  });

  const register_button = document.getElementById('register-btn');
  register_button.addEventListener('click', function (event) {
    window.location.href = "/register";
  })
});