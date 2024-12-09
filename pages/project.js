function init(user_data) {
  username = user_data.username
// 获取URL中的参数（假设课题名参数名为project_name，可根据实际修改）
const urlParams = new URLSearchParams(window.location.search);
const projectName = urlParams.get('lesson_name');

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
  uploadButton = document.getElementById("upload-button")
  uploadButton.addEventListener('click', () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*,application/pdf';
    const formData = new FormData();
    formData.append('lesson_name', projectName);
    formData.append('username', username)
    input.multiple = true
    input.onchange= async () => {
      for (var i = 0; i < input.files.length; i++) {
        var file = input.files[i]
        if (file) {
          if (file.size > 5 * 1024 * 1024) { // 5MB
            alert('文件大小超过限制，请选择小于 5MB 的文件,skip');
            continue
          }
          formData.append('files', file);
        }
      }
      const response = await fetch('/upload', {
        method: 'POST',
        body: formData
      });
      const result = await response.text();
      // console.log(result);
      alert(result)
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

const backButton = document.getElementById('back-usercenter-button');
  backButton.addEventListener('click', function() {
    window.location.href = '/user_center';
  });

}

window.onload = function () {
  CheckAuth().then(result=>JSON.parse(result)).then(user_data=>init(user_data));
};