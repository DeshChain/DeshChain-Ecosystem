/* DeshChain Landing Page - Main JavaScript */

class DeshChainApp {
    constructor() {
        this.navbar = document.querySelector('.navbar');
        this.navLinks = document.querySelectorAll('.nav-link');
        this.navHamburger = document.getElementById('nav-hamburger');
        this.navMenu = document.getElementById('nav-menu');
        this.themeToggle = document.getElementById('theme-toggle');
        this.languageBtn = document.getElementById('language-btn');
        this.statNumbers = document.querySelectorAll('.stat-number[data-count]');
        
        this.isMenuOpen = false;
        this.currentTheme = localStorage.getItem('deshchain-theme') || 'light';
        this.currentLanguage = localStorage.getItem('deshchain-language') || 'en';
        
        this.init();
    }
    
    init() {
        this.setupScrollEffects();
        this.setupNavigation();
        this.setupThemeToggle();
        this.setupLanguageSelector();
        this.setupCountingAnimation();
        this.setupSmoothScrolling();
        this.applyStoredPreferences();
        this.setupIntersectionObserver();
    }
    
    setupScrollEffects() {
        let lastScrollY = window.scrollY;
        
        window.addEventListener('scroll', () => {
            const scrollY = window.scrollY;
            
            // Navbar scroll effect
            if (scrollY > 100) {
                this.navbar.classList.add('scrolled');
            } else {
                this.navbar.classList.remove('scrolled');
            }
            
            // Hide/show navbar on scroll
            if (scrollY > lastScrollY && scrollY > 200) {
                this.navbar.style.transform = 'translateY(-100%)';
            } else {
                this.navbar.style.transform = 'translateY(0)';
            }
            
            lastScrollY = scrollY;
        });
    }
    
    setupNavigation() {
        // Mobile menu toggle
        this.navHamburger?.addEventListener('click', () => {
            this.toggleMobileMenu();
        });
        
        // Close menu when clicking nav links
        this.navLinks.forEach(link => {
            link.addEventListener('click', () => {
                if (this.isMenuOpen) {
                    this.toggleMobileMenu();
                }
            });
        });
        
        // Close menu when clicking outside
        document.addEventListener('click', (e) => {
            if (this.isMenuOpen && !this.navMenu.contains(e.target) && !this.navHamburger.contains(e.target)) {
                this.toggleMobileMenu();
            }
        });
    }
    
    toggleMobileMenu() {
        this.isMenuOpen = !this.isMenuOpen;
        this.navMenu.classList.toggle('active');
        this.navHamburger.classList.toggle('active');
        document.body.classList.toggle('menu-open');
    }
    
    setupThemeToggle() {
        this.themeToggle?.addEventListener('click', () => {
            this.toggleTheme();
        });
    }
    
    toggleTheme() {
        this.currentTheme = this.currentTheme === 'light' ? 'dark' : 'light';
        this.applyTheme();
        localStorage.setItem('deshchain-theme', this.currentTheme);
    }
    
    applyTheme() {
        document.documentElement.setAttribute('data-theme', this.currentTheme);
        this.themeToggle.textContent = this.currentTheme === 'light' ? '🌙' : '☀️';
    }
    
    setupLanguageSelector() {
        const languages = {
            'en': { flag: '🇮🇳', name: 'EN' },
            'hi': { flag: '🇮🇳', name: 'हिं' },
            'bn': { flag: '🇮🇳', name: 'বাং' },
            'te': { flag: '🇮🇳', name: 'తె' },
            'mr': { flag: '🇮🇳', name: 'मरा' },
            'ta': { flag: '🇮🇳', name: 'தமி' },
            'gu': { flag: '🇮🇳', name: 'ગુ' },
            'kn': { flag: '🇮🇳', name: 'ಕನ್ನ' },
            'ml': { flag: '🇮🇳', name: 'മല' },
            'pa': { flag: '🇮🇳', name: 'ਪੰ' },
            'or': { flag: '🇮🇳', name: 'ଓଡ଼ି' },
            'as': { flag: '🇮🇳', name: 'অস' }
        };
        
        this.languageBtn?.addEventListener('click', () => {
            const langKeys = Object.keys(languages);
            const currentIndex = langKeys.indexOf(this.currentLanguage);
            const nextIndex = (currentIndex + 1) % langKeys.length;
            this.currentLanguage = langKeys[nextIndex];
            
            this.applyLanguage();
            localStorage.setItem('deshchain-language', this.currentLanguage);
        });
    }
    
