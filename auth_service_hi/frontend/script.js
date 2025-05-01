document.addEventListener('DOMContentLoaded', function() {
    // API endpoints
    const API_URL = 'http://localhost:8080';
    const SIGNUP_ENDPOINT = `${API_URL}/api/auth/signup`;
    const LOGIN_ENDPOINT = `${API_URL}/api/auth/login`;
    const PROFILE_ENDPOINT = `${API_URL}/api/auth/profile`;

    // DOM elements
    const loginTab = document.getElementById('loginTab');
    const signupTab = document.getElementById('signupTab');
    const profileTab = document.getElementById('profileTab');
    const loginForm = document.getElementById('loginForm');
    const signupForm = document.getElementById('signupForm');
    const profilePage = document.getElementById('profilePage');
    const logoutButton = document.getElementById('logoutButton');

    // Messages
    const loginMessage = document.getElementById('loginMessage');
    const signupMessage = document.getElementById('signupMessage');

    // Check if user is already logged in
    checkAuthStatus();

    // Tab switching
    loginTab.addEventListener('click', () => showTab('login'));
    signupTab.addEventListener('click', () => showTab('signup'));
    profileTab.addEventListener('click', () => showTab('profile'));

    // Form submissions
    document.getElementById('login').addEventListener('submit', handleLogin);
    document.getElementById('signup').addEventListener('submit', handleSignup);
    logoutButton.addEventListener('click', handleLogout);

    // Tab switching function
    function showTab(tabName) {
        // Hide all tabs
        loginTab.classList.remove('active');
        signupTab.classList.remove('active');
        profileTab.classList.remove('active');
        loginForm.classList.remove('active');
        signupForm.classList.remove('active');
        profilePage.classList.remove('active');

        // Show selected tab
        if (tabName === 'login') {
            loginTab.classList.add('active');
            loginForm.classList.add('active');
        } else if (tabName === 'signup') {
            signupTab.classList.add('active');
            signupForm.classList.add('active');
        } else if (tabName === 'profile') {
            profileTab.classList.add('active');
            profilePage.classList.add('active');
        }
    }

    // Check authentication status
    function checkAuthStatus() {
        const token = localStorage.getItem('token');
        if (token) {
            // Show profile tab and hide login/signup tabs
            profileTab.classList.remove('hidden');
            fetchProfile();
            showTab('profile');
        } else {
            // Show login/signup tabs and hide profile tab
            profileTab.classList.add('hidden');
            showTab('login');
        }
    }

    // Handle login form submission
    async function handleLogin(event) {
        event.preventDefault();
        
        const email = document.getElementById('loginEmail').value;
        const password = document.getElementById('loginPassword').value;
        
        try {
            const response = await fetch(LOGIN_ENDPOINT, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password })
            });
            
            const data = await response.json();
            
            if (response.ok) {
                // Store token in localStorage
                localStorage.setItem('token', data.token);
                // Update UI
                loginMessage.textContent = 'Login successful!';
                loginMessage.className = 'message success';
                // Redirect to profile page
                profileTab.classList.remove('hidden');
                fetchProfile();
                showTab('profile');
            } else {
                loginMessage.textContent = data.error || 'Login failed. Please try again.';
                loginMessage.className = 'message error';
            }
        } catch (error) {
            loginMessage.textContent = 'An error occurred. Please try again later.';
            loginMessage.className = 'message error';
            console.error('Login error:', error);
        }
    }

    // Handle signup form submission
    async function handleSignup(event) {
        event.preventDefault();
        
        const name = document.getElementById('signupName').value;
        const email = document.getElementById('signupEmail').value;
        const phone_number = document.getElementById('signupPhone').value;
        const password = document.getElementById('signupPassword').value;
        const role = document.getElementById('signupRole').value;
        
        try {
            const response = await fetch(SIGNUP_ENDPOINT, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    name,
                    email,
                    phone_number,
                    password,
                    role
                })
            });
            
            const data = await response.json();
            
            if (response.ok) {
                signupMessage.textContent = 'Signup successful! Please login with your new account.';
                signupMessage.className = 'message success';
                // Clear form
                document.getElementById('signup').reset();
                // Switch to login tab after a brief delay
                setTimeout(() => showTab('login'), 2000);
            } else {
                signupMessage.textContent = data.error || 'Signup failed. Please try again.';
                signupMessage.className = 'message error';
            }
        } catch (error) {
            signupMessage.textContent = 'An error occurred. Please try again later.';
            signupMessage.className = 'message error';
            console.error('Signup error:', error);
        }
    }

    // Fetch profile data
    async function fetchProfile() {
        const token = localStorage.getItem('token');
        
        if (!token) {
            showTab('login');
            return;
        }
        
        try {
            const response = await fetch(PROFILE_ENDPOINT, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                }
            });
            
            if (response.ok) {
                const user = await response.json();
                // Populate profile data
                document.getElementById('profileName').textContent = user.name;
                document.getElementById('profileEmail').textContent = user.email;
                document.getElementById('profilePhone').textContent = user.phone_number;
                document.getElementById('profileRole').textContent = user.role;
                document.getElementById('profileId').textContent = user.id;
                document.getElementById('profileCreated').textContent = new Date(user.created_at).toLocaleString();
            } else {
                // Token might be invalid or expired
                handleLogout();
            }
        } catch (error) {
            console.error('Profile fetch error:', error);
            // On error, clear token and redirect to login
            handleLogout();
        }
    }

    // Handle logout
    function handleLogout() {
        localStorage.removeItem('token');
        profileTab.classList.add('hidden');
        showTab('login');
        loginMessage.textContent = 'You have been logged out.';
        loginMessage.className = 'message';
    }
});