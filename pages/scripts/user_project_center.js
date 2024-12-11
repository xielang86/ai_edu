function RenderFileList(files, itemsPerPage, currentPage) {
  const fileItemList = document.getElementById('file-list');
  fileItemList.innerHTML = ''; // 先清空之前的内容

  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  for (let i = startIndex; i < endIndex && i < files.length; i++) {
    const file = files[i]
    const fileName = file.name;
    const iconPath = getIconPath(fileName);
    const fileItem = document.createElement('div');
    fileItem.classList.add('file-item');
    const img = document.createElement('img');
    img.src = iconPath;
    img.alt = '文件图标';
    img.height = "32";
    const a = document.createElement('a');
    a.href = file.cloud_path;
    a.textContent = fileName;
    fileItem.appendChild(img);
    fileItem.appendChild(a);
    fileItemList.appendChild(fileItem);
  }
}

function RenderPagination(files, currentPage, itemsPerPage, totalPages) {
  const pagination = document.querySelector('.pagination');
  pagination.innerHTML = '';
  for (let i = 1; i <= totalPages; i++) {
    const pageItem = document.createElement('span');
    pageItem.classList.add('page-item');
    pageItem.textContent = i;
    if (i === currentPage) {
        pageItem.classList.add('current');
    }
    pageItem.addEventListener('click', function () {
        currentPage = i;
        RenderFileList(files, itemsPerPage, currentPage);
        RenderPagination(files, totalPages, currentPage);
    });
    pagination.appendChild(pageItem);
  }
  if (totalPages > 5) {
    const ellipsis = document.createElement('span');
    ellipsis.textContent = '...';
    ellipsis.classList.add('ellipsis');
    pagination.appendChild(ellipsis);
  }
  const nextBtn = document.createElement('button');
  nextBtn.textContent = '下一页';
  nextBtn.addEventListener('click', function () {
    if (currentPage < totalPages) {
        currentPage++;
        RenderFileList(files, itemsPerPage, currentPage);
        RenderPagination(files, totalPages, currentPage);
    }
  });
  const firstBtn = document.createElement('button');
  firstBtn.textContent = '首页';
  firstBtn.addEventListener('click', function () {
      currentPage = 1;
      RenderFileItems(files, itemsPerPage, currentPage);
      RenderPagination(files, totalPages, currentPage);
  });
  const lastBtn = document.createElement('button');
  lastBtn.textContent = '尾页';
  lastBtn.addEventListener('click', function () {
      currentPage = totalPages;
      RenderFileItems(files, itemsPerPage, currentPage);
      RenderPagination(files, totalPages, currentPage);
  });
  pagination.appendChild(nextBtn);
  pagination.appendChild(firstBtn);
  pagination.appendChild(lastBtn);
}

async function InitRenderFileList(username, to_last_page) {
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

    const itemsPerPage = 4;
    // 总页数
    let totalPages = Math.ceil(files.length / itemsPerPage);
    // 当前页码，初始化为1
    let currentPage = 1;
    if (to_last_page) {
      currentPage = totalPages
    }

    RenderFileList(files, itemsPerPage, currentPage)
    RenderPagination(files, currentPage, itemsPerPage, totalPages)
  }
}

// 跳转到个人中心页面的函数（目前只是简单跳转，实际需配置正确的页面URL等）
function DoUpload(username) {
  return function() {
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
    // update  files and direct to the last page
    InitRenderFileList(username, true)
  }
}

function init(user_data) {
  username = user_data.username
  role = user_data.role
  document.getElementById("username").textContent = username

  InitRenderFileList(username, false)

  uploadButton = document.getElementById("upload-button")
  uploadButton.addEventListener('click', DoUpload(username))

}

// 页面加载完成后渲染文件夹列表
window.onload = function () {
  const result = CheckAuth();
  result.then(user_data_str=>{return JSON.parse(user_data_str)}).then(user_data=>init(user_data));
};