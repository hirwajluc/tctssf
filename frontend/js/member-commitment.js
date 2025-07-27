document.addEventListener('DOMContentLoaded', function() {
    // Check authentication and role
    if (!requireAuth() || !requireRole('member')) return;
    
    // Initialize page
    initializePage();
    
    // Load member data
    loadMemberData();
    
    // Setup form handler
    setupCommitmentForm();
});

async function loadMemberData() {
    try {
        const dashboardData = await apiCall('/dashboard');
        const currentCommitment = dashboardData.savings.monthly_commitment || 0;
        
        // Update current commitment display
        document.getElementById('currentCommitment').value = currentCommitment;
        
        // Update deduction summary
        updateDeductionSummary(currentCommitment);
        
        // Set minimum date to next month
        setMinimumDate();
        
    } catch (error) {
        showError('Failed to load member data: ' + error.message);
    }
}

function updateDeductionSummary(commitment) {
    const voluntaryAmount = document.getElementById('voluntaryAmount');
    const totalDeduction = document.getElementById('totalDeduction');
    
    if (voluntaryAmount) {
        voluntaryAmount.textContent = formatCurrency(commitment);
    }
    
    if (totalDeduction) {
        totalDeduction.textContent = formatCurrency(commitment + 5000); // 5k social contribution
    }
}

function setMinimumDate() {
    const effectiveDateInput = document.getElementById('effectiveDate');
    if (effectiveDateInput) {
        const nextMonth = new Date();
        nextMonth.setMonth(nextMonth.getMonth() + 1);
        nextMonth.setDate(1);
        
        const minDate = nextMonth.toISOString().split('T')[0];
        effectiveDateInput.min = minDate;
        effectiveDateInput.value = minDate;
    }
}

function setupCommitmentForm() {
    const commitmentForm = document.getElementById('commitmentForm');
    
    if (commitmentForm) {
        commitmentForm.addEventListener('submit', handleCommitmentUpdate);
    }
    
    // Real-time update of deduction summary
    const newCommitmentInput = document.getElementById('newCommitment');
    if (newCommitmentInput) {
        newCommitmentInput.addEventListener('input', function() {
            const newValue = parseFloat(this.value) || 0;
            updateDeductionSummary(newValue);
        });
    }
}

async function handleCommitmentUpdate(e) {
    e.preventDefault();
    
    const submitBtn = e.target.querySelector('button[type="submit"]');
    const originalText = submitBtn.innerHTML;
    
    // Show loading state
    showLoading(submitBtn);
    
    const formData = new FormData(e.target);
    const commitmentData = {
        new_commitment: parseFloat(formData.get('new_commitment')),
        effective_date: formData.get('effective_date')
    };
    
    // Validate minimum commitment
    if (commitmentData.new_commitment < 5000) {
        showError('Minimum commitment is RWF 5,000');
        hideLoading(submitBtn, originalText);
        return;
    }
    
    // Validate future date
    const effectiveDate = new Date(commitmentData.effective_date);
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    
    if (effectiveDate < tomorrow) {
        showError('Effective date must be at least tomorrow');
        hideLoading(submitBtn, originalText);
        return;
    }
    
    try {
        await apiCall('/savings/update-commitment', {
            method: 'POST',
            body: JSON.stringify(commitmentData)
        });
        
        showSuccess('Monthly commitment updated successfully! Changes will take effect from the specified date.');
        
        // Reload data to show updated values
        setTimeout(() => {
            loadMemberData();
        }, 1000);
        
    } catch (error) {
        showError(error.message || 'Failed to update commitment');
    } finally {
        hideLoading(submitBtn, originalText);
    }
}