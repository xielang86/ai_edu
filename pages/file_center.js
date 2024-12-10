// 获取页面URL中的role参数（简单示意，实际可能需要更完善的URL解析）
// const urlParams = new URLSearchParams(window.location.search);
/// const role = urlParams.get('role');

// 模拟从后端获取的文件和文件夹对应关系数据（实际需通过AJAX等请求后端接口）
const fileFolderData = {
  "student": {
    "all": [],
    // "私有": [],
    // "课题相关": ["课题1", "课题2"],
    // "证书": ["证书1", "证书2"]
  },
  "teacher": {
    "all": [],
    "学生上传": [],
    // "课题相关": ["课题A", "课题B"],
    // "初始课题": ["初始课题1", "初始课题2"]
  }
};

// 渲染文件夹列表
function renderFolderList(role, username) {
  const folderSection = document.getElementById('folderSection');
  const folders = fileFolderData[role];
  for (const folder in folders) {
    const button = document.createElement('button');
    button.textContent = folder;
    button.onclick = function () {
      renderSubfolderList(folder, role, username);
    };
    folderSection.appendChild(button);
  }
}

// 渲染子文件夹列表
function renderSubfolderList(selectedFolder, role, username) {
  const subfolderSection = document.getElementById('subfolderSection');
  subfolderSection.innerHTML = ''; // 先清空之前的内容
  const subfolders = fileFolderData[role][selectedFolder];
  if (subfolders.length === 0) {
    return;
  }
  subfolders.forEach(subfolder => {
    const button = document.createElement('button');
    button.textContent = subfolder;
    button.onclick = function () {
      renderFileList(subfolder, role, username);
    };
    subfolderSection.appendChild(button);
  });
}

// 渲染文件列表
async function renderFileList(selectedSubfolder, role, username) {
  const fileListSection = document.getElementById('fileListSection');
  fileListSection.innerHTML = ''; // 先清空之前的内容
  // 这里假设根据子文件夹能获取到对应的文件列表（实际需从后端获取对应关系）
  let data = {
    username : username
  }

  const response = await fetch('/get_all_file', {
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

  if (result.data != null && result.data.file != null) {
    files = await result.data.file
    files.forEach(file => {
      const fileLink = document.createElement('a');
      filename = file.name
      fileLink.href = file.cloud_path;
      fileLink.textContent = file.name;

      // 根据文件扩展名设置不同的下载属性或其他行为
      const fileExtension = filename.split('.').pop().toLowerCase();
      if (['doc', 'docx'].includes(fileExtension)) {
        fileLink.target = '_blank';
      } else if (fileExtension === 'pdf') {
        fileLink.download = filename;
        fileLink.target = '_blank';
      } else if (['jpg', 'png', 'jpeg', 'gif'].includes(fileExtension)) {
        fileLink.target = '_blank';
      }

      const listItem = document.createElement('li');
      listItem.appendChild(fileLink);
      fileListSection.appendChild(listItem);
    });
  }
}

// 跳转到个人中心页面的函数（目前只是简单跳转，实际需配置正确的页面URL等）
function goToPersonalCenter() {
  window.location.href = "personal_center.html";
}

function init(user_data) {
  username = user_data.username
  role = user_data.role
  document.getElementById("username").textContent = username
  document.getElementById("page_title").textContent = "用户课题中心"
  uploadButton = document.getElementById("upload-button")
  uploadButton.addEventListener('click', () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*,application/pdf';
    const formData = new FormData();
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

  renderFileList("all", role, username)
  });

  renderFolderList(role)
  renderSubfolderList("all", role, username)
  renderFileList("all", role, username)
}

// 页面加载完成后渲染文件夹列表
window.onload = function () {
  const result = CheckAuth();
  // result.then(user_data=>renderFolderList(user_data)).then(user_data=>renderFileList(user_data));
  result.then(user_data_str=>{return JSON.parse(user_data_str)}).then(user_data=>init(user_data));
};