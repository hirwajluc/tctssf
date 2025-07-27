/*!
 * Chart.js v4.4.0
 * https://www.chartjs.org
 * (c) 2023 Chart.js Contributors
 * Released under the MIT License
 */

// This is a minimal Chart.js implementation for local use
// For full functionality, download from https://github.com/chartjs/Chart.js

(function (global, factory) {
    typeof exports === 'object' && typeof module !== 'undefined' ? factory(exports) :
    typeof define === 'function' && define.amd ? define(['exports'], factory) :
    (global = typeof globalThis !== 'undefined' ? globalThis : global || self, factory(global.Chart = {}));
})(this, (function (exports) {
    'use strict';

    // Chart.js Core
    class Chart {
        constructor(ctx, config) {
            this.ctx = ctx;
            this.config = config;
            this.data = config.data;
            this.options = config.options || {};
            this.type = config.type;
            this.chart = null;
            
            this.init();
        }

        init() {
            if (this.ctx && this.ctx.getContext) {
                this.chart = this.ctx.getContext('2d');
            } else if (typeof this.ctx === 'string') {
                const canvas = document.getElementById(this.ctx);
                if (canvas) {
                    this.chart = canvas.getContext('2d');
                }
            } else if (this.ctx.getContext) {
                this.chart = this.ctx.getContext('2d');
            }

            if (this.chart) {
                this.render();
            }
        }

        render() {
            const canvas = this.chart.canvas;
            const width = canvas.width;
            const height = canvas.height;
            
            // Clear canvas
            this.chart.clearRect(0, 0, width, height);
            
            // Draw based on chart type
            switch(this.type) {
                case 'doughnut':
                    this.renderDoughnut();
                    break;
                case 'bar':
                    this.renderBar();
                    break;
                case 'line':
                    this.renderLine();
                    break;
                default:
                    this.renderDefault();
            }
        }

        renderDoughnut() {
            const canvas = this.chart.canvas;
            const centerX = canvas.width / 2;
            const centerY = canvas.height / 2;
            const radius = Math.min(centerX, centerY) - 20;
            
            const data = this.data.datasets[0].data;
            const labels = this.data.labels;
            const colors = this.data.datasets[0].backgroundColor;
            
            const total = data.reduce((sum, value) => sum + value, 0);
            let currentAngle = -Math.PI / 2;
            
            // Draw segments
            data.forEach((value, index) => {
                const sliceAngle = (value / total) * 2 * Math.PI;
                
                this.chart.beginPath();
                this.chart.arc(centerX, centerY, radius, currentAngle, currentAngle + sliceAngle);
                this.chart.arc(centerX, centerY, radius * 0.6, currentAngle + sliceAngle, currentAngle, true);
                this.chart.closePath();
                this.chart.fillStyle = colors[index] || `hsl(${index * 60}, 70%, 50%)`;
                this.chart.fill();
                
                currentAngle += sliceAngle;
            });

            // Draw legend if enabled
            if (this.options.plugins && this.options.plugins.legend !== false) {
                this.drawLegend(labels, colors);
            }
        }

        renderBar() {
            const canvas = this.chart.canvas;
            const padding = 40;
            const chartWidth = canvas.width - (padding * 2);
            const chartHeight = canvas.height - (padding * 2);
            
            const data = this.data.datasets[0].data;
            const labels = this.data.labels;
            const colors = this.data.datasets[0].backgroundColor;
            
            const maxValue = Math.max(...data);
            const barWidth = chartWidth / data.length * 0.8;
            const barSpacing = chartWidth / data.length * 0.2;
            
            // Draw bars
            data.forEach((value, index) => {
                const barHeight = (value / maxValue) * chartHeight;
                const x = padding + (index * (barWidth + barSpacing)) + barSpacing / 2;
                const y = canvas.height - padding - barHeight;
                
                this.chart.fillStyle = colors[index] || `hsl(${index * 60}, 70%, 50%)`;
                this.chart.fillRect(x, y, barWidth, barHeight);
                
                // Draw label
                this.chart.fillStyle = '#333';
                this.chart.font = '12px Arial';
                this.chart.textAlign = 'center';
                this.chart.fillText(labels[index], x + barWidth / 2, canvas.height - padding + 20);
                
                // Draw value
                this.chart.fillText(value.toString(), x + barWidth / 2, y - 5);
            });
            
            // Draw axes
            this.chart.strokeStyle = '#ccc';
            this.chart.lineWidth = 1;
            this.chart.beginPath();
            this.chart.moveTo(padding, padding);
            this.chart.lineTo(padding, canvas.height - padding);
            this.chart.lineTo(canvas.width - padding, canvas.height - padding);
            this.chart.stroke();
        }

        renderLine() {
            const canvas = this.chart.canvas;
            const padding = 40;
            const chartWidth = canvas.width - (padding * 2);
            const chartHeight = canvas.height - (padding * 2);
            
            const data = this.data.datasets[0].data;
            const labels = this.data.labels;
            const borderColor = this.data.datasets[0].borderColor || '#3498db';
            const backgroundColor = this.data.datasets[0].backgroundColor || 'rgba(52, 152, 219, 0.1)';
            
            const maxValue = Math.max(...data);
            const minValue = Math.min(...data);
            const range = maxValue - minValue || 1;
            
            // Calculate points
            const points = data.map((value, index) => ({
                x: padding + (index / (data.length - 1)) * chartWidth,
                y: padding + ((maxValue - value) / range) * chartHeight
            }));
            
            // Draw fill area if enabled
            if (this.data.datasets[0].fill) {
                this.chart.fillStyle = backgroundColor;
                this.chart.beginPath();
                this.chart.moveTo(points[0].x, canvas.height - padding);
                points.forEach(point => {
                    this.chart.lineTo(point.x, point.y);
                });
                this.chart.lineTo(points[points.length - 1].x, canvas.height - padding);
                this.chart.closePath();
                this.chart.fill();
            }
            
            // Draw line
            this.chart.strokeStyle = borderColor;
            this.chart.lineWidth = this.data.datasets[0].borderWidth || 2;
            this.chart.beginPath();
            points.forEach((point, index) => {
                if (index === 0) {
                    this.chart.moveTo(point.x, point.y);
                } else {
                    this.chart.lineTo(point.x, point.y);
                }
            });
            this.chart.stroke();
            
            // Draw points
            points.forEach(point => {
                this.chart.fillStyle = borderColor;
                this.chart.beginPath();
                this.chart.arc(point.x, point.y, 4, 0, 2 * Math.PI);
                this.chart.fill();
            });
            
            // Draw labels
            this.chart.fillStyle = '#333';
            this.chart.font = '12px Arial';
            this.chart.textAlign = 'center';
            labels.forEach((label, index) => {
                if (points[index]) {
                    this.chart.fillText(label, points[index].x, canvas.height - padding + 20);
                }
            });
            
            // Draw axes
            this.chart.strokeStyle = '#ccc';
            this.chart.lineWidth = 1;
            this.chart.beginPath();
            this.chart.moveTo(padding, padding);
            this.chart.lineTo(padding, canvas.height - padding);
            this.chart.lineTo(canvas.width - padding, canvas.height - padding);
            this.chart.stroke();
        }

        renderDefault() {
            const canvas = this.chart.canvas;
            this.chart.fillStyle = '#f0f0f0';
            this.chart.fillRect(0, 0, canvas.width, canvas.height);
            
            this.chart.fillStyle = '#666';
            this.chart.font = '16px Arial';
            this.chart.textAlign = 'center';
            this.chart.fillText('Chart Type Not Supported', canvas.width / 2, canvas.height / 2);
        }

        drawLegend(labels, colors) {
            const canvas = this.chart.canvas;
            const legendY = canvas.height - 30;
            const legendItemWidth = canvas.width / labels.length;
            
            this.chart.font = '12px Arial';
            this.chart.textAlign = 'center';
            
            labels.forEach((label, index) => {
                const x = (index + 0.5) * legendItemWidth;
                
                // Draw color box
                this.chart.fillStyle = colors[index] || `hsl(${index * 60}, 70%, 50%)`;
                this.chart.fillRect(x - 6, legendY - 6, 12, 12);
                
                // Draw label
                this.chart.fillStyle = '#333';
                this.chart.fillText(label, x, legendY + 20);
            });
        }

        destroy() {
            if (this.chart && this.chart.canvas) {
                this.chart.clearRect(0, 0, this.chart.canvas.width, this.chart.canvas.height);
            }
        }

        update() {
            this.render();
        }
    }

    // Export Chart
    if (typeof module !== 'undefined' && module.exports) {
        module.exports = Chart;
    } else if (typeof window !== 'undefined') {
        window.Chart = Chart;
    }

    exports.Chart = Chart;
    exports.default = Chart;

}));

// Auto-initialize Chart as global
if (typeof window !== 'undefined') {
    window.Chart = window.Chart || exports.Chart;
}