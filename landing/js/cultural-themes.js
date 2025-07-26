/* DeshChain Cultural Themes - Dynamic Festival Management */

class CulturalThemeManager {
    constructor() {
        this.festivals = {
            'diwali': {
                name: 'Diwali',
                englishName: 'Festival of Lights',
                dates: this.calculateDiwaliDates(),
                colors: {
                    primary: '#FFD700',
                    secondary: '#FF4500',
                    accent: '#FF6347',
                    gradient: 'linear-gradient(135deg, #FFD700, #FF4500)'
                },
                symbols: ['🪔', '✨', '🎆', '🕯️'],
                greeting: {
                    hi: 'दीपावली की शुभकामनाएं!',
                    en: 'Happy Diwali!',
                    gu: 'દિવાળીની શુભેચ્છાઓ!',
                    mr: 'दिवाळीच्या हार्दिक शुभेच्छा!'
                },
                message: 'May this Diwali illuminate the path to financial prosperity for all!'
            },
            'holi': {
                name: 'Holi',
                englishName: 'Festival of Colors',
                dates: this.calculateHoliDates(),
                colors: {
                    primary: '#FF69B4',
                    secondary: '#32CD32',
                    accent: '#FFD700',
                    gradient: 'linear-gradient(135deg, #FF69B4, #32CD32, #FFD700)'
                },
                symbols: ['🎨', '🌈', '💚', '💛'],
                greeting: {
                    hi: 'होली की शुभकामनाएं!',
                    en: 'Happy Holi!',
                    gu: 'હોળીની શુભેચ્છાઓ!',
                    mr: 'होळीच्या हार्दिक शुभेच्छा!'
                },
                message: 'Let colors of joy paint your blockchain journey!'
            },
            'dussehra': {
                name: 'Dussehra',
                englishName: 'Victory of Good over Evil',
                dates: this.calculateDussehraDates(),
                colors: {
                    primary: '#FF4500',
                    secondary: '#FFD700',
                    accent: '#DC143C',
                    gradient: 'linear-gradient(135deg, #FF4500, #FFD700)'
                },
                symbols: ['🏹', '👑', '🔥', '⚔️'],
                greeting: {
                    hi: 'विजयादशमी की शुभकामनाएं!',
                    en: 'Happy Dussehra!',
                    gu: 'દશેરાની શુભેચ્છાઓ!',
                    mr: 'दसऱ्याच्या हार्दिक शुभेच्छा!'
                },
                message: 'Victory of transparency over corruption, decentralization over centralization!'
            },
            'navratri': {
                name: 'Navratri',
                englishName: 'Nine Nights of Divine',
                dates: this.calculateNavratriDates(),
                colors: {
                    primary: '#FF1493',
                    secondary: '#FF8C00',
                    accent: '#9400D3',
                    gradient: 'linear-gradient(135deg, #FF1493, #FF8C00, #9400D3)'
                },
                symbols: ['💃', '🕉️', '🌺', '🎭'],
                greeting: {
                    hi: 'नवरात्रि की शुभकामनाएं!',
                    en: 'Happy Navratri!',
                    gu: 'નવરાત્રીની શુભેચ્છાઓ!',
                    mr: 'नवरात्रीच्या हार्दिक शुभेच्छा!'
                },
                message: 'Nine nights of celebrating divine blockchain innovation!'
            },
            'karva-chauth': {
                name: 'Karva Chauth',
                englishName: 'Festival of Love & Devotion',
                dates: this.calculateKarvaChauthDates(),
                colors: {
                    primary: '#DC143C',
                    secondary: '#FFD700',
                    accent: '#FF69B4',
                    gradient: 'linear-gradient(135deg, #DC143C, #FFD700)'
                },
                symbols: ['🌙', '💖', '🕯️', '💍'],
                greeting: {
                    hi: 'करवा चौथ की शुभकामनाएं!',
                    en: 'Happy Karva Chauth!',
                    gu: 'કરવા ચૌથની શુભેચ્છાઓ!',
                    mr: 'करवा चौथीच्या हार्दिक शुभेच्छा!'
                },
                message: 'Devotion to innovation, commitment to financial inclusion!'
            },
            'ganesh-chaturthi': {
                name: 'Ganesh Chaturthi',
                englishName: 'Lord Ganesha Festival',
                dates: this.calculateGaneshChaturthiDates(),
                colors: {
                    primary: '#FF8C00',
                    secondary: '#FF4500',
                    accent: '#FFD700',
                    gradient: 'linear-gradient(135deg, #FF8C00, #FF4500)'
                },
                symbols: ['🐘', '🕉️', '🌺', '🍬'],
                greeting: {
                    hi: 'गणेश चतुर्थी की शुभकामनाएं!',
                    en: 'Happy Ganesh Chaturthi!',
                    gu: 'ગણેશ ચતુર્થીની શુભેચ્છાઓ!',
                    mr: 'गणेश चतुर्थीच्या हार्दिक शुभेच्छा!'
                },
                message: 'May Lord Ganesha remove all obstacles from your DeFi journey!'
            }
        };
        
        this.currentTheme = null;
        this.currentLanguage = localStorage.getItem('deshchain-language') || 'en';
        
        this.init();
    }
    
