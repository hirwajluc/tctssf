document.addEventListener('DOMContentLoaded', function() {
    // Check authentication and role
    if (!requireAuth() || !requireRole(['admin', 'superadmin'])) return;
    
    // Initialize page
    initializePage();
    
    // Load dashboard data
    loadAdminDashboardData();
});

async function loadAdminDashboardData() {
    try {
        const analyticsData = await apiCall('/admin/reports');
        
        // Update stats
        updateAdminStatsGrid(analyticsData);
        
        // Update charts
        setTimeout(() => {
            drawYearlyComparisonChart(analyticsData);
            drawFinancialChart(analyticsData);
        }, 100);
        
    } catch (error) {
        showError('Failed to load dashboard data: ' + error.message);
    }
}

// Updated the admin-dashboard.js file
// Replaced the updateAdminStatsGrid function with this corrected version:

function updateAdminStatsGrid(data) {
    const statsGrid = document.getElementById('statsGrid');
    
    const stats = [
        {
            title: 'Total Users',
            value: data.stats.total_users || 0,
            icon: 'fas fa-users',
            iconClass: 'balance',
            change: `${data.stats.active_users || 0} active`
        },
        {
            title: 'Personal Savings', // CORRECTED: Label clarified
            value: formatCurrency(data.stats.total_savings || 0),
            icon: 'fas fa-piggy-bank',
            iconClass: 'commitment',
            change: `${data.stats.total_members || 0} members`
        },
        {
            title: 'Social Fund', // CORRECTED: Now shows only social contributions
            value: formatCurrency(data.stats.total_social_funds || 0),
            icon: 'fas fa-heart',
            iconClass: 'contribution',
            change: 'Community contributions only' // CORRECTED: Clarified description
        },
        {
            title: 'Pending Loans',
            value: data.stats.pending_loans || 0,
            icon: 'fas fa-clock',
            iconClass: 'loans',
            change: `${data.stats.total_loans || 0} total loans`
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

function drawYearlyComparisonChart(data) {
    const canvas = document.getElementById('yearlyComparisonChart');
    if (!canvas) return;
    
    const ctx = canvas.getContext('2d');
    const width = canvas.width;
    const height = canvas.height;
    const margin = 60;
    const chartWidth = width - 2 * margin;
    const chartHeight = height - 2 * margin;
    
    const thisYear = data.stats.this_year_funds || 0;
    const lastYear = data.stats.last_year_funds || 0;
    const maxValue = Math.max(thisYear, lastYear, 1000000); // Minimum 1M for scale
    
    const barWidth = chartWidth / 4; // 2 bars with spacing
    
    // Clear canvas
    ctx.clearRect(0, 0, width, height);
    
    // Draw grid lines
    ctx.strokeStyle = '#e2e8f0';
    ctx.lineWidth = 1;
    for (let i = 0; i <= 5; i++) {
        const y = margin + (i * chartHeight / 5);
        ctx.beginPath();
        ctx.moveTo(margin, y);
        ctx.lineTo(margin + chartWidth, y);
        ctx.stroke();
    }
    
    // Draw Y-axis labels
    ctx.fillStyle = '#718096';
    ctx.font = '12px Inter';
    ctx.textAlign = 'right';
    for (let i = 0; i <= 5; i++) {
        const value = maxValue - (i * maxValue / 5);
        const y = margin + (i * chartHeight / 5);
        ctx.fillText(formatCurrency(value), margin - 10, y + 4);
    }
    
    // Draw This Year bar
    const thisYearHeight = (thisYear / maxValue) * chartHeight;
    const thisYearX = margin + barWidth * 0.5;
    const thisYearY = margin + chartHeight - thisYearHeight;
    
    ctx.fillStyle = '#667eea';
    ctx.fillRect(thisYearX, thisYearY, barWidth, thisYearHeight);
    
    // Draw Last Year bar
    const lastYearHeight = (lastYear / maxValue) * chartHeight;
    const lastYearX = margin + barWidth * 2;
    const lastYearY = margin + chartHeight - lastYearHeight;
    
    ctx.fillStyle = '#48bb78';
    ctx.fillRect(lastYearX, lastYearY, barWidth, lastYearHeight);
    
    // Draw values on top of bars
    ctx.fillStyle = '#2d3748';
    ctx.font = 'bold 12px Inter';
    ctx.textAlign = 'center';
    
    if (thisYearHeight > 20) {
        ctx.fillText(formatCurrency(thisYear), thisYearX + barWidth/2, thisYearY - 8);
    }
    if (lastYearHeight > 20) {
        ctx.fillText(formatCurrency(lastYear), lastYearX + barWidth/2, lastYearY - 8);
    }
    
    // Draw X-axis labels
    ctx.fillStyle = '#2d3748';
    ctx.font = '14px Inter';
    ctx.fillText('2025', thisYearX + barWidth/2, margin + chartHeight + 25);
    ctx.fillText('2024', lastYearX + barWidth/2, margin + chartHeight + 25);
    
    // Draw percentage change
    if (lastYear > 0) {
        const change = ((thisYear - lastYear) / lastYear * 100).toFixed(1);
        const changeText = change >= 0 ? `+${change}%` : `${change}%`;
        const changeColor = change >= 0 ? '#48bb78' : '#fc8181';
        
        ctx.fillStyle = changeColor;
        ctx.font = 'bold 14px Inter';
        ctx.textAlign = 'center';
        ctx.fillText(`${changeText} growth`, width/2, margin - 20);
    }
}

// Updated the drawFinancialChart function in admin-dashboard.js
// Replaced the existing function with this corrected version:

function drawFinancialChart(data) {
    const canvas = document.getElementById('financialChart');
    if (!canvas) return;
    
    const ctx = canvas.getContext('2d');
    const width = canvas.width;
    const height = canvas.height;
    const margin = 60;
    const chartWidth = width - 2 * margin;
    const chartHeight = height - 2 * margin;
    
    // CORRECTED: Updated values and labels for clarity
    const values = [
        data.stats.total_savings || 0,        // Personal savings only
        data.stats.total_social_funds || 0,   // Social contributions only  
        data.stats.total_loan_amount || 0     // Total loans disbursed
    ];
    
    // CORRECTED: Updated labels to be more specific
    const labels = ['Personal Savings', 'Social Fund', 'Loans Disbursed'];
    const colors = ['#48bb78', '#ed8936', '#667eea'];
    
    const maxValue = Math.max(...values, 100000); // Minimum scale
    const barWidth = chartWidth / (labels.length * 2);
    
    // Clear canvas
    ctx.clearRect(0, 0, width, height);
    
    // Draw grid lines
    ctx.strokeStyle = '#e2e8f0';
    ctx.lineWidth = 1;
    for (let i = 0; i <= 5; i++) {
        const y = margin + (i * chartHeight / 5);
        ctx.beginPath();
        ctx.moveTo(margin, y);
        ctx.lineTo(margin + chartWidth, y);
        ctx.stroke();
    }
    
    // Draw Y-axis labels
    ctx.fillStyle = '#718096';
    ctx.font = '12px Inter';
    ctx.textAlign = 'right';
    for (let i = 0; i <= 5; i++) {
        const value = maxValue - (i * maxValue / 5);
        const y = margin + (i * chartHeight / 5);
        ctx.fillText(formatCurrency(value), margin - 10, y + 4);
    }
    
    // Draw bars
    labels.forEach((label, index) => {
        const barHeight = Math.max((values[index] / maxValue) * chartHeight, 2);
        const x = margin + (index * 2 + 0.5) * barWidth;
        const y = margin + chartHeight - barHeight;
        
        ctx.fillStyle = colors[index];
        ctx.fillRect(x, y, barWidth, barHeight);
        
        // Draw value on top of bar
        if (barHeight > 20) {
            ctx.fillStyle = '#2d3748';
            ctx.font = 'bold 12px Inter';
            ctx.textAlign = 'center';
            ctx.fillText(formatCurrency(values[index]), x + barWidth/2, y - 8);
        }
        
        // Draw label
        ctx.fillStyle = '#2d3748';
        ctx.font = '14px Inter';
        ctx.fillText(label, x + barWidth/2, margin + chartHeight + 25);
    });
    
    // Add chart title
    ctx.fillStyle = '#2d3748';
    ctx.font = 'bold 16px Inter';
    ctx.textAlign = 'center';
    ctx.fillText('Financial Overview', width/2, margin - 20);
}