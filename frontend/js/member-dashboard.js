document.addEventListener('DOMContentLoaded', function() {
    // Check authentication and role
    if (!requireAuth() || !requireRole('member')) return;
    
    // Initialize page
    initializePage();
    
    // Load dashboard data
    loadDashboardData();
});

async function loadDashboardData() {
    try {
        const dashboardData = await apiCall('/dashboard');
        
        // Update stats
        updateStatsGrid(dashboardData);
        
        // Update charts
        setTimeout(() => {
            drawSavingsProgressChart(dashboardData);
            drawMonthlyBreakdownChart(dashboardData);
        }, 100);
        
        // Update transactions table
        updateTransactionsTable(dashboardData.transactions || []);
        
    } catch (error) {
        showError('Failed to load dashboard data: ' + error.message);
    }
}

function updateStatsGrid(data) {
    const statsGrid = document.getElementById('statsGrid');
    
    const stats = [
        {
            title: 'Current Balance',
            value: formatCurrency(data.savings.current_balance || 0),
            icon: 'fas fa-wallet',
            iconClass: 'balance',
            change: data.savings.current_balance > data.savings.monthly_commitment ? '+Good Progress' : 'Keep Saving'
        },
        {
            title: 'Monthly Commitment',
            value: formatCurrency(data.savings.monthly_commitment || 0),
            icon: 'fas fa-calendar-check',
            iconClass: 'commitment',
            change: 'Auto-deducted'
        },
        {
            title: 'Social Contributions',
            value: formatCurrency(data.savings.social_contributions || 0),
            icon: 'fas fa-heart',
            iconClass: 'contribution',
            change: 'RWF 5,000/month'
        },
        {
            title: 'Active Loans',
            value: data.activeLoans || 0,
            icon: 'fas fa-credit-card',
            iconClass: 'loans',
            change: data.activeLoans > 0 ? 'Manage carefully' : 'Available to apply'
        }
    ];
    
    statsGrid.innerHTML = stats.map(stat => `
        <div class="stat-card">
            <div class="stat-header">
                <h3>${stat.title}</h3>
                <div class="stat-icon ${stat.iconClass}">
                    <i class="${stat.icon}"></i>
                </div>
            </div>
            <div class="stat-value">${stat.value}</div>
            <div class="stat-change">
                <i class="fas fa-info-circle"></i>
                ${stat.change}
            </div>
        </div>
    `).join('');
}

function updateTransactionsTable(transactions) {
    const tbody = document.getElementById('transactionsBody');
    
    if (transactions.length === 0) {
        tbody.innerHTML = createEmptyTableRow(4, 'No transactions yet');
    } else {
        tbody.innerHTML = transactions.map(tx => `
            <tr>
                <td>${formatDate(tx.created_at)}</td>
                <td>
                    <span style="display: flex; align-items: center;">
                        <i class="fas fa-${tx.type === 'savings' ? 'piggy-bank' : 'heart'}" style="margin-right: 8px; color: #667eea;"></i>
                        ${tx.type.replace('_', ' ').toUpperCase()}
                    </span>
                </td>
                <td><strong>${formatCurrency(tx.amount)}</strong></td>
                <td>${tx.description}</td>
            </tr>
        `).join('');
    }
}

function drawSavingsProgressChart(data) {
    const canvas = document.getElementById('savingsProgressChart');
    if (!canvas) return;
    
    const ctx = canvas.getContext('2d');
    const width = canvas.width;
    const height = canvas.height;
    const margin = 40;
    
    const currentBalance = data.savings.current_balance || 0;
    const monthlyCommitment = data.savings.monthly_commitment || 0;
    const target = monthlyCommitment * 12; // Annual target
    
    // Clear canvas
    ctx.clearRect(0, 0, width, height);
    
    if (target === 0) {
        ctx.fillStyle = '#718096';
        ctx.font = '16px Inter';
        ctx.textAlign = 'center';
        ctx.fillText('Set your monthly commitment', width/2, height/2);
        return;
    }
    
    const progress = Math.min(currentBalance / target, 1);
    const progressWidth = (width - 2 * margin) * progress;
    
    // Draw background bar
    ctx.fillStyle = '#e2e8f0';
    ctx.fillRect(margin, height/2 - 20, width - 2 * margin, 40);
    
    // Draw progress bar
    ctx.fillStyle = '#48bb78';
    ctx.fillRect(margin, height/2 - 20, progressWidth, 40);
    
    // Draw labels
    ctx.fillStyle = '#2d3748';
    ctx.font = 'bold 16px Inter';
    ctx.textAlign = 'center';
    ctx.fillText('Annual Savings Progress', width/2, height/2 - 40);
    ctx.fillText(`${formatCurrency(currentBalance)} / ${formatCurrency(target)}`, width/2, height/2 + 50);
    ctx.fillText(`${Math.round(progress * 100)}%`, width/2, height/2 + 5);
}

function drawMonthlyBreakdownChart(data) {
    const canvas = document.getElementById('monthlyBreakdownChart');
    if (!canvas) return;
    
    const ctx = canvas.getContext('2d');
    const centerX = canvas.width / 2;
    const centerY = canvas.height / 2;
    const radius = 100;
    
    const monthlyCommitment = data.savings.monthly_commitment || 0;
    const socialContribution = 5000;
    const total = monthlyCommitment + socialContribution;
    
    // Clear canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    if (total === 0) {
        ctx.fillStyle = '#718096';
        ctx.font = '16px Inter';
        ctx.textAlign = 'center';
        ctx.fillText('No commitments set', centerX, centerY);
        return;
    }
    
    const commitmentAngle = (monthlyCommitment / total) * 2 * Math.PI;
    const socialAngle = (socialContribution / total) * 2 * Math.PI;
    
    // Draw commitment slice
    ctx.beginPath();
    ctx.moveTo(centerX, centerY);
    ctx.arc(centerX, centerY, radius, 0, commitmentAngle);
    ctx.closePath();
    ctx.fillStyle = '#667eea';
    ctx.fill();
    
    // Draw social slice
    ctx.beginPath();
    ctx.moveTo(centerX, centerY);
    ctx.arc(centerX, centerY, radius, commitmentAngle, commitmentAngle + socialAngle);
    ctx.closePath();
    ctx.fillStyle = '#ed8936';
    ctx.fill();
    
    // Draw labels
    ctx.fillStyle = '#2d3748';
    ctx.font = 'bold 12px Inter';
    ctx.textAlign = 'left';
    
    // Commitment label
    ctx.fillStyle = '#667eea';
    ctx.fillRect(centerX + radius + 20, centerY - 30, 15, 15);
    ctx.fillStyle = '#2d3748';
    ctx.fillText(`Voluntary: ${formatCurrency(monthlyCommitment)}`, centerX + radius + 40, centerY - 18);
    
    // Social label
    ctx.fillStyle = '#ed8936';
    ctx.fillRect(centerX + radius + 20, centerY - 5, 15, 15);
    ctx.fillStyle = '#2d3748';
    ctx.fillText(`Social: ${formatCurrency(socialContribution)}`, centerX + radius + 40, centerY + 7);
}