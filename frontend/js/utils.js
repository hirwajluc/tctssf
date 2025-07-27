// API base URL
const API_BASE_URL = '/api';

// Generic API call function
async function apiCall(endpoint, options = {}) {
    const token = localStorage.getItem('authToken');
    
    const defaultOptions = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            ...(token && { 'Authorization': token })
        }
    };
    
    const finalOptions = { ...defaultOptions, ...options };
    
    // Merge headers properly
    if (options.headers) {
        finalOptions.headers = { ...defaultOptions.headers, ...options.headers };
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, finalOptions);
        
        // Handle different response types
        const contentType = response.headers.get('content-type');
        let data;
        
        if (contentType && contentType.includes('application/json')) {
            data = await response.json();
        } else {
            data = await response.text();
        }
        
        if (!response.ok) {
            // If it's JSON with an error message, use that
            if (typeof data === 'object' && data.error) {
                throw new Error(data.error);
            }
            // Otherwise use the status text or a generic message
            throw new Error(data || response.statusText || `HTTP ${response.status}`);
        }
        
        return data;
    } catch (error) {
        console.error('API call failed:', error);
        
        // Handle network errors
        if (error.name === 'TypeError' && error.message === 'Failed to fetch') {
            throw new Error('Network error. Please check your connection.');
        }
        
        // Handle authentication errors
        if (error.message.includes('Unauthorized') || error.message.includes('Invalid session')) {
            // Clear invalid session data
            localStorage.removeItem('authToken');
            localStorage.removeItem('currentUser');
            
            // Redirect to login if not already there
            if (!window.location.pathname.includes('index.html') && window.location.pathname !== '/') {
                window.location.href = '/';
            }
        }
        
        throw error;
    }
}