    applyLanguage() {
        const languages = {
            'en': { flag: '🇮🇳', name: 'EN' },
            'hi': { flag: '🇮🇳', name: 'हिं' },
            'bn': { flag: '🇮🇳', name: 'বাং' },
            'te': { flag: '🇮🇳', name: 'తె' },
            'mr': { flag: '🇮🇳', name: 'मरा' },
            'ta': { flag: '🇮🇳', name: 'தமி' },
            'gu': { flag: '🇮🇳', name: 'ગુ' },
            'kn': { flag: '🇮🇳', name: 'ಕನ್ನ' },
            'ml': { flag: '🇮🇳', name: 'മല' },
            'pa': { flag: '🇮🇳', name: 'ਪੰ' },
            'or': { flag: '🇮🇳', name: 'ଓଡ଼ି' },
            'as': { flag: '🇮🇳', name: 'অস' }
        };
        
        const lang = languages[this.currentLanguage];
        if (lang && this.languageBtn) {
            this.languageBtn.textContent = `${lang.flag} ${lang.name}`;
        }
        
        // Apply language class to body for CSS styling
        document.body.className = document.body.className.replace(/\blang-\w+\b/g, '');
        document.body.classList.add(`lang-${this.currentLanguage}`);
    }
    
    setupCountingAnimation() {
        const animateCount = (element, start, end, duration) => {
            const startTime = performance.now();
            const range = end - start;
            
            const animate = (currentTime) => {
                const elapsed = currentTime - startTime;
                const progress = Math.min(elapsed / duration, 1);
                
                // Easing function
                const easeOutCubic = 1 - Math.pow(1 - progress, 3);
                const current = Math.floor(start + (range * easeOutCubic));
                
                element.textContent = current.toLocaleString();
                
                if (progress < 1) {
                    requestAnimationFrame(animate);
                }
            };
            
            requestAnimationFrame(animate);
        };
        
        // Trigger counting animation when stats come into view
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    const element = entry.target;
                    const endValue = parseInt(element.dataset.count);
                    animateCount(element, 0, endValue, 2000);
                    observer.unobserve(element);
                }
            });
        }, { threshold: 0.5 });
        
        this.statNumbers.forEach(stat => observer.observe(stat));
    }
    
    setupSmoothScrolling() {
        this.navLinks.forEach(link => {
            link.addEventListener('click', (e) => {
                const href = link.getAttribute('href');
                if (href.startsWith('#')) {
                    e.preventDefault();
                    const target = document.querySelector(href);
                    if (target) {
                        const offsetTop = target.offsetTop - 80; // Account for fixed navbar
                        window.scrollTo({
                            top: offsetTop,
                            behavior: 'smooth'
                        });
                    }
                }
            });
        });
    }
    
    setupIntersectionObserver() {
        // Add fade-in animation for sections
        const observerOptions = {
            threshold: 0.1,
            rootMargin: '0px 0px -50px 0px'
        };
        
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.classList.add('animate-in');
                }
            });
        }, observerOptions);
        
        // Observe all sections
        document.querySelectorAll('section').forEach(section => {
            observer.observe(section);
        });
        
        // Observe cards and other elements
        document.querySelectorAll('.impact-card, .module-card, .community-card, .infographic-item').forEach(card => {
            observer.observe(card);
        });
    }
    
    applyStoredPreferences() {
        this.applyTheme();
        this.applyLanguage();
    }
}