    init() {
        this.checkActiveFestival();
        this.setupThemeToggler();
        this.setupLanguageListener();
    }
    
    // Festival Date Calculations (simplified - in production would use lunar calendar APIs)
    calculateDiwaliDates() {
        const year = new Date().getFullYear();
        return [`${year}-10-20`, `${year}-11-15`]; // Approximate range
    }
    
    calculateHoliDates() {
        const year = new Date().getFullYear();
        return [`${year}-03-08`, `${year}-03-15`]; // Approximate range
    }
    
    calculateDussehraDates() {
        const year = new Date().getFullYear();
        return [`${year}-09-25`, `${year}-10-05`]; // Approximate range
    }
    
    calculateNavratriDates() {
        const year = new Date().getFullYear();
        return [`${year}-09-15`, `${year}-09-24`]; // Approximate range
    }
    
    calculateKarvaChauthDates() {
        const year = new Date().getFullYear();
        return [`${year}-10-28`, `${year}-11-02`]; // Approximate range
    }
    
    calculateGaneshChaturthiDates() {
        const year = new Date().getFullYear();
        return [`${year}-08-22`, `${year}-09-01`]; // Approximate range
    }
    
    checkActiveFestival() {
        const today = new Date();
        const currentDate = today.toISOString().split('T')[0];
        
        for (const [festivalId, festival] of Object.entries(this.festivals)) {
            const [startDate, endDate] = festival.dates;
            if (currentDate >= startDate && currentDate <= endDate) {
                this.activateFestivalTheme(festivalId);
                break;
            }
        }
    }
    
    activateFestivalTheme(festivalId) {
        const festival = this.festivals[festivalId];
        if (!festival) return;
        
        this.currentTheme = festivalId;
        
        // Apply theme colors
        this.applyThemeColors(festival.colors);
        
        // Add theme class to body
        document.body.classList.add(`theme-${festivalId}`);
        
        // Create festival banner
        this.createFestivalBanner(festival);
        
        // Add cultural patterns
        this.addCulturalPatterns(festivalId);
        
        // Update greeting message
        this.updateGreetingMessage(festival);
        
        // Store active theme
        localStorage.setItem('deshchain-active-festival', festivalId);
    }
    
    applyThemeColors(colors) {
        const root = document.documentElement;
        root.style.setProperty('--festival-primary', colors.primary);
        root.style.setProperty('--festival-secondary', colors.secondary);
        root.style.setProperty('--festival-accent', colors.accent);
        root.style.setProperty('--festival-gradient', colors.gradient);
    }
    
    createFestivalBanner(festival) {
        // Remove existing banner
        const existingBanner = document.querySelector('.festival-banner');
        if (existingBanner) {
            existingBanner.remove();
        }
        
        const banner = document.createElement('div');
        banner.className = 'festival-banner';
        banner.innerHTML = `
            <div class="festival-content">
                <span class="festival-symbols">${festival.symbols.join(' ')}</span>
                <span class="festival-greeting">${festival.greeting[this.currentLanguage]}</span>
                <span class="festival-message">${festival.message}</span>
            </div>
        `;
        
        // Insert before navbar
        const navbar = document.querySelector('.navbar');
        navbar.parentNode.insertBefore(banner, navbar);
        
        // Add animation
        setTimeout(() => {
            banner.classList.add('active');
        }, 100);
    }
    
    addCulturalPatterns(festivalId) {
        const hero = document.querySelector('.hero');
        if (!hero) return;
        
        const patterns = {
            'diwali': this.createDiyaPattern(),
            'holi': this.createColorSplashPattern(),
            'dussehra': this.createArrowPattern(),
            'navratri': this.createDancePattern(),
            'karva-chauth': this.createMoonPattern(),
            'ganesh-chaturthi': this.createGaneshPattern()
        };
        
        if (patterns[festivalId]) {
            hero.appendChild(patterns[festivalId]);
        }
    }
    
    createDiyaPattern() {
        const pattern = document.createElement('div');
        pattern.className = 'festival-pattern diya-pattern';
        pattern.innerHTML = '🪔'.repeat(20);
        return pattern;
    }
    
    createColorSplashPattern() {
        const pattern = document.createElement('div');
        pattern.className = 'festival-pattern color-splash-pattern';
        pattern.innerHTML = '🎨💚💛💙💜❤️'.repeat(5);
        return pattern;
    }
    
    createArrowPattern() {
        const pattern = document.createElement('div');
        pattern.className = 'festival-pattern arrow-pattern';
        pattern.innerHTML = '🏹'.repeat(15);
        return pattern;
    }
    
    createDancePattern() {
        const pattern = document.createElement('div');
        pattern.className = 'festival-pattern dance-pattern';
        pattern.innerHTML = '💃🕺'.repeat(10);
        return pattern;
    }
    
    createMoonPattern() {
        const pattern = document.createElement('div');
        pattern.className = 'festival-pattern moon-pattern';
        pattern.innerHTML = '🌙✨'.repeat(12);
        return pattern;
    }
    
    createGaneshPattern() {
        const pattern = document.createElement('div');
        pattern.className = 'festival-pattern ganesh-pattern';
        pattern.innerHTML = '🐘🕉️🌺'.repeat(8);
        return pattern;
    }
    
    updateGreetingMessage(festival) {
        const heroTitle = document.querySelector('.hero-title');
        if (heroTitle) {
            const greeting = document.createElement('div');
            greeting.className = 'festival-greeting-overlay';
            greeting.textContent = festival.greeting[this.currentLanguage];
            
            heroTitle.appendChild(greeting);
            
            // Auto-remove after 5 seconds
            setTimeout(() => {
                greeting.classList.add('fade-out');
                setTimeout(() => greeting.remove(), 1000);
            }, 5000);
        }
    }
    
    setupThemeToggler() {
        // Add festival theme selector button
        const navbar = document.querySelector('.nav-actions');
        if (navbar) {
            const themeSelector = document.createElement('button');
            themeSelector.className = 'festival-theme-selector';
            themeSelector.innerHTML = '🎭';
            themeSelector.title = 'Change Festival Theme';
            
            themeSelector.addEventListener('click', () => {
                this.showThemeSelector();
            });
            
            navbar.insertBefore(themeSelector, navbar.firstChild);
        }
    }
    
    showThemeSelector() {
        const modal = document.createElement('div');
        modal.className = 'theme-selector-modal';
        modal.innerHTML = `
            <div class="theme-selector-content">
                <h3>Choose Festival Theme</h3>
                <div class="theme-options">
                    ${Object.entries(this.festivals).map(([id, festival]) => `
                        <button class="theme-option ${this.currentTheme === id ? 'active' : ''}" 
                                data-theme="${id}">
                            <span class="theme-symbols">${festival.symbols.join(' ')}</span>
                            <span class="theme-name">${festival.name}</span>
                            <span class="theme-english">${festival.englishName}</span>
                        </button>
                    `).join('')}
                    <button class="theme-option ${!this.currentTheme ? 'active' : ''}" 
                            data-theme="default">
                        <span class="theme-symbols">🌟</span>
                        <span class="theme-name">Default</span>
                        <span class="theme-english">No Festival Theme</span>
                    </button>
                </div>
                <button class="close-modal">×</button>
            </div>
        `;
        
        document.body.appendChild(modal);
        
        // Event listeners
        modal.querySelector('.close-modal').addEventListener('click', () => {
            modal.remove();
        });
        
        modal.querySelectorAll('.theme-option').forEach(option => {
            option.addEventListener('click', () => {
                const themeId = option.dataset.theme;
                this.switchTheme(themeId);
                modal.remove();
            });
        });
        
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.remove();
            }
        });
    }
    
    switchTheme(themeId) {
        // Remove current theme
        if (this.currentTheme) {
            document.body.classList.remove(`theme-${this.currentTheme}`);
        }
        
        // Remove festival banner
        const banner = document.querySelector('.festival-banner');
        if (banner) banner.remove();
        
        // Remove festival patterns
        const patterns = document.querySelectorAll('.festival-pattern');
        patterns.forEach(pattern => pattern.remove());
        
        if (themeId === 'default') {
            this.currentTheme = null;
            localStorage.removeItem('deshchain-active-festival');
        } else {
            this.activateFestivalTheme(themeId);
        }
    }
    
    setupLanguageListener() {
        document.addEventListener('languageChanged', (e) => {
            this.currentLanguage = e.detail.language;
            if (this.currentTheme) {
                const festival = this.festivals[this.currentTheme];
                this.updateGreetingMessage(festival);
            }
        });
    }
    
    // Cultural wisdom quotes (context-aware)
    getCulturalQuote() {
        const quotes = {
            'diwali': [
                'सर्वे भवन्तु सुखिनः सर्वे सन्तु निरामयाः - May all be happy and healthy',
                'तमसो मा ज्योतिर्गमय - Lead us from darkness to light'
            ],
            'holi': [
                'वसुधैव कुटुम्बकम् - The world is one family',
                'सर्वे भद्राणि पश्यन्तु - May everyone see auspiciousness'
            ],
            'default': [
                'अहिंसा परमो धर्मः - Non-violence is the highest virtue',
                'सत्यमेव जयते - Truth alone triumphs',
                'कर्मण्येवाधिकारस्ते - You have the right to action alone'
            ]
        };
        
        const themeQuotes = quotes[this.currentTheme] || quotes.default;
        return themeQuotes[Math.floor(Math.random() * themeQuotes.length)];
    }
    
    // Regional celebrations
    getRegionalCelebration() {
        const regional = {
            'bengal': 'Durga Puja',
            'punjab': 'Vaisakhi',
            'kerala': 'Onam',
            'tamil': 'Pongal',
            'maharashtra': 'Gudi Padwa',
            'gujarat': 'Navratri',
            'rajasthan': 'Teej',
            'karnataka': 'Ugadi'
        };
        
        // Could be enhanced with geolocation or user preference
        return regional[Math.floor(Math.random() * Object.keys(regional).length)];
    }
}

