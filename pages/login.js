// document.addEventListener('DOMContentLoaded', function () {
//   var form = document.getElementById('myForm');
// 
//   form.addEventListener('submit', function (event) {
//     // 阻止表单默认提交行为
//     event.preventDefault();
// 
//     var username = document.getElementById('username').value;
//     var password = document.getElementById('password').value;
//     var illegalChars = /[^\w]/;
// 
//     if (illegalChars.test(username)) {
//       alert('用户名只能包含字母、数字和下划线');
//       return;
//     }
// 
//     // 验证密码长度是否符合要求
//     if (password.length < 8) {
//       alert('密码长度至少为8位');
//       return;
//     }
// 
//     // 验证通过，可以继续提交表单
//     this.submit();
//   });
// });
function validateForm() {
   var username = document.getElementById("username").value;
   var password = document.getElementById("password").value;
   var illegalChars = /[^\w]/;
   if (illegalChars.test(username)) {
     alert("用户名只能包含字母、数字和下划线");
     return false;
   }
   // 验证密码长度是否符合要求
   if (password.length < 8) {
     alert("密码长度至少为8位");
     return false;
   }
   return true;
 }