// Festival Theme Manager
class FestivalThemeManager {
    constructor() {
        this.themes = {
            'diwali': { 
                name: 'Diwali', 
                dates: ['10-20', '11-15'], 
                colors: ['#FFD700', '#FF4500'] 
            },
            'holi': { 
                name: 'Holi', 
                dates: ['03-08', '03-15'], 
                colors: ['#FF69B4', '#32CD32', '#FFD700'] 
            },
            'dussehra': { 
                name: 'Dussehra', 
                dates: ['09-25', '10-05'], 
                colors: ['#FF4500', '#FFD700'] 
            },
            'navratri': { 
                name: 'Navratri', 
                dates: ['09-15', '09-24'], 
                colors: ['#FF1493', '#FF8C00', '#9400D3'] 
            }
        };
        
        this.checkFestivalTheme();
    }
    
    checkFestivalTheme() {
        const today = new Date();
        const monthDay = `${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`;
        
        for (const [themeId, theme] of Object.entries(this.themes)) {
            const [startDate, endDate] = theme.dates;
            if (monthDay >= startDate && monthDay <= endDate) {
                this.applyFestivalTheme(themeId);
                break;
            }
        }
    }
    
    applyFestivalTheme(themeId) {
        document.body.classList.add(`theme-${themeId}`);
        
        // Add festival banner
        const banner = document.createElement('div');
        banner.className = 'festival-banner';
        banner.innerHTML = `🎉 Celebrating ${this.themes[themeId].name} with DeshChain! 🎉`;
        
        const navbar = document.querySelector('.navbar');
        navbar.parentNode.insertBefore(banner, navbar);
    }
}

// Utility Functions
const Utils = {
    debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    },
    
    throttle(func, limit) {
        let inThrottle;
        return function() {
            const args = arguments;
            const context = this;
            if (!inThrottle) {
                func.apply(context, args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        };
    },
    
    formatNumber(num) {
        if (num >= 1000000) {
            return (num / 1000000).toFixed(1) + 'M';
        }
        if (num >= 1000) {
            return (num / 1000).toFixed(1) + 'K';
        }
        return num.toString();
    }
};

// Additional navbar scroll styles
const navbarStyles = `
.navbar.scrolled {
    background: rgba(255, 255, 255, 0.98);
    box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
}

.navbar {
    transition: all 0.3s ease;
}

.nav-menu.active {
    display: flex;
    position: fixed;
    top: 80px;
    left: 0;
    right: 0;
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(20px);
    flex-direction: column;
    padding: 2rem;
    border-bottom: 1px solid var(--border);
    z-index: 999;
}

.nav-hamburger.active span:nth-child(1) {
    transform: rotate(45deg) translate(5px, 5px);
}

.nav-hamburger.active span:nth-child(2) {
    opacity: 0;
}

.nav-hamburger.active span:nth-child(3) {
    transform: rotate(-45deg) translate(7px, -6px);
}

body.menu-open {
    overflow: hidden;
}

/* Animation classes */
section {
    opacity: 0;
    transform: translateY(30px);
    transition: all 0.6s ease;
}

section.animate-in {
    opacity: 1;
    transform: translateY(0);
}

.impact-card,
.module-card,
.community-card,
.infographic-item {
    opacity: 0;
    transform: translateY(20px);
    transition: all 0.5s ease;
}

.impact-card.animate-in,
.module-card.animate-in,
.community-card.animate-in,
.infographic-item.animate-in {
    opacity: 1;
    transform: translateY(0);
}

@media (min-width: 769px) {
    .nav-menu.active {
        position: static;
        background: none;
        flex-direction: row;
        padding: 0;
        border: none;
    }
}
`;

// Inject additional styles
const additionalStyleSheet = document.createElement('style');
additionalStyleSheet.textContent = navbarStyles;
document.head.appendChild(additionalStyleSheet);

// Initialize everything when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new DeshChainApp();
    new FestivalThemeManager();
    
    // Add loading complete class
    document.body.classList.add('loaded');
});

// Handle page load
window.addEventListener('load', () => {
    document.body.classList.add('page-loaded');
});

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { DeshChainApp, FestivalThemeManager, Utils };
}