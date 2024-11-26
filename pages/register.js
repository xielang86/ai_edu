let verificationCode; // 用于存储后端发送过来的验证码

function sendVerificationCode() {
  const phone = document.getElementById('phone').value;
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
function toggleFields() {
  const role = document.getElementById('role').value;
  const parentFields = document.getElementById('parentFields');
  if (role === '家长') {
    parentFields.style.display = 'block';
  } else {
    parentFields.style.display = 'none';
  }
}

function enableOtherInput() {
  const graduateSchool = document.getElementById('graduateSchool').value;
  const otherInput = document.getElementById('graduateSchoolOther');
  if (graduateSchool === '其他') {
    otherInput.style.display = 'block';
  } else {
    otherInput.style.display = 'none';
  }
}
// NOTE(all): form event listenner must be after the dom content loaded, if remove the document.addEvent, it would lead to null value for formvalue
// if u just put the js code block at the end of html page, it would be ok
document.addEventListener('DOMContentLoaded', function () {
  const form = document.getElementById('registrationForm');
  form.addEventListener('submit', async function (e) {
    e.preventDefault();
    // const inputVerificationCode = document.getElementById('verificationCode').value;
    // TODO(xl): waiting for really code
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
          window.location.href = "./login";
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
  });
});