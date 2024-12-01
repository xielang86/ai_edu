let globalUsername
document.addEventListener('DOMContentLoaded', function () {
// 获取URL中的参数（假设课题名参数名为project_name，可根据实际修改）
const urlParams = new URLSearchParams(window.location.search);
if (typeof urlParams === "undefined") {
  // TODO(*): get role from a post
  username = globalUsername
  role = "student"
} else {
  const projectName = urlParams.get('project_name');
  username = urlParms.get("username")
  role = urlParams.get("role")
}

// 设置学生信息区域的内容
document.getElementById('student-username').textContent = username;
document.getElementById('project-name').textContent = projectName;

// 模拟设置课题进度（这里只是示例，实际需根据后端数据来设置颜色）
const progressCircles = document.querySelectorAll('.progress-circle');
progressCircles.forEach((circle, index) => {
  if (index < 3) {
    circle.classList.add('completed');
  }
});

// 问答中心按钮点击事件
document.getElementById('qa-button').addEventListener('click', () => {
  window.location.href = '/qa?username='; // 替换为实际问答中心页面地址
});

// 文件上传按钮点击事件
document.getElementById('upload-button').addEventListener('click', () => {
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = 'image/*,.pdf'; // 限制文件格式为图片和pdf
  input.multiple = false; // 一次只允许上传一个文件
  input.onchange = () => {
    const file = input.files[0];
    if (file) {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('username', username);
      formData.append('role', role);
      // 发送文件上传请求到后端的/upload接口
      fetch('/upload', {
        method: 'POST',
        body: formData
      })
      .then(response => response.text())
      .then(result => console.log(result))
      .catch(error => console.error('上传文件出错:', error));
    }
  };
  input.click();
});

// 我的文件集按钮点击事件（这里可补充实际功能代码）
document.getElementById('my-files-button').addEventListener('click', () => {
  // 比如跳转到我的文件集页面或者获取展示相关文件列表等操作
});

// 老师的文件集按钮点击事件（这里可补充实际功能代码）
document.getElementById('teacher-files-button').addEventListener('click', () => {
  // 比如跳转到老师的文件集页面或者获取展示相关文件列表等操作
});
});