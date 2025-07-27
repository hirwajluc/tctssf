// Authentication handling
document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
});

async function handleLogin(e) {
    e.preventDefault();
    
    const submitBtn = e.target.querySelector('button[type="submit"]');
    const originalText = submitBtn.innerHTML;
    
    // Clear any previous messages
    hideError();
    hideSuccess();
    
    // Show loading state
    showLoading(submitBtn);
    
    const formData = new FormData(e.target);
    const loginData = {
        email: formData.get('email'),
        password: formData.get('password')
    };
    
    try {
        const response = await apiCall('/login', {
            method: 'POST',
            body: JSON.stringify(loginData)
        });
        
        // Store session data
        localStorage.setItem('authToken', response.token);
        localStorage.setItem('currentUser', JSON.stringify(response.user));
        
        showSuccess('Login successful! Redirecting...');
        
        // Redirect based on role
        setTimeout(() => {
            if (response.user.role === 'member') {
                window.location.href = 'member-dashboard.html';
            } else if (response.user.role === 'admin' || response.user.role === 'superadmin') {
                window.location.href = 'admin-dashboard.html';
            } else if (response.user.role === 'treasurer') {
                window.location.href = 'treasurer-dashboard.html';
            } else {
                window.location.href = '/';
            }
        }, 1000);
        
    } catch (error) {
        // Show error message in red
        showError(error.message || 'Login failed. Please check your credentials and try again.');
        hideLoading(submitBtn, originalText);
    }
}

// Check if user should be redirected from login page
function checkLoginRedirect() {
    const token = localStorage.getItem('authToken');
    const user = localStorage.getItem('currentUser');
    
    if (token && user) {
        try {
            const userData = JSON.parse(user);
            // Verify token is still valid
            apiCall('/profile').then(() => {
                // Token is valid, redirect to dashboard
                if (userData.role === 'member') {
                    window.location.href = 'member-dashboard.html';
                } else if (userData.role === 'admin' || userData.role === 'superadmin') {
                    window.location.href = 'admin-dashboard.html';
                } else if (userData.role === 'treasurer') {
                    window.location.href = 'treasurer-dashboard.html';
                }
            }).catch(() => {
                // Token is invalid, clear storage
                localStorage.removeItem('authToken');
                localStorage.removeItem('currentUser');
            });
        } catch (error) {
            // Invalid stored data, clear it
            localStorage.removeItem('authToken');
            localStorage.removeItem('currentUser');
        }
    }
}

// Logout function
function logout() {
    localStorage.removeItem('authToken');
    localStorage.removeItem('currentUser');
    window.location.href = '/';
}

// Check if user is authenticated
function isAuthenticated() {
    const token = localStorage.getItem('authToken');
    const user = localStorage.getItem('currentUser');
    return token && user;
}

// Get current user data
function getCurrentUser() {
    const user = localStorage.getItem('currentUser');
    if (user) {
        try {
            return JSON.parse(user);
        } catch (error) {
            return null;
        }
    }
    return null;
}

// Get authentication token
function getAuthToken() {
    return localStorage.getItem('authToken');
}

// Check if user has specific role
function hasRole(role) {
    const user = getCurrentUser();
    return user && user.role === role;
}

// Check if user has any of the specified roles
function hasAnyRole(roles) {
    const user = getCurrentUser();
    return user && roles.includes(user.role);
}

// ===== NEW FUNCTIONS - ADD THESE AFTER THE EXISTING FUNCTIONS =====

/**
 * Require user to be authenticated
 * Redirects to login page if not authenticated
 * @returns {boolean} true if authenticated, false if redirected
 */
function requireAuth() {
    if (!isAuthenticated()) {
        console.log('No authentication found, redirecting to login');
        window.location.href = '/';
        return false;
    }
    return true;
}

/**
 * Require user to have specific role(s)
 * Redirects to appropriate dashboard if wrong role
 * @param {string|string[]} allowedRoles - Single role or array of roles
 * @returns {boolean} true if user has required role, false if redirected
 */
function requireRole(allowedRoles) {
    const user = getCurrentUser();
    
    if (!user) {
        console.log('No user found for role check');
        window.location.href = '/';
        return false;
    }
    
    // Convert single role to array for consistency
    const roles = Array.isArray(allowedRoles) ? allowedRoles : [allowedRoles];
    
    if (!roles.includes(user.role)) {
        console.log(`User role '${user.role}' not in allowed roles:`, roles);
        // Redirect to appropriate dashboard based on user's actual role
        redirectToUserDashboard(user.role);
        return false;
    }
    
    return true;
}

/**
 * Redirect user to their appropriate dashboard based on role
 * @param {string} role - User's role
 */
function redirectToUserDashboard(role) {
    switch (role) {
        case 'member':
            if (!window.location.pathname.includes('member-')) {
                window.location.href = 'member-dashboard.html';
            }
            break;
        case 'admin':
        case 'superadmin':
            if (!window.location.pathname.includes('admin-')) {
                window.location.href = 'admin-dashboard.html';
            }
            break;
        case 'treasurer':
            if (!window.location.pathname.includes('treasurer-')) {
                window.location.href = 'treasurer-dashboard.html';
            }
            break;
        default:
            window.location.href = '/';
    }
}

/**
 * Update user display elements on the page
 * Call this after successful authentication to show user info
 */
function updateUserDisplay() {
    const user = getCurrentUser();
    if (!user) return;
    
    // Update user name displays
    const userNameElements = document.querySelectorAll('.user-name, [data-user-name]');
    userNameElements.forEach(element => {
        element.textContent = `${user.first_name} ${user.last_name}`;
    });
    
    // Update user email displays
    const userEmailElements = document.querySelectorAll('.user-email, [data-user-email]');
    userEmailElements.forEach(element => {
        element.textContent = user.email;
    });
    
    // Update user role displays
    const userRoleElements = document.querySelectorAll('.user-role, [data-user-role]');
    userRoleElements.forEach(element => {
        element.textContent = user.role.charAt(0).toUpperCase() + user.role.slice(1);
    });
    
    // Update account number displays
    const accountNumberElements = document.querySelectorAll('.account-number, [data-account-number]');
    accountNumberElements.forEach(element => {
        element.textContent = user.account_number;
    });
}

/**
 * Initialize authentication for the current page
 * Call this in DOMContentLoaded if you need to set up auth
 */
function initializePage() {
    // Update user display elements
    if (isAuthenticated()) {
        updateUserDisplay();
    }
    
    // Set up logout buttons
    const logoutButtons = document.querySelectorAll('.logout-btn, [data-logout]');
    logoutButtons.forEach(button => {
        button.addEventListener('click', (e) => {
            e.preventDefault();
            logout();
        });
    });
}

/**
 * Create empty table row for when no data is available
 * @param {number} colspan - Number of columns to span
 * @param {string} message - Message to display
 * @returns {string} HTML string for empty row
 */
function createEmptyTableRow(colspan, message = 'No data available') {
    return `
        <tr>
            <td colspan="${colspan}" style="text-align: center; padding: 40px; color: #718096;">
                <i class="fas fa-inbox" style="font-size: 24px; margin-bottom: 10px; display: block;"></i>
                ${message}
            </td>
        </tr>
    `;
}