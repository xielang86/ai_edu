let verificationCode; // 用于存储后端发送过来的验证码

// NOTE(all): form event listenner must be after the dom content loaded, if remove the document.addEvent, it would lead to null value for formvalue
// if u just put the js code block at the end of html page, it would be ok
async function DoRegister() {
 const inputVerificationCode = verificationCode;
 if (inputVerificationCode === verificationCode) {
   // 收集表单数据，根据角色判断是否包含家长相关字段
   const formData = {
     role: document.getElementById('role').value,
     username: document.getElementById('username').value,
     phone: document.getElementById('phone').value,
     password: document.getElementById('password').value,
   };
   if (formData.role === '家长') {
     formData.graduateSchool = document.getElementById('graduateSchool').value;
     formData.major = document.getElementById('major').value;
     formData.degree = document.getElementById('degree').value;
     formData.jobDirection = document.getElementById('jobDirection').value;
   }

   // 使用fetch发送POST请求到后端的/register_post接口
   try {
     const response = await fetch('/register_post', {
       method: 'POST',
       headers: {
         'Content-Type': 'application/json'
       },
       body: JSON.stringify(formData)
     });
     const result = await response.json();
     if (response.status == 200 && result.status ==='success') {
       alert('注册成功');
       window.location.href = "./index";
     } else {
       alert(result.message || '注册失败，请稍后再试');
     }
   } catch (error) {
     console.error(error);
     alert('网络异常，请稍后再试');
   }
 } else {
   const error = document.createElement('p');
   error.className = 'error';
   error.textContent = '验证码错误，请重新输入';
   form.appendChild(error);
 }
}

const backButton = document.getElementById('confirm-btn');
backButton.addEventListener('click', function() {
  DoRegister()
});