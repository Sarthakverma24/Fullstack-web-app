document.getElementById('login-form').addEventListener('submit', async (e) => {
  e.preventDefault();

  const userLogin = {
    user_name: document.getElementById('user_name').value.trim(),
    password: document.getElementById('password').value.trim()
  };

  try {
    const response = await fetch('http://localhost:8081/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(userLogin)
    });

    if (!response.ok) {
      throw new Error('Failed to log-in');
    }

    const result = await response.json();
    alert('Login successful! Server says: ' + result.message);

    e.target.reset();
  } catch (error) {
    alert('Error: ' + error.message);
  }
});

document.getElementById('forget-btn').addEventListener('click', function () {
  document.getElementById('forget-password-modal').classList.remove('hidden');
});

// Hide the modal
document.getElementById('close-modal').addEventListener('click', function () {
  document.getElementById('forget-password-modal').classList.add('hidden');
});

// Handle recover form
document.getElementById('recover-form').addEventListener('submit',async function (e) {
  e.preventDefault();
  const UserDetails={
   gmail : document.getElementById('recover-email').value.trim(),
   phone : document.getElementById('recover-phone').value.trim()

  }
  try{
    const response =await fetch('http://localhost:8081/recover',{
      method:'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(UserDetails)
    });
    if(!response.ok){
      throw new Error('Failed to find user');
    }
    const result = await response.json();
    alert('Found password' + result.message);
    window.location.href = "../login/login.html";
    e.target.reset();
  }catch (error) {
    alert('Error: ' + error.message);
  }
  document.getElementById('forget-password-modal').classList.add('hidden');
});
