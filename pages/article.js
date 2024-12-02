// 检查浏览器是否支持 MediaDevices API（用于摄像头访问）
if ('mediaDevices' in navigator && 'getUserMedia' in navigator.mediaDevices) {
  const captureButton = document.getElementById('captureButton');
  const uploadButton = document.getElementById('uploadButton');

  // 拍摄作文按钮点击事件
  captureButton.addEventListener('click', async () => {
    try {
      const formData = new FormData();
      formData.append('grade', document.getElementById('grade').value);
      formData.append('content_type', document.getElementById('content_type').value);
      formData.append('language', document.getElementById('content_language').value);
      formData.append('username', "default")

      const stream = await navigator.mediaDevices.getUserMedia({ video: true });
      const videoTrack = stream.getVideoTracks()[0];
      const imageCapture = new ImageCapture(videoTrack);
      const blob = await imageCapture.takePhoto();
      // 关闭视频流
      videoTrack.stop();
      // 这里将拍摄的图片发送到后端进行 OCR 处理
      // 假设后端的上传接口是 /upload

      formData.append('image', blob);
      console.log(formData);
      const response = await fetch('/upload', {
        method: 'POST',
        body: formData
      });
      const result = await response.text();
      console.log(result);
    } catch (error) {
      console.error('拍摄或上传错误:', error);
    }
  });

  // 作文上传按钮点击事件
  uploadButton.addEventListener('click', () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*,application/pdf';
    const formData = new FormData();
    formData.append('grade', document.getElementById('grade').value);
    formData.append('content_type', document.getElementById('content_type').value);
    formData.append('language', document.getElementById('content_language').value);
    formData.append('username', "default")
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
      console.log(formData);
      const response = await fetch('/upload_ocr', {
        method: 'POST',
        body: formData
      });
      const result = await response.text();
      console.log(result);
    };
    input.click();
  });
} else {
  console.log('浏览器不支持摄像头访问或文件上传。');
}