// Show error message
function showError(message) {
    const errorDiv = document.getElementById('errorMsg');
    if (errorDiv) {
        errorDiv.textContent = message;
        errorDiv.classList.remove('hidden');
        
        // Auto-hide after 5 seconds
        setTimeout(() => {
            hideError();
        }, 5000);
        
        // Scroll to top to make sure message is visible
        errorDiv.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
}

// Hide error message
function hideError() {
    const errorDiv = document.getElementById('errorMsg');
    if (errorDiv) {
        errorDiv.classList.add('hidden');
    }
}

// Show success message
function showSuccess(message) {
    const successDiv = document.getElementById('successMsg');
    if (successDiv) {
        successDiv.textContent = message;
        successDiv.classList.remove('hidden');
        
        // Auto-hide after 3 seconds
        setTimeout(() => {
            hideSuccess();
        }, 3000);
        
        // Scroll to top to make sure message is visible
        successDiv.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
}

// Hide success message
function hideSuccess() {
    const successDiv = document.getElementById('successMsg');
    if (successDiv) {
        successDiv.classList.add('hidden');
    }
}

// Show loading state on button
function showLoading(button) {
    if (button) {
        button.disabled = true;
        button.dataset.originalText = button.innerHTML;
        button.innerHTML = '<i class="fas fa-spinner fa-spin" style="margin-right: 8px;"></i>Loading...';
    }
}

// Hide loading state on button
function hideLoading(button, originalText = null) {
    if (button) {
        button.disabled = false;
        const text = originalText || button.dataset.originalText || 'Submit';
        button.innerHTML = text;
    }
}

// Format currency
function formatCurrency(amount) {
    return new Intl.NumberFormat('en-RW', {
        style: 'currency',
        currency: 'RWF',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0
    }).format(amount);
}

// Format number with commas
function formatNumber(num) {
    return new Intl.NumberFormat().format(num);
}

// Format date
function formatDate(dateString, options = {}) {
    const defaultOptions = {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    };
    
    const finalOptions = { ...defaultOptions, ...options };
    
    try {
        return new Date(dateString).toLocaleDateString('en-US', finalOptions);
    } catch (error) {
        return dateString;
    }
}

// Format datetime
function formatDateTime(dateString) {
    return formatDate(dateString, {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

// Validate email format
function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// Validate phone number (Rwanda format)
function isValidPhone(phone) {
    const phoneRegex = /^(\+250|250)?[0-9]{9}$/;
    return phoneRegex.test(phone.replace(/\s/g, ''));
}

// Sanitize input to prevent XSS
function sanitizeInput(input) {
    const div = document.createElement('div');
    div.textContent = input;
    return div.innerHTML;
}

// Debounce function for search inputs
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Copy text to clipboard
async function copyToClipboard(text) {
    try {
        await navigator.clipboard.writeText(text);
        showSuccess('Copied to clipboard!');
    } catch (error) {
        // Fallback for older browsers
        const textArea = document.createElement('textarea');
        textArea.value = text;
        document.body.appendChild(textArea);
        textArea.select();
        try {
            document.execCommand('copy');
            showSuccess('Copied to clipboard!');
        } catch (err) {
            showError('Failed to copy to clipboard');
        }
        document.body.removeChild(textArea);
    }
}

// Download data as CSV
function downloadCSV(data, filename) {
    const csv = data.map(row => row.map(field => `"${field}"`).join(',')).join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
}

// Confirm dialog
function confirmAction(message, callback) {
    if (confirm(message)) {
        callback();
    }
}

// Initialize tooltips (if using a tooltip library)
function initializeTooltips() {
    // Add tooltip initialization code if needed
}

// Initialize page
function initializePage() {
    // Common page initialization
    initializeTooltips();
    
    // Set focus on first input
    const firstInput = document.querySelector('input:not([type="hidden"])');
    if (firstInput) {
        firstInput.focus();
    }
}

// Handle form submission with loading state
function handleFormSubmit(form, submitHandler) {
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const submitBtn = form.querySelector('button[type="submit"]');
        const originalText = submitBtn.innerHTML;
        
        try {
            hideError();
            hideSuccess();
            showLoading(submitBtn);
            
            await submitHandler(e);
        } catch (error) {
            showError(error.message || 'An error occurred. Please try again.');
        } finally {
            hideLoading(submitBtn, originalText);
        }
    });
}

// Auto-initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', initializePage);

// Add these functions to the END of your existing utils.js file

/**
 * Get user initials from first and last name
 * @param {string} firstName - User's first name
 * @param {string} lastName - User's last name
 * @returns {string} User initials (e.g., "JD" for John Doe)
 */
function getUserInitials(firstName, lastName) {
    if (!firstName && !lastName) return '??';
    
    const first = firstName ? firstName.charAt(0).toUpperCase() : '';
    const last = lastName ? lastName.charAt(0).toUpperCase() : '';
    
    return first + last;
}

/**
 * Create empty table row with custom icon (enhanced version)
 * @param {number} colspan - Number of columns to span
 * @param {string} message - Message to display
 * @param {string} icon - FontAwesome icon class (optional)
 * @returns {string} HTML string for empty row
 */
function createEmptyTableRow(colspan, message = 'No data available', icon = 'fas fa-inbox') {
    return `
        <tr>
            <td colspan="${colspan}" style="text-align: center; padding: 40px; color: #718096;">
                <i class="${icon}" style="font-size: 24px; margin-bottom: 10px; display: block;"></i>
                ${message}
            </td>
        </tr>
    `;
}

/**
 * Validate Rwanda phone number (enhanced version)
 * @param {string} phone - Phone number to validate
 * @returns {boolean} True if valid Rwanda phone number
 */
function isValidRwandaPhone(phone) {
    if (!phone) return false;
    
    // Remove spaces and hyphens
    const cleanPhone = phone.replace(/[\s\-]/g, '');
    
    // Rwanda phone patterns:
    // +250XXXXXXXXX (12 digits total)
    // 250XXXXXXXXX (11 digits total)
    // 07XXXXXXXX or 08XXXXXXXX (10 digits total)
    const patterns = [
        /^\+250[0-9]{9}$/,  // +250 followed by 9 digits
        /^250[0-9]{9}$/,    // 250 followed by 9 digits
        /^0[78][0-9]{8}$/   // 07 or 08 followed by 8 digits
    ];
    
    return patterns.some(pattern => pattern.test(cleanPhone));
}

/**
 * Format phone number for display
 * @param {string} phone - Phone number to format
 * @returns {string} Formatted phone number
 */
function formatPhoneNumber(phone) {
    if (!phone) return '';
    
    // Remove all non-digit characters except +
    const cleaned = phone.replace(/[^\d+]/g, '');
    
    // If it starts with +250, format as +250 XXX XXX XXX
    if (cleaned.startsWith('+250')) {
        const number = cleaned.substring(4);
        if (number.length === 9) {
            return `+250 ${number.substring(0, 3)} ${number.substring(3, 6)} ${number.substring(6)}`;
        }
    }
    
    // If it starts with 250, format as +250 XXX XXX XXX
    if (cleaned.startsWith('250') && cleaned.length === 12) {
        const number = cleaned.substring(3);
        return `+250 ${number.substring(0, 3)} ${number.substring(3, 6)} ${number.substring(6)}`;
    }
    
    // If it starts with 07 or 08, format as +250 XXX XXX XXX
    if ((cleaned.startsWith('07') || cleaned.startsWith('08')) && cleaned.length === 10) {
        const number = cleaned.substring(1);
        return `+250 ${number.substring(0, 3)} ${number.substring(3, 6)} ${number.substring(6)}`;
    }
    
    return phone; // Return original if no pattern matches
}