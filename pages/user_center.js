// 模拟从后端获取老师名单数据（实际需通过AJAX等方式请求后端接口）
const teachers = ["张老师", "李老师", "王老师", "赵老师", "刘老师", "陈老师", "杨老师", "黄老师", "周老师", "吴老师"];

// 渲染老师列表按钮
function renderTeacherList() {
  const teacherListDiv = document.getElementById('teacherList');
  teachers.forEach((teacher, index) => {
    const button = document.createElement('button');
    button.className = 'teacher-button';
    button.textContent = teacher;
    button.onclick = function () {
      // 这里可以添加点击老师按钮后的具体逻辑，比如查看老师详情等，目前暂为空
    };
    teacherListDiv.appendChild(button);
  });
}

// 跳转到文件柜页面的函数（目前只是简单跳转，实际需配置正确的页面URL等）
function goToFileCabinet() {
  window.location.href = "/file_center";
}

// 跳转到课题中心页面的函数
function goToProjectCenter() {
  window.location.href = "/project_center";
}

// 跳转到个人信息页面的函数
function goToPersonalInfo() {
  window.location.href = "/personal_info";
}

// 跳转到财务中心页面的函数
function goToFinanceCenter() {
  window.location.href = "/finance_center";
}

// 页面加载完成后渲染老师列表
window.onload = function () {
  CheckAuth();
  // TODO(*): use ajax fetch teach list from web server
  renderTeacherList();
};