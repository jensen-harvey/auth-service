package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/herb-immortal/auth_service_hi/pkg/auth"
	"github.com/herb-immortal/auth_service_hi/pkg/database"
	"github.com/herb-immortal/auth_service_hi/pkg/utils"
)

// HTML for the simple auth UI
const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Auth Service</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
            font-family: Arial, sans-serif;
        }

        body {
            background-color: #f5f5f5;
            padding: 20px;
        }

        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: white;
            border-radius: 8px;
            padding: 30px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        }

        h1 {
            text-align: center;
            margin-bottom: 20px;
            color: #2c3e50;
        }

        h2 {
            margin-bottom: 20px;
            color: #3498db;
        }

        .tabs {
            display: flex;
            margin-bottom: 20px;
            border-bottom: 1px solid #ddd;
        }

        .tab {
            background: none;
            border: none;
            padding: 10px 20px;
            font-size: 16px;
            cursor: pointer;
            outline: none;
            color: #7f8c8d;
        }

        .tab.active {
            color: #3498db;
            border-bottom: 3px solid #3498db;
            font-weight: bold;
        }

        .tab.hidden {
            display: none;
        }

        .form-container {
            display: none;
        }

        .form-container.active {
            display: block;
        }

        .form-group {
            margin-bottom: 15px;
        }

        label {
            display: block;
            margin-bottom: 5px;
            color: #555;
        }

        input, select {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
        }

        .btn {
            background-color: #3498db;
            color: white;
            border: none;
            padding: 10px 15px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            width: 100%;
            margin-top: 10px;
        }

        .btn:hover {
            background-color: #2980b9;
        }

        .message {
            margin-top: 15px;
            padding: 10px;
            border-radius: 4px;
        }

        .message.success {
            background-color: #d4edda;
            color: #155724;
        }

        .message.error {
            background-color: #f8d7da;
            color: #721c24;
        }

        .profile-item {
            margin-bottom: 15px;
            padding: 10px;
            background-color: #f9f9f9;
            border-radius: 4px;
        }

        #logoutButton {
            background-color: #e74c3c;
        }

        #logoutButton:hover {
            background-color: #c0392b;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Herb Immortal Authentication</h1>
        
        <!-- Navigation Tabs -->
        <div class="tabs">
            <button id="loginTab" class="tab active">Login</button>
            <button id="signupTab" class="tab">Sign Up</button>
            <button id="profileTab" class="tab hidden">Profile</button>
        </div>

        <!-- Login Form -->
        <div id="loginForm" class="form-container active">
            <h2>Login</h2>
            <form id="login">
                <div class="form-group">
                    <label for="loginEmail">Email</label>
                    <input type="email" id="loginEmail" required>
                </div>
                <div class="form-group">
                    <label for="loginPassword">Password</label>
                    <input type="password" id="loginPassword" required>
                </div>
                <button type="submit" class="btn">Login</button>
                <p id="loginMessage" class="message"></p>
            </form>
        </div>

        <!-- Signup Form -->
        <div id="signupForm" class="form-container">
            <h2>Sign Up</h2>
            <form id="signup">
                <div class="form-group">
                    <label for="signupName">Full Name</label>
                    <input type="text" id="signupName" required>
                </div>
                <div class="form-group">
                    <label for="signupEmail">Email</label>
                    <input type="email" id="signupEmail" required>
                </div>
                <div class="form-group">
                    <label for="signupPhone">Phone Number</label>
                    <input type="text" id="signupPhone" required>
                </div>
                <div class="form-group">
                    <label for="signupPassword">Password</label>
                    <input type="password" id="signupPassword" required>
                </div>
                <div class="form-group">
                    <label for="signupRole">Role</label>
                    <select id="signupRole" required>
                        <option value="customer">Customer</option>
                        <option value="admin">Admin</option>
                        <option value="healer">Healer</option>
                        <option value="vendor">Vendor</option>
                    </select>
                </div>
                <button type="submit" class="btn">Sign Up</button>
                <p id="signupMessage" class="message"></p>
            </form>
        </div>

        <!-- Profile Page -->
        <div id="profilePage" class="form-container">
            <h2>User Profile</h2>
            <div id="profileData">
                <div class="profile-item">
                    <strong>Name:</strong> <span id="profileName"></span>
                </div>
                <div class="profile-item">
                    <strong>Email:</strong> <span id="profileEmail"></span>
                </div>
                <div class="profile-item">
                    <strong>Phone:</strong> <span id="profilePhone"></span>
                </div>
                <div class="profile-item">
                    <strong>Role:</strong> <span id="profileRole"></span>
                </div>
                <div class="profile-item">
                    <strong>User ID:</strong> <span id="profileId"></span>
                </div>
                <div class="profile-item">
                    <strong>Account Created:</strong> <span id="profileCreated"></span>
                </div>
            </div>
            <button id="logoutButton" class="btn">Logout</button>
        </div>
    </div>

    <script>
        document.addEventListener('DOMContentLoaded', function() {
            // API endpoints
            const API_URL = window.location.origin;
            const SIGNUP_ENDPOINT = API_URL + '/api/auth/signup';
            const LOGIN_ENDPOINT = API_URL + '/api/auth/login';
            const PROFILE_ENDPOINT = API_URL + '/api/auth/profile';

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
                            'Authorization': 'Bearer ' + token,
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
    </script>
</body>
</html>
`

func main() {
	// Database configuration
	// In a production environment, these would be environment variables
	dbConfig := &database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "auth_service",
		SSLMode:  "disable",
	}

	// Connect to database
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create database tables if they don't exist
	if err := database.CreateTables(db); err != nil {
		log.Fatalf("Failed to create database tables: %v", err)
	}

	// Initialize repositories
	userRepo := database.NewUserRepository(db)

	// Initialize JWT token manager
	// In a production environment, these would be environment variables
	jwtSecret := "your-secret-key"        // Use a strong secret key in production
	jwtIssuer := "auth-service"           // Your application name
	jwtTTL := 24 * time.Hour              // Token validity duration
	tokenManager := utils.NewTokenManager(jwtSecret, jwtIssuer, jwtTTL)

	// Initialize authentication service
	authService := auth.NewAuthService(userRepo, tokenManager)

	// Initialize HTTP handler
	httpHandler := auth.NewHTTPHandler(authService)

	// Set up HTTP router
	mux := http.NewServeMux()
	httpHandler.SetupRoutes(mux)
	
	// Add a handler for the root path to serve the UI
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If the path is exactly "/", serve the HTML UI
		if r.URL.Path == "/" {
			fmt.Fprint(w, htmlTemplate)
			return
		}
		
		// Otherwise, return 404 Not Found
		http.NotFound(w, r)
	})
	
	// Print welcome message with usage information
	log.Println("=================================================")
	log.Println("Auth Service is running!")
	log.Println("API Endpoints:")
	log.Println("  POST http://localhost:8080/api/auth/signup - Create a new user")
	log.Println("  POST http://localhost:8080/api/auth/login - Login")
	log.Println("  GET http://localhost:8080/api/auth/profile - Get user profile (protected)")
	log.Println("Frontend:")
	log.Println("  http://localhost:8080/ - Web interface")
	log.Println("=================================================")

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}