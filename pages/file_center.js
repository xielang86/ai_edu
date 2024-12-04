// 获取页面URL中的role参数（简单示意，实际可能需要更完善的URL解析）
// const urlParams = new URLSearchParams(window.location.search);
/// const role = urlParams.get('role');

// 模拟从后端获取的文件和文件夹对应关系数据（实际需通过AJAX等请求后端接口）
const fileFolderData = {
  "student": {
    "all": [],
    "私有": [],
    "课题相关": ["课题1", "课题2"],
    "证书": ["证书1", "证书2"]
  },
  "teacher": {
    "all": [],
    "学生上传": ["学生1", "学生2"],
    "课题相关": ["课题A", "课题B"],
    "初始课题": ["初始课题1", "初始课题2"]
  }
};

// 渲染文件夹列表
function renderFolderList(role) {
  const folderSection = document.getElementById('folderSection');
  const folders = fileFolderData[role];
  for (const folder in folders) {
    const button = document.createElement('button');
    button.textContent = folder;
    button.onclick = function () {
      renderSubfolderList(folder);
    };
    folderSection.appendChild(button);
  }
}

// 渲染子文件夹列表
function renderSubfolderList(selectedFolder) {
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
      renderFileList(subfolder);
    };
    subfolderSection.appendChild(button);
  });
}

// 渲染文件列表
function renderFileList(selectedSubfolder, role) {
  const fileListSection = document.getElementById('fileListSection');
  fileListSection.innerHTML = ''; // 先清空之前的内容
  // 这里假设根据子文件夹能获取到对应的文件列表（实际需从后端获取对应关系）
  const files = []; // 模拟文件列表，实际需替换为真实数据
  files.forEach(file => {
    const button = document.createElement('button');
    button.textContent = file;
    button.onclick = function () {
      // 这里添加点击文件名后下载或打开文件的具体逻辑，目前暂为空
    };
    fileListSection.appendChild(button);
  });
}

// 跳转到个人中心页面的函数（目前只是简单跳转，实际需配置正确的页面URL等）
function goToPersonalCenter() {
  window.location.href = "personal_center.html";
}

function init(user_data) {
  username = user_data.username
  role = user_data.role
  document.getElementById("username").textContent = username
  renderFolderList(role)
  renderSubfolderList("all", role)
  renderFileList("all", role)
}

// 页面加载完成后渲染文件夹列表
window.onload = function () {
  const result = CheckAuth();
  // result.then(user_data=>renderFolderList(user_data)).then(user_data=>renderFileList(user_data));
  result.then(user_data_str=>{return JSON.parse(user_data_str)}).then(user_data=>init(user_data));
};