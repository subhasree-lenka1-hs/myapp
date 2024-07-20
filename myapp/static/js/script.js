function validateForm() {
    const dob = document.getElementById('dob').value;
    const email = document.getElementById('email').value;
    const mobile = document.getElementById('mobile').value;
    const password = document.getElementById('password').value;

    const today = new Date();
    const birthDate = new Date(dob);
    const age = today.getFullYear() - birthDate.getFullYear();
    const month = today.getMonth() - birthDate.getMonth(); 

    if (month < 0 || (month === 0 && today.getDate() < birthDate.getDate())) {
        age--;
    }

    if (age < 16) {
        alert('You must be at least 16 years old.');
        return false;
    }

    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailPattern.test(email)) {
        alert('Please enter a valid email address.');
        return false;
    }

    const mobilePattern = /^\d{10}$/;
    if (!mobilePattern.test(mobile)) {
        alert('Please enter a valid mobile number (10 digits).');
        return false;
    }

   
    const passwordPattern = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/;
    if (!passwordPattern.test(password)) {
        alert('Password must be at least 8 characters long and include at least one letter and one number.');
        return false;
    }

    return true;
}

function togglePasswordVisibility() {
    var passwordField = document.getElementById("password");
    if (passwordField.type === "password") {
        passwordField.type = "text";
    } else {
        passwordField.type = "password";
    }
}

window.onload = function() {
    document.getElementById("registrationForm").reset();
};