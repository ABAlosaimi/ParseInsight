// Application state
let opsChart = null;
let timeChart = null;

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    loadAvailableLibraries();
    setupEventListeners();
});

// Load available parser libraries from API
async function loadAvailableLibraries() {
    try {
        const response = await fetch('/api/libraries');
        const data = await response.json();

        const container = document.getElementById('librariesContainer');
        container.innerHTML = '';

        data.libraries.forEach(lib => {
            const label = document.createElement('label');
            label.innerHTML = `
                <input type="checkbox" value="${lib}" checked>
                ${lib}
            `;
            container.appendChild(label);
        });
    } catch (error) {
        console.error('Failed to load libraries:', error);
    }
}

// Setup event listeners
function setupEventListeners() {
    document.getElementById('runBenchmark').addEventListener('click', runBenchmark);
}

// Run benchmark
async function runBenchmark() {
    const message = document.getElementById('httpMessage').value.trim();
    const messageType = document.getElementById('messageType').value;
    const iterations = parseInt(document.getElementById('iterations').value);
    const concurrency = parseInt(document.getElementById('concurrency').value);

    // Get selected libraries
    const libraryCheckboxes = document.querySelectorAll('#librariesContainer input[type="checkbox"]:checked');
    const libraries = Array.from(libraryCheckboxes).map(cb => cb.value);

    // Validation
    if (!message) {
        alert('Please enter an HTTP message');
        return;
    }

    if (libraries.length === 0) {
        alert('Please select at least one parser library');
        return;
    }

    // Show loading
    showLoading(true);
    hideResults();

    try {
        const response = await fetch('/api/benchmark', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                message: message,
                message_type: messageType,
                iterations: iterations,
                concurrency: concurrency,
                libraries: libraries,
            }),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Benchmark failed');
        }

        const data = await response.json();
        displayResults(data);
    } catch (error) {
        alert('Error: ' + error.message);
    } finally {
        showLoading(false);
    }
}

// Show/hide loading indicator
function showLoading(show) {
    document.getElementById('loading').style.display = show ? 'block' : 'none';
    document.getElementById('runBenchmark').disabled = show;
}

// Hide results section
function hideResults() {
    document.getElementById('resultsSection').style.display = 'none';
}

// Display benchmark results
function displayResults(data) {
    // Show results section
    document.getElementById('resultsSection').style.display = 'block';

    // Display recommendation
    document.getElementById('recommendation').textContent = data.recommendation;

    // Display table
    displayTable(data.results);

    // Display charts
    displayCharts(data.results);

    // Scroll to results
    document.getElementById('resultsSection').scrollIntoView({ behavior: 'smooth' });
}

// Display results table
function displayTable(results) {
    const tbody = document.getElementById('resultsTableBody');
    tbody.innerHTML = '';

    results.forEach(result => {
        const row = document.createElement('tr');

        let statusBadge = '';
        if (!result.success) {
            statusBadge = '<span class="status-badge status-error">Failed</span>';
        } else if (result.winner) {
            statusBadge = '<span class="status-badge status-winner">Winner</span>';
        } else {
            statusBadge = '<span class="status-badge status-success">Success</span>';
        }

        const opsPerSec = result.success ? formatNumber(result.ops_per_second) : '-';
        const avgTime = result.success ? formatDuration(result.avg_time_per_parse) : '-';
        const memory = result.success ? formatBytes(result.memory_allocated) : '-';
        const allocs = result.success ? formatNumber(result.allocs_per_op) : '-';

        row.innerHTML = `
            <td><strong>${result.library}</strong></td>
            <td>${opsPerSec}</td>
            <td>${avgTime}</td>
            <td>${memory}</td>
            <td>${allocs}</td>
            <td>${statusBadge}${result.error ? '<br><small>' + result.error + '</small>' : ''}</td>
        `;

        tbody.appendChild(row);
    });
}

// Display charts
function displayCharts(results) {
    const successfulResults = results.filter(r => r.success);

    if (successfulResults.length === 0) {
        return;
    }

    const labels = successfulResults.map(r => r.library);
    const opsData = successfulResults.map(r => r.ops_per_second);
    const timeData = successfulResults.map(r => r.avg_time_per_parse / 1000); // Convert to microseconds

    // Operations per second chart
    if (opsChart) {
        opsChart.destroy();
    }

    const opsCtx = document.getElementById('opsChart').getContext('2d');
    opsChart = new Chart(opsCtx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Operations per Second',
                data: opsData,
                backgroundColor: successfulResults.map(r =>
                    r.winner ? 'rgba(251, 191, 36, 0.8)' : 'rgba(59, 130, 246, 0.8)'
                ),
                borderColor: successfulResults.map(r =>
                    r.winner ? 'rgb(251, 191, 36)' : 'rgb(59, 130, 246)'
                ),
                borderWidth: 2,
            }],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false,
                },
                title: {
                    display: true,
                    text: 'Operations per Second (Higher is Better)',
                    font: {
                        size: 14,
                        weight: 'bold',
                    },
                },
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        callback: function(value) {
                            return formatNumber(value);
                        },
                    },
                },
            },
        },
    });

    // Average time per parse chart
    if (timeChart) {
        timeChart.destroy();
    }

    const timeCtx = document.getElementById('timeChart').getContext('2d');
    timeChart = new Chart(timeCtx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Average Time per Parse (μs)',
                data: timeData,
                backgroundColor: successfulResults.map(r =>
                    r.winner ? 'rgba(16, 185, 129, 0.8)' : 'rgba(139, 92, 246, 0.8)'
                ),
                borderColor: successfulResults.map(r =>
                    r.winner ? 'rgb(16, 185, 129)' : 'rgb(139, 92, 246)'
                ),
                borderWidth: 2,
            }],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false,
                },
                title: {
                    display: true,
                    text: 'Average Time per Parse (Lower is Better)',
                    font: {
                        size: 14,
                        weight: 'bold',
                    },
                },
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        callback: function(value) {
                            return value.toFixed(2) + ' μs';
                        },
                    },
                },
            },
        },
    });
}

// Format number with commas
function formatNumber(num) {
    return num.toLocaleString('en-US', { maximumFractionDigits: 0 });
}

// Format duration in nanoseconds to human readable
function formatDuration(ns) {
    if (ns < 1000) {
        return ns.toFixed(0) + ' ns';
    } else if (ns < 1000000) {
        return (ns / 1000).toFixed(2) + ' μs';
    } else if (ns < 1000000000) {
        return (ns / 1000000).toFixed(2) + ' ms';
    } else {
        return (ns / 1000000000).toFixed(2) + ' s';
    }
}

// Format bytes to human readable
function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i];
}
