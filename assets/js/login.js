let loginForm = document.querySelector('.form-login');
let id_user = document.querySelector('#id_user');
let password = document.querySelector('#password');
loginForm.addEventListener('submit', (e) => {
    let data = new FormData();
    data.append('id_user', id_user.value);
    data.append('password', password.value);

    e.preventDefault();
    fetch('http://localhost:9000/login_process', {
        method: 'POST',
        body: data
    }).then((res) => res.text()).then((data) => alert('data has been save'), window.history.go());
    if (data != null) {
        alert("selamat anda berhasil login")
    } else {
        alert("login gagal")
    }
})