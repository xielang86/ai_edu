// 获取URL中的参数
function getUrlParams() {
  const params = {};
  const urlSearchParams = new URLSearchParams(window.location.search);
  for (const [key, value] of urlSearchParams.entries()) {
    params[key] = value;
  }
  return params;
}

function createTeacherList(teachers) {
  const teacherListSection = document.createElement('div');
  teacherListSection.classList.add('section');
  const title = document.createElement('h2');
  title.textContent = '老师列表';
  teacherListSection.appendChild(title);

  row = document.createElement('div');
  row.classList.add('row');
  teacherListSection.appendChild(row);

  teachers.forEach((teacher, index) => {
    if (index % 4 === 0) {
      const newRow = document.createElement('div');
      newRow.classList.add('row');
      teacherListSection.appendChild(newRow);
      row = newRow;
    }
    const button = document.createElement('button');
    button.textContent = teacher.name;
    button.addEventListener('click', () => {
      // 这里假设跳转到学生课题列表页，实际需替换真实链接
      window.location.href = `student_lessons.html?id=${student.id}`;
    });
    row.appendChild(button);
  });

  return teacherListSection;
}

function createProjectList(projects) {
  const projectListSection = document.createElement('div');
  projectListSection.classList.add('section');
  const title = document.createElement('h2');
  title.textContent = '课题列表页';
  projectListSection.appendChild(title);

  const row = document.createElement('div');
  row.classList.add('row');
  projectListSection.appendChild(row);

  projects.forEach(project => {
    const button = document.createElement('button');
    button.textContent = project.name;
    row.appendChild(button);
  });

  const addProjectButton = document.createElement('button');
  addProjectButton.textContent = '增加课题';
  addProjectButton.addEventListener('click', () => {
    // 这里简单模拟弹出对话框，实际需要用更完善的模态框组件等
    const dialog = document.createElement('dialog');
    const dialogContent = document.createElement('div');
    const descriptionInput = document.createElement('input');
    descriptionInput.placeholder = '描述课题';
    const uploadButton = document.createElement('button');
    uploadButton.textContent = '上传文件';
    const fileList = document.createElement('div');
    const confirmButton = document.createElement('button');
    confirmButton.textContent = '确认';
    const cancelButton = document.createElement('button');
    cancelButton.textContent = '取消';

    dialogContent.appendChild(descriptionInput);
    dialogContent.appendChild(uploadButton);
    dialogContent.appendChild(fileList);
    dialogContent.appendChild(confirmButton);
    dialogContent.appendChild(cancelButton);
    dialog.appendChild(dialogContent);
    document.body.appendChild(dialog);

    cancelButton.addEventListener('click', () => {
      dialog.close();
    });
    confirmButton.addEventListener('click', () => {
      // 这里模拟添加课题逻辑，实际需与后端交互并更新列表展示
      const newProject = {
        name: descriptionInput.value
      };
      projects.push(newProject);
      updateProjectList(projectListSection, projects);
      dialog.close();
    });
  });

  projectListSection.appendChild(addProjectButton);

  return projectListSection;
}

function updateProjectList(projectListSection, projects) {
  const row = projectListSection.querySelector('.row');
  row.innerHTML = '';
  projects.forEach(project => {
    const button = document.createElement('button');
    button.textContent = project.name;
    row.appendChild(button);
  });
}

function CreateStudentList(students) {
  const studentListSection = document.createElement('div');
  studentListSection.classList.add('section');
  const title = document.createElement('h2');
  title.textContent = '老师列表';
  studentListSection.appendChild(title);

  const row = document.createElement('div');
  row.classList.add('row');
  studentListSection.appendChild(row);

  students.forEach((student, index) => {
    if (index % 4 === 0) {
      const newRow = document.createElement('div');
      newRow.classList.add('row');
      studentListSection.appendChild(newRow);
      row = newRow;
    }
    const button = document.createElement('button');
    button.textContent = student.name;
    button.addEventListener('click', () => {
      // 假设跳转到老师个人介绍页，需替换真实链接
      window.location.href = `teacher_info.html?id=${teacher.id}`;
    });
    row.appendChild(button);
  });

  return studentListSection;
}

function createLessonList(lessons) {
  const lessonListSection = document.createElement('div');
  lessonListSection.classList.add('section');
  const title = document.createElement('h2');
  title.textContent = '课题中心';
  lessonListSection.appendChild(title);

  const row = document.createElement('div');
  row.classList.add('row');
  lessonListSection.appendChild(row);

  lessons.forEach(lesson => {
    const button = document.createElement('button');
    button.textContent = lesson.name;
    row.appendChild(button);
  });

  return lessonListSection;
}

function createFileCenter() {
  const fileCenterSection = document.createElement('div');
  fileCenterSection.classList.add('section');
  const button = document.createElement('button');
  button.textContent = '文件中心';
  button.addEventListener('click', () => {
    window.location.href = '/file_center';
  });
  fileCenterSection.appendChild(button);
  return fileCenterSection;
}

function createPersonalInfo() {
  const personalInfoSection = document.createElement('div');
  personalInfoSection.classList.add('section');
  const button = document.createElement('button');
  button.textContent = '个人信息';
  button.addEventListener('click', () => {
    window.location.href = 'personal_page.html';
  });
  personalInfoSection.appendChild(button);
  return personalInfoSection;
}

function createFinanceCenter() {
  const financeCenterSection = document.createElement('div');
  financeCenterSection.classList.add('section');
  const button = document.createElement('button');
  button.textContent = '财务中心';
  button.addEventListener('click', () => {
    window.location.href = 'payment_page.html';
  });
  financeCenterSection.appendChild(button);
  return financeCenterSection;
}

function init(user_data) {
  username = user_data.username
  role = user_data.role
  document.getElementById("username").textContent = username
  const app = document.getElementById('user_center');
  
  if (role ==='student') {
    // 模拟获取老师数据，实际需从后端获取
    const teachers = [
      { name: '张老师' }, { name: '李老师' }, { name: '王老师' }, { name: '赵老师' },
      { name: '孙老师' }, { name: '刘老师' }, { name: '陈老师' }, { name: '杨老师' }
    ];
    const teacherListSection = createTeacherList(teachers);
    const fileCenterSection = createFileCenter();
    // 模拟获取课题数据，实际需从后端获取
    const lessons = [
      { name: '课题1' }, { name: '课题2' }, { name: '课题3' }
    ];
    const lessonListSection = createLessonList(lessons);
    const personalInfoSection = createPersonalInfo();
    const financeCenterSection = createFinanceCenter();

    app.appendChild(teacherListSection);
    app.appendChild(fileCenterSection);
    app.appendChild(lessonListSection);
    app.appendChild(personalInfoSection);
    app.appendChild(financeCenterSection);
  } else if (role === 'teacher') {
    // 模拟获取学生数据，实际需从后端获取
    const students = [
      { name: '学生1', id: 1 }, { name: '学生2', id: 2 }, { name: '学生3', id: 3 }, { name: '学生4', id: 4 },
      { name: '学生5', id: 5 }, { name: '学生6', id: 6 }, { name: '学生7', id: 7 }, { name: '学生8', id: 8 }
    ];
    const studentListSection = createStudentList(students);
    const projectListSection = createProjectList([]);
    const fileCenterSection = createFileCenter();

    app.appendChild(studentListSection);
    app.appendChild(projectListSection);
    app.appendChild(fileCenterSection);
  }
}

// 页面加载完成后渲染老师列表
window.onload = function () {
  const result = CheckAuth();
  result.then(user_data_str=>{return JSON.parse(user_data_str)}).then(user_data=>init(user_data));
};