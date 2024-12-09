// 获取URL中的参数
function getUrlParams() {
  const params = {};
  const urlSearchParams = new URLSearchParams(window.location.search);
  for (const [key, value] of urlSearchParams.entries()) {
    params[key] = value;
  }
  return params;
}

async function createTeacherList() {
  let data = {
    username : username
  }

  const response = await fetch('/get_all_teacher', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });

  const result = await response.json();
  if (response.status == 200 && result.status ==='success') {
  } else {
    alert(result.message || 'get lesson failed for current user');
  }
  // result.data
  console.log(result.data)

  const teacherListSection = document.createElement('div');
  teacherListSection.classList.add('section');
  const title = document.createElement('h2');
  title.textContent = '老师列表';
  teacherListSection.appendChild(title);

  row = document.createElement('div');
  row.classList.add('row');
  teacherListSection.appendChild(row);

  if (result.data === null) {
    return teacherListSection;
  }

  await result.data.forEach((teacher, index) => {
    if (index % 4 === 0) {
      const newRow = document.createElement('div');
      newRow.classList.add('row');
      teacherListSection.appendChild(newRow);
      row = newRow;
    }
    const button = document.createElement('button');
    button.textContent = teacher.username;
    button.addEventListener('click', () => {
      // 假设跳转到老师个人介绍页，需替换真实链接
      window.location.href = `/personal_desc?username=${teacher.username}`;
    });
    row.appendChild(button);
  });

  return teacherListSection;
}

async function createProjectList() {
  let data = {
    username : username
  }
  const response = await fetch('/get_all_project', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });

  const result = response.json();
  if (response.status == 200 && result.status ==='success') {
  } else {
    alert(result.message || 'get lesson failed for current user');
  }
  // result.data
  console.log(result.data)
  
  const projectListSection = document.createElement('div');
  projectListSection.classList.add('section');
  const title = document.createElement('h2');
  title.textContent = '课题列表页';
  projectListSection.appendChild(title);

  const row = document.createElement('div');
  row.classList.add('row');
  projectListSection.appendChild(row);

  result.data.forEach(project => {
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

async function CreateStudentList() {
  // TODO(xl): 
  let data = {
    username : username
  }

  const response = await fetch('/get_all_student', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });

  const result = response.json();
  if (response.status == 200 && result.status ==='success') {
  } else {
    alert(result.message || 'get lesson failed for current user');
  }
  // result.data
  const studentListSection = document.createElement('div');
  studentListSection.classList.add('section');
  const title = document.createElement('h2');
  title.textContent = '学生列表';
  studentListSection.appendChild(title);

  const row = document.createElement('div');
  row.classList.add('row');
  studentListSection.appendChild(row);

  if (result.data === null) {
    return studentListSection;
  }

  result.data.forEach((student, index) => {
    if (index % 4 === 0) {
      const newRow = document.createElement('div');
      newRow.classList.add('row');
      studentListSection.appendChild(newRow);
      row = newRow;
    }
    const button = document.createElement('button');
    button.textContent = username;
    button.addEventListener('click', () => {
      // 这里假设跳转到学生课题列表页，实际需替换真实链接
      window.location.href = `student_lessons.html?id=${student.id}`;
    });
    row.appendChild(button);
  });

  return studentListSection;
}

async function createStudentProjectList(username) {
  // TODO(xl): safe problem, server must check jwt ?
  let data = {
    username : username
  }

  const response = await fetch('/get_all_project', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });

  const result = await response.json();
  if (response.status == 200 && result.status ==='success') {
  } else {
    alert(result.message || 'get lesson failed for current user');
  }
  const lessonListSection = document.createElement('div');
  lessonListSection.classList.add('section');
  const title = document.createElement('h2');
  title.textContent = '课题中心';
  title.addEventListener('click', () => {
    window.location.href = '/project_center';
  });
  lessonListSection.appendChild(title);

  const row = document.createElement('div');
  row.classList.add('row');
  lessonListSection.appendChild(row);

  if (result.data === null) {
    return lessonListSection;
  }

  result.data.forEach(lesson => {
    const button = document.createElement('button');
    button.textContent = lesson.name;
    button.addEventListener('click', () => {
      window.location.href = '/project?lesson_name='+lesson.name;
    });
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
    window.location.href = '/personal_desc?username=username';
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
    window.location.href = '/payment?username=username';
  });
  financeCenterSection.appendChild(button);
  return financeCenterSection;
}

async function init(user_data) {
  username = user_data.username
  role = user_data.role
  document.getElementById("username").textContent = username
  const app = document.getElementById('user_center');
  
  if (role ==='student') {
    // 模拟获取老师数据，实际需从后端获取
    const teacherListSection = await createTeacherList();
    const lessonListSection = await createStudentProjectList(username);

    const fileCenterSection = createFileCenter();
    const personalInfoSection = createPersonalInfo();
    const financeCenterSection = createFinanceCenter();

    app.appendChild(teacherListSection);
    app.appendChild(fileCenterSection);
    app.appendChild(lessonListSection);
    app.appendChild(personalInfoSection);
    app.appendChild(financeCenterSection);
  } else if (role === 'teacher') {
    // 模拟获取学生数据，实际需从后端获取
    const studentListSection = await createStudentList();
    const projectListSection = await createProjectList();

    const fileCenterSection = createFileCenter();

    app.appendChild(studentListSection);
    app.appendChild(projectListSection);
    app.appendChild(fileCenterSection);
  }
  return username + "succ loaded"
}

// 页面加载完成后渲染老师列表
window.onload = function () {
  CheckAuth().then(result=>JSON.parse(result)).then(user_data=>init(user_data));
};