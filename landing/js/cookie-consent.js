/* DeshChain Cookie Consent Management */

class CookieConsent {
    constructor() {
        this.consentBanner = document.getElementById('cookie-consent');
        this.acceptBtn = document.getElementById('cookie-accept');
        this.declineBtn = document.getElementById('cookie-decline');
        
        this.init();
    }
    
    init() {
        // Check if consent already given
        if (!this.hasConsent()) {
            this.showBanner();
        }
        
        // Bind events
        this.acceptBtn?.addEventListener('click', () => this.acceptCookies());
        this.declineBtn?.addEventListener('click', () => this.declineCookies());
    }
    
    hasConsent() {
        return localStorage.getItem('deshchain-cookie-consent') !== null;
    }
    
    showBanner() {
        this.consentBanner?.classList.remove('hidden');
        
        // Animate in
        setTimeout(() => {
            this.consentBanner?.classList.add('show');
        }, 100);
    }
    
    hideBanner() {
        this.consentBanner?.classList.add('hiding');
        
        setTimeout(() => {
            this.consentBanner?.classList.add('hidden');
            this.consentBanner?.classList.remove('show', 'hiding');
        }, 300);
    }
    
    acceptCookies() {
        localStorage.setItem('deshchain-cookie-consent', 'accepted');
        localStorage.setItem('deshchain-consent-date', new Date().toISOString());
        
        this.hideBanner();
        this.enableAnalytics();
        
        // Dispatch event for other components
        document.dispatchEvent(new CustomEvent('cookiesAccepted'));
    }
    
    declineCookies() {
        localStorage.setItem('deshchain-cookie-consent', 'declined');
        localStorage.setItem('deshchain-consent-date', new Date().toISOString());
        
        this.hideBanner();
        
        // Dispatch event for other components
        document.dispatchEvent(new CustomEvent('cookiesDeclined'));
    }
    
    enableAnalytics() {
        // Only enable analytics if user consented
        if (this.getConsentStatus() === 'accepted') {
            // Initialize any analytics here
            console.log('Analytics enabled - user consented to cookies');
        }
    }
    
    getConsentStatus() {
        return localStorage.getItem('deshchain-cookie-consent');
    }
    
    revokeConsent() {
        localStorage.removeItem('deshchain-cookie-consent');
        localStorage.removeItem('deshchain-consent-date');
        this.showBanner();
    }
}

// Cookie Consent Styles
const consentStyles = `
.cookie-consent {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    background: rgba(26, 26, 26, 0.95);
    backdrop-filter: blur(20px);
    color: white;
    padding: 1.5rem;
    border-top: 2px solid var(--saffron);
    z-index: 10000;
    transform: translateY(100%);
    transition: transform 0.3s ease-out;
}

.cookie-consent.show {
    transform: translateY(0);
}

.cookie-consent.hiding {
    transform: translateY(100%);
}

.cookie-content {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 2rem;
    flex-wrap: wrap;
}

.cookie-content p {
    margin: 0;
    flex: 1;
    min-width: 300px;
}

.cookie-buttons {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-wrap: wrap;
}

.btn-accept {
    background: var(--gradient-primary);
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: var(--radius-lg);
    font-weight: 500;
    cursor: pointer;
    transition: var(--transition-base);
}

.btn-accept:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(255, 153, 51, 0.4);
}

.btn-decline {
    background: transparent;
    color: rgba(255, 255, 255, 0.8);
    border: 1px solid rgba(255, 255, 255, 0.3);
    padding: 0.75rem 1.5rem;
    border-radius: var(--radius-lg);
    font-weight: 500;
    cursor: pointer;
    transition: var(--transition-base);
}

.btn-decline:hover {
    background: rgba(255, 255, 255, 0.1);
    color: white;
}

.cookie-link {
    color: var(--saffron);
    text-decoration: underline;
    font-size: 0.875rem;
}

.cookie-link:hover {
    color: var(--marigold);
}

@media (max-width: 768px) {
    .cookie-content {
        flex-direction: column;
        text-align: center;
        gap: 1rem;
    }
    
    .cookie-buttons {
        justify-content: center;
        width: 100%;
    }
    
    .btn-accept,
    .btn-decline {
        flex: 1;
        min-width: 120px;
    }
}

.hidden {
    display: none !important;
}
`;

// Inject styles
const styleSheet = document.createElement('style');
styleSheet.textContent = consentStyles;
document.head.appendChild(styleSheet);

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new CookieConsent();
});

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = CookieConsent;
}