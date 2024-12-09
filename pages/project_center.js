// 模拟从后台获取课题数据以及是否参加的状态，实际中需要通过fetch等发送POST请求获取真实数据

const topicContainer = document.getElementById('topic-container');
function DoRenderAllProject(projects) {
  projects.forEach(project => {
    const topicDiv = document.createElement('div');
    topicDiv.classList.add('topic-item');

    const topicLink = document.createElement('a');
    topicLink.href = `/project?project_name=${project.name}`; // 假设课题主页根据课题id来区分，实际按真实链接设置
    topicLink.textContent = project.name;
    topicLink.classList.add('topic-link');

    const joinBtn = document.createElement('button');
    joinBtn.textContent = '参加课题';
    if (project.joined) {
        joinBtn.disabled = true;
        joinBtn.classList.add('disabled');
    } else {
        joinBtn.addEventListener('click', () => {
            alert('需要线下缴费，请拨打电话咨询');
        });
    }

    const cancelBtn = document.createElement('button');
    cancelBtn.textContent = '取消参加';
    cancelBtn.disabled =!project.joined;
    if (project.joined) {
        cancelBtn.addEventListener('click', () => {
            // 这里可以添加取消参加的逻辑，比如发送请求到后台更新状态等，暂不详细实现
            alert('已取消参加该课题');
        });
    }

    topicDiv.appendChild(topicLink);
    topicDiv.appendChild(joinBtn);
    topicDiv.appendChild(cancelBtn);

    topicContainer.appendChild(topicDiv);

  });
}

// 渲染课题列表
function init(user_data) {
  let data = {
    username : user_data.username,
    need_all: 1,
    role: user_data.role
  }

 // CheckAuth().then(result=>JSON.parse(result)).then(user_data=>init(user_data));

  fetch('/get_all_project', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  }).then(response=>{
    if (response.status == 200) {
      return response.json()
    }
    return null
  }).then(response_data=> {
    if (response_data == null || response_data.data == null) {
      alert(response_data.message || 'get lesson failed for current user');
    } else {
      DoRenderAllProject(response_data.data)
    }
  });

  const backButton = document.getElementById('backToUserCenter');
  backButton.addEventListener('click', function() {
    window.location.href = '/user_center';
  });


}

window.onload = function () {
  CheckAuth().then(result=>JSON.parse(result)).then(user_data=>init(user_data));
};