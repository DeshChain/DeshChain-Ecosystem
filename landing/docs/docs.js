// DeshChain Documentation JavaScript
document.addEventListener('DOMContentLoaded', function() {
    // Initialize search functionality
    initializeSearch();
    
    // Initialize module tabs
    initializeModuleTabs();
    
    // Initialize theme toggle
    initializeTheme();
    
    // Initialize smooth scrolling
    initializeSmoothScroll();
    
    // Initialize code highlighting
    initializeCodeHighlight();
});

// Search Functionality
function initializeSearch() {
    const searchInput = document.getElementById('doc-search');
    const searchBtn = document.querySelector('.search-btn');
    
    if (!searchInput) return;
    
    // Search data - would be loaded from a JSON file in production
    const searchData = [
        // Modules
        { title: "NAMO Token", type: "module", url: "./modules/namo.html", keywords: ["token", "utility", "native", "cultural"] },
        { title: "Trade Finance", type: "module", url: "./modules/tradefinance.html", keywords: ["basel", "swift", "mt700", "mt710", "regulatory"] },
        { title: "Money Order DEX", type: "module", url: "./modules/moneyorder.html", keywords: ["dex", "exchange", "traditional", "money order"] },
        { title: "Gram Pension", type: "module", url: "./modules/grampension.html", keywords: ["pension", "retirement", "rural", "guaranteed returns"] },
        { title: "Identity System", type: "module", url: "./modules/identity.html", keywords: ["did", "w3c", "kyc", "india stack", "aadhaar"] },
        { title: "Sikkebaaz", type: "module", url: "./modules/sikkebaaz.html", keywords: ["memecoin", "launchpad", "anti-dump", "desi"] },
        { title: "Charitable Trust", type: "module", url: "./modules/charitabletrust.html", keywords: ["charity", "donation", "social impact", "transparency"] },
        
        // Guides
        { title: "Installation Guide", type: "guide", url: "./guides/installation.html", keywords: ["install", "setup", "quick start"] },
        { title: "Node Setup", type: "guide", url: "./guides/node-setup.html", keywords: ["node", "validator", "setup", "configuration"] },
        { title: "Basel III Integration", type: "guide", url: "./guides/basel-iii-integration.html", keywords: ["basel", "regulatory", "compliance", "capital"] },
        { title: "SWIFT MT Processing", type: "guide", url: "./guides/swift-mt-processing.html", keywords: ["swift", "mt", "messaging", "banking"] },
        
        // API
        { title: "REST API", type: "api", url: "./api/rest.html", keywords: ["rest", "api", "http", "endpoints"] },
        { title: "gRPC API", type: "api", url: "./api/grpc.html", keywords: ["grpc", "protobuf", "rpc"] },
        { title: "WebSocket API", type: "api", url: "./api/websocket.html", keywords: ["websocket", "ws", "realtime", "events"] },
        
        // Resources
        { title: "Whitepaper", type: "resource", url: "./resources/whitepaper.html", keywords: ["whitepaper", "economics", "technical"] },
        { title: "Glossary", type: "resource", url: "./resources/glossary.html", keywords: ["terms", "definitions", "glossary"] },
        { title: "FAQ", type: "resource", url: "./resources/faq.html", keywords: ["faq", "questions", "help"] }
    ];
    
    // Search function
    function performSearch(query) {
        if (!query || query.length < 2) {
            hideSearchResults();
            return;
        }
        
        const lowerQuery = query.toLowerCase();
        const results = searchData.filter(item => {
            const titleMatch = item.title.toLowerCase().includes(lowerQuery);
            const keywordMatch = item.keywords.some(keyword => 
                keyword.toLowerCase().includes(lowerQuery)
            );
            return titleMatch || keywordMatch;
        });
        
        displaySearchResults(results, query);
    }
    
    // Display search results
    function displaySearchResults(results, query) {
        // Create or get results container
        let resultsContainer = document.getElementById('search-results');
        if (!resultsContainer) {
            resultsContainer = document.createElement('div');
            resultsContainer.id = 'search-results';
            resultsContainer.className = 'search-results';
            searchInput.parentElement.appendChild(resultsContainer);
        }
        
        if (results.length === 0) {
            resultsContainer.innerHTML = `
                <div class="no-results">
                    <p>No results found for "${query}"</p>
                </div>
            `;
        } else {
            const resultsHTML = results.map(result => `
                <a href="${result.url}" class="search-result-item">
                    <div class="result-type">${result.type}</div>
                    <div class="result-title">${highlightMatch(result.title, query)}</div>
                </a>
            `).join('');
            
            resultsContainer.innerHTML = `
                <div class="search-results-header">
                    <span>${results.length} results for "${query}"</span>
                    <button class="close-search" onclick="hideSearchResults()">×</button>
                </div>
                <div class="search-results-list">
                    ${resultsHTML}
                </div>
            `;
        }
        
        resultsContainer.classList.add('active');
    }
    
    // Highlight matching text
    function highlightMatch(text, query) {
        const regex = new RegExp(`(${query})`, 'gi');
        return text.replace(regex, '<mark>$1</mark>');
    }
    
    // Hide search results
    window.hideSearchResults = function() {
        const resultsContainer = document.getElementById('search-results');
        if (resultsContainer) {
            resultsContainer.classList.remove('active');
            setTimeout(() => {
                resultsContainer.remove();
            }, 300);
        }
    };
    
    // Event listeners
    searchInput.addEventListener('input', (e) => {
        const query = e.target.value.trim();
        performSearch(query);
    });
    
    searchInput.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            hideSearchResults();
            searchInput.blur();
        }
    });
    
    searchBtn.addEventListener('click', () => {
        const query = searchInput.value.trim();
        performSearch(query);
    });
    
    // Click outside to close
    document.addEventListener('click', (e) => {
        if (!e.target.closest('.search-container')) {
            hideSearchResults();
        }
    });
}

// Module Tabs Functionality
function initializeModuleTabs() {
    const tabButtons = document.querySelectorAll('.tab-btn');
    const moduleCategories = document.querySelectorAll('.module-category');
    
    if (tabButtons.length === 0) return;
    
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const category = button.getAttribute('data-category');
            
            // Update active button
            tabButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');
            
            // Update active category
            moduleCategories.forEach(cat => {
                if (cat.id === category) {
                    cat.classList.add('active');
                } else {
                    cat.classList.remove('active');
                }
            });
            
            // Animate entrance
            const activeCategory = document.getElementById(category);
            if (activeCategory) {
                const cards = activeCategory.querySelectorAll('.module-doc-card');
                cards.forEach((card, index) => {
                    card.style.opacity = '0';
                    card.style.transform = 'translateY(20px)';
                    setTimeout(() => {
                        card.style.transition = 'all 0.3s ease';
                        card.style.opacity = '1';
                        card.style.transform = 'translateY(0)';
                    }, index * 50);
                });
            }
        });
    });
}

// Theme Toggle
function initializeTheme() {
    const themeToggle = document.getElementById('theme-toggle');
    const currentTheme = localStorage.getItem('theme') || 'light';
    
    // Apply saved theme
    document.documentElement.setAttribute('data-theme', currentTheme);
    updateThemeIcon(currentTheme);
    
    if (themeToggle) {
        themeToggle.addEventListener('click', () => {
            const theme = document.documentElement.getAttribute('data-theme');
            const newTheme = theme === 'light' ? 'dark' : 'light';
            
            document.documentElement.setAttribute('data-theme', newTheme);
            localStorage.setItem('theme', newTheme);
            updateThemeIcon(newTheme);
        });
    }
    
    function updateThemeIcon(theme) {
        if (themeToggle) {
            themeToggle.textContent = theme === 'light' ? '🌙' : '☀️';
        }
    }
}

// Smooth Scrolling
function initializeSmoothScroll() {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                const offset = 80; // Navbar height
                const targetPosition = target.getBoundingClientRect().top + window.pageYOffset - offset;
                
                window.scrollTo({
                    top: targetPosition,
                    behavior: 'smooth'
                });
            }
        });
    });
}

// Code Highlighting (simplified)
function initializeCodeHighlight() {
    const codeBlocks = document.querySelectorAll('pre code');
    
    codeBlocks.forEach(block => {
        // Add line numbers
        const lines = block.textContent.split('\n');
        const numberedLines = lines.map((line, index) => {
            return `<span class="line-number">${index + 1}</span>${line}`;
        }).join('\n');
        
        block.innerHTML = numberedLines;
        
        // Add copy button
        const copyBtn = document.createElement('button');
        copyBtn.className = 'copy-code-btn';
        copyBtn.textContent = 'Copy';
        copyBtn.onclick = () => copyCode(block);
        
        block.parentElement.style.position = 'relative';
        block.parentElement.appendChild(copyBtn);
    });
}

// Copy code function
function copyCode(codeBlock) {
    const code = codeBlock.textContent.replace(/^\d+/gm, ''); // Remove line numbers
    navigator.clipboard.writeText(code).then(() => {
        const copyBtn = codeBlock.parentElement.querySelector('.copy-code-btn');
        copyBtn.textContent = 'Copied!';
        setTimeout(() => {
            copyBtn.textContent = 'Copy';
        }, 2000);
    });
}

// Add search results styles
const searchStyles = `
<style>
.search-results {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    margin-top: 0.5rem;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-lg);
    max-height: 400px;
    overflow-y: auto;
    opacity: 0;
    transform: translateY(-10px);
    transition: all 0.3s ease;
    z-index: 1000;
}

.search-results.active {
    opacity: 1;
    transform: translateY(0);
}

.search-results-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    border-bottom: 1px solid var(--border);
    font-weight: 500;
}

.close-search {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: var(--text-secondary);
    transition: color 0.2s;
}

.close-search:hover {
    color: var(--text-primary);
}

.search-result-item {
    display: block;
    padding: 1rem;
    border-bottom: 1px solid var(--border);
    text-decoration: none;
    color: var(--text-primary);
    transition: background-color 0.2s;
}

.search-result-item:hover {
    background: rgba(255, 153, 51, 0.05);
}

.search-result-item:last-child {
    border-bottom: none;
}

.result-type {
    font-size: 0.75rem;
    color: var(--primary);
    text-transform: uppercase;
    font-weight: 600;
    margin-bottom: 0.25rem;
}

.result-title {
    font-weight: 500;
}

.result-title mark {
    background: rgba(255, 153, 51, 0.3);
    color: inherit;
    padding: 0 2px;
    border-radius: 2px;
}

.no-results {
    padding: 2rem;
    text-align: center;
    color: var(--text-secondary);
}

.copy-code-btn {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    padding: 0.25rem 0.75rem;
    background: var(--primary);
    color: white;
    border: none;
    border-radius: var(--radius-sm);
    font-size: 0.75rem;
    cursor: pointer;
    transition: all 0.2s;
}

.copy-code-btn:hover {
    background: var(--accent);
    transform: translateY(-1px);
}

.line-number {
    display: inline-block;
    width: 2.5em;
    padding-right: 1em;
    color: #666;
    text-align: right;
    user-select: none;
}
</style>
`;

// Inject search styles
document.head.insertAdjacentHTML('beforeend', searchStyles);