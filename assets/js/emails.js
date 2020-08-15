let emailRequestForm = document.querySelector('.email-request-form');

emailRequestForm.addEventListener('submit',(e)=>{
    e.preventDefault();
    fetch('http://localhost:3000/emails',{
        method: 'POST',
        headers: {
            'Content-Type':'application/json'
        },
        body: JSON.stringify({
            name: document.getElementById('name').value,
            email: document.getElementById('email').value,
            text: document.getElementById('message').value,
        })
    }).then((resp)=>resp.text()).then((data)=> alert('terimasih sudah memnerikan pesan, kami akan segera membalas dan memfollow up anda'), window.history.go());

})