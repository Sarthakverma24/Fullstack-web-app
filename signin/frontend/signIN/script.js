document.getElementById('signin-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const userData = {
    user_id:document.getElementById('user_id').value.trim(),
    user_name: document.getElementById('user_name').value.trim(),
    phone_no: document.getElementById('phone_no').value.trim(),
    age: parseFloat(document.getElementById('age').value),
    gmail: document.getElementById('gmail').value.trim(),
    password:document.getElementById('password').value.trim()
  };

  try {
    const response = await fetch('http://localhost:8081/signin', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(userData)
    });

    if (!response.ok) {
      throw new Error('Failed to sign in');
    }

    const result = await response.json();
    alert( result.message);

    Login();

    // Clear form inputs after successful sign-in
    e.target.reset();
  } catch (error) {
    alert('Error: ' + error.message);
  }
});


function Login(){
 window.location.href = "../login/login.html";
}