// Cultural Theme Styles
const culturalThemeStyles = `
/* Festival Banner Styles */
.festival-banner {
    background: var(--festival-gradient, var(--gradient-cultural));
    color: white;
    text-align: center;
    padding: 0.75rem;
    font-weight: 500;
    position: relative;
    overflow: hidden;
    z-index: 999;
    transform: translateY(-100%);
    transition: transform 0.5s ease;
}

.festival-banner.active {
    transform: translateY(0);
}

.festival-content {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    flex-wrap: wrap;
}

.festival-symbols {
    font-size: 1.2rem;
}

.festival-greeting {
    font-weight: 600;
}

.festival-message {
    font-style: italic;
    opacity: 0.9;
}

/* Festival Patterns */
.festival-pattern {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    pointer-events: none;
    overflow: hidden;
    opacity: 0.1;
    font-size: 2rem;
    z-index: 1;
}

.diya-pattern {
    animation: twinkle 3s ease-in-out infinite;
}

.color-splash-pattern {
    animation: colorWave 4s ease-in-out infinite;
}

.arrow-pattern {
    animation: shoot 2s ease-in-out infinite;
}

.dance-pattern {
    animation: dance 3s ease-in-out infinite;
}

.moon-pattern {
    animation: moonGlow 4s ease-in-out infinite;
}

.ganesh-pattern {
    animation: bless 5s ease-in-out infinite;
}

/* Festival Theme Selector */
.festival-theme-selector {
    background: transparent;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: var(--radius-md);
    transition: var(--transition-fast);
}

.festival-theme-selector:hover {
    background: rgba(255, 153, 51, 0.1);
}

.theme-selector-modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10000;
}

.theme-selector-content {
    background: var(--surface);
    border-radius: var(--radius-xl);
    padding: 2rem;
    max-width: 600px;
    width: 90%;
    position: relative;
}

.theme-options {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
    margin-top: 1rem;
}

.theme-option {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 1rem;
    background: var(--background);
    border: 2px solid var(--border);
    border-radius: var(--radius-lg);
    cursor: pointer;
    transition: var(--transition-base);
}

.theme-option:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-md);
}

.theme-option.active {
    border-color: var(--primary);
    background: rgba(255, 153, 51, 0.1);
}

.theme-symbols {
    font-size: 2rem;
    margin-bottom: 0.5rem;
}

.theme-name {
    font-weight: 600;
    margin-bottom: 0.25rem;
}

.theme-english {
    font-size: 0.875rem;
    color: var(--text-secondary);
}

.close-modal {
    position: absolute;
    top: 1rem;
    right: 1rem;
    background: none;
    border: none;
    font-size: 2rem;
    cursor: pointer;
    color: var(--text-secondary);
}

/* Greeting Overlay */
.festival-greeting-overlay {
    position: absolute;
    top: -2rem;
    left: 0;
    right: 0;
    background: var(--festival-gradient, var(--gradient-primary));
    color: white;
    padding: 0.5rem 1rem;
    border-radius: var(--radius-lg);
    font-size: 1rem;
    text-align: center;
    animation: greetingFloat 5s ease-in-out;
}

.festival-greeting-overlay.fade-out {
    opacity: 0;
    transform: translateY(-20px);
    transition: all 1s ease;
}

/* Animations */
@keyframes twinkle {
    0%, 100% { opacity: 0.05; }
    50% { opacity: 0.15; }
}

@keyframes colorWave {
    0%, 100% { transform: rotate(0deg) scale(1); }
    50% { transform: rotate(180deg) scale(1.1); }
}

@keyframes shoot {
    0%, 100% { transform: translateX(0); }
    50% { transform: translateX(10px); }
}

@keyframes dance {
    0%, 100% { transform: rotate(0deg); }
    25% { transform: rotate(5deg); }
    75% { transform: rotate(-5deg); }
}

@keyframes moonGlow {
    0%, 100% { opacity: 0.1; filter: brightness(1); }
    50% { opacity: 0.2; filter: brightness(1.5); }
}

@keyframes bless {
    0%, 100% { transform: scale(1); }
    50% { transform: scale(1.05); }
}

@keyframes greetingFloat {
    0% { 
        opacity: 0; 
        transform: translateY(20px); 
    }
    10%, 90% { 
        opacity: 1; 
        transform: translateY(0); 
    }
    100% { 
        opacity: 0; 
        transform: translateY(-20px); 
    }
}

/* Responsive */
@media (max-width: 768px) {
    .festival-content {
        flex-direction: column;
        gap: 0.5rem;
    }
    
    .theme-options {
        grid-template-columns: 1fr;
    }
    
    .festival-pattern {
        font-size: 1.5rem;
    }
}

/* Reduced Motion */
@media (prefers-reduced-motion: reduce) {
    .festival-pattern {
        animation: none !important;
    }
    
    .festival-greeting-overlay {
        animation: none !important;
    }
}
`;

// Inject cultural theme styles
const culturalStyleSheet = document.createElement('style');
culturalStyleSheet.textContent = culturalThemeStyles;
document.head.appendChild(culturalStyleSheet);

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new CulturalThemeManager();
});

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = CulturalThemeManager;
}