let callMeForm = document.querySelector('.call-me-form');
callMeForm.addEventListener('submit',(e)=>{
    let phoneInput =document.querySelector('input');
    e.preventDefault();
    fetch('http://localhost:3000/callback-request',{
        method: 'POST',
        headers: {
            'Content-Type':'application/json'
        },
        body: JSON.stringify({
            phoneNumber: phoneInput.value
        })
    }).then((resp)=>resp.text()).then(()=> alert('we will call back you as soon as possible'), window.history.go());

})