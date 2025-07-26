/* DeshChain Landing Page - Advanced Animations */

class AnimationController {
    constructor() {
        this.animations = new Map();
        this.observers = new Map();
        this.isReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
        
        this.init();
    }
    
    init() {
        if (!this.isReducedMotion) {
            this.setupParallaxEffects();
            this.setupHeroAnimations();
            this.setupCardAnimations();
            this.setupFloatingElements();
            this.setupBackgroundAnimations();
        }
        
        this.setupAccessibleAlternatives();
    }
    
    setupParallaxEffects() {
        const parallaxElements = document.querySelectorAll('.mandala-container, .cultural-pattern');
        
        const handleParallax = this.throttle(() => {
            const scrolled = window.pageYOffset;
            const parallax = scrolled * 0.5;
            
            parallaxElements.forEach(element => {
                element.style.transform = `translateY(${parallax}px)`;
            });
        }, 16); // ~60fps
        
        window.addEventListener('scroll', handleParallax);
    }
    
    setupHeroAnimations() {
        const heroContent = document.querySelector('.hero-content');
        const heroVisual = document.querySelector('.hero-visual');
        
        if (heroContent && heroVisual) {
            // Stagger animation for hero elements
            const heroElements = [
                '.hero-badge',
                '.hero-title',
                '.hero-subtitle',
                '.hero-stats',
                '.hero-actions'
            ];
            
            heroElements.forEach((selector, index) => {
                const element = heroContent.querySelector(selector);
                if (element) {
                    element.style.opacity = '0';
                    element.style.transform = 'translateY(30px)';
                    
                    setTimeout(() => {
                        element.style.transition = 'all 0.8s cubic-bezier(0.25, 0.46, 0.45, 0.94)';
                        element.style.opacity = '1';
                        element.style.transform = 'translateY(0)';
                    }, index * 200);
                }
            });
            
            // Hero visual animation
            setTimeout(() => {
                heroVisual.style.opacity = '1';
                heroVisual.style.transform = 'scale(1)';
            }, 600);
        }
    }
    
    setupCardAnimations() {
        const cardSelectors = [
            '.impact-card',
            '.module-card',
            '.community-card',
            '.infographic-item'
        ];
        
        cardSelectors.forEach(selector => {
            this.setupStaggeredAnimation(selector, {
                threshold: 0.2,
                stagger: 100,
                animation: 'slideUp'
            });
        });
    }
    
    setupStaggeredAnimation(selector, options = {}) {
        const elements = document.querySelectorAll(selector);
        const { threshold = 0.1, stagger = 100, animation = 'fadeIn' } = options;
        
        const observer = new IntersectionObserver((entries) => {
            entries.forEach((entry, index) => {
                if (entry.isIntersecting) {
                    setTimeout(() => {
                        this.playAnimation(entry.target, animation);
                    }, index * stagger);
                    observer.unobserve(entry.target);
                }
            });
        }, { threshold });
        
        elements.forEach(element => {
            this.prepareElement(element, animation);
            observer.observe(element);
        });
        
        this.observers.set(selector, observer);
    }
    
    prepareElement(element, animationType) {
        switch (animationType) {
            case 'fadeIn':
                element.style.opacity = '0';
                break;
            case 'slideUp':
                element.style.opacity = '0';
                element.style.transform = 'translateY(30px)';
                break;
            case 'slideRight':
                element.style.opacity = '0';
                element.style.transform = 'translateX(-30px)';
                break;
            case 'scale':
                element.style.opacity = '0';
                element.style.transform = 'scale(0.9)';
                break;
        }
    }
    
    playAnimation(element, animationType) {
        element.style.transition = 'all 0.6s cubic-bezier(0.25, 0.46, 0.45, 0.94)';
        
        switch (animationType) {
            case 'fadeIn':
                element.style.opacity = '1';
                break;
            case 'slideUp':
                element.style.opacity = '1';
                element.style.transform = 'translateY(0)';
                break;
            case 'slideRight':
                element.style.opacity = '1';
                element.style.transform = 'translateX(0)';
                break;
            case 'scale':
                element.style.opacity = '1';
                element.style.transform = 'scale(1)';
                break;
        }
        
        element.classList.add('animated');
    }
    
    setupFloatingElements() {
        const floatingElements = document.querySelectorAll('.hero-card, .token-logo');
        
        floatingElements.forEach(element => {
            element.style.animation = 'float 6s ease-in-out infinite';
        });
    }
    
    setupBackgroundAnimations() {
        // Animated gradient backgrounds
        const gradientElements = document.querySelectorAll('.hero, .technology, .whitepaper');
        
        gradientElements.forEach(element => {
            element.classList.add('animated-gradient');
        });
        
        // Pulsing cultural patterns
        const culturalPatterns = document.querySelectorAll('.pattern-lotus, .pattern-rangoli');
        culturalPatterns.forEach(pattern => {
            pattern.classList.add('cultural-pulse');
        });
    }
    
    setupAccessibleAlternatives() {
        // For users who prefer reduced motion
        if (this.isReducedMotion) {
            document.body.classList.add('reduced-motion');
            
            // Ensure all elements are visible without animation
            const allAnimatedElements = document.querySelectorAll('[style*="opacity"], [style*="transform"]');
            allAnimatedElements.forEach(element => {
                element.style.opacity = '1';
                element.style.transform = 'none';
            });
        }
    }
    
    // Utility methods
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
    }
    
    cleanup() {
        this.observers.forEach(observer => observer.disconnect());
        this.animations.clear();
        this.observers.clear();
    }
}

// Custom CSS animations
const animationStyles = `
/* Keyframe Animations */
@keyframes float {
    0%, 100% {
        transform: translateY(0px);
    }
    50% {
        transform: translateY(-10px);
    }
}

@keyframes pulse {
    0%, 100% {
        transform: scale(1);
        opacity: 0.8;
    }
    50% {
        transform: scale(1.05);
        opacity: 1;
    }
}

@keyframes shimmer {
    0% {
        background-position: -200% 0;
    }
    100% {
        background-position: 200% 0;
    }
}

@keyframes slideInFromLeft {
    0% {
        opacity: 0;
        transform: translateX(-50px);
    }
    100% {
        opacity: 1;
        transform: translateX(0);
    }
}

@keyframes slideInFromRight {
    0% {
        opacity: 0;
        transform: translateX(50px);
    }
    100% {
        opacity: 1;
        transform: translateX(0);
    }
}

@keyframes fadeInUp {
    0% {
        opacity: 0;
        transform: translateY(30px);
    }
    100% {
        opacity: 1;
        transform: translateY(0);
    }
}

@keyframes scaleIn {
    0% {
        opacity: 0;
        transform: scale(0.9);
    }
    100% {
        opacity: 1;
        transform: scale(1);
    }
}

@keyframes rotateIn {
    0% {
        opacity: 0;
        transform: rotate(-10deg) scale(0.9);
    }
    100% {
        opacity: 1;
        transform: rotate(0deg) scale(1);
    }
}

/* Animation Classes */
.animated-gradient {
    background-size: 200% 200%;
    animation: gradientShift 8s ease infinite;
}

@keyframes gradientShift {
    0% { background-position: 0% 50%; }
    50% { background-position: 100% 50%; }
    100% { background-position: 0% 50%; }
}

.cultural-pulse {
    animation: pulse 4s ease-in-out infinite;
}

.shimmer-effect {
    background: linear-gradient(
        90deg,
        rgba(255, 255, 255, 0) 0%,
        rgba(255, 255, 255, 0.2) 50%,
        rgba(255, 255, 255, 0) 100%
    );
    background-size: 200% 100%;
    animation: shimmer 2s infinite;
}

/* Hero Visual Enhancements */
.hero-visual {
    opacity: 0;
    transform: scale(0.9);
    transition: all 1s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.hero-visual.loaded {
    opacity: 1;
    transform: scale(1);
}

/* Mandala Rotation */
.mandala {
    animation: mandalaRotate 60s linear infinite;
}

.mandala-2 {
    animation: mandalaRotate 45s linear infinite reverse;
}

.mandala-3 {
    animation: mandalaRotate 90s linear infinite;
}

/* Hover Animations */
.hover-lift {
    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.hover-lift:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}

.hover-glow:hover {
    box-shadow: 0 0 20px rgba(255, 153, 51, 0.4);
}

/* Loading States */
.loading-skeleton {
    background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
    background-size: 200% 100%;
    animation: shimmer 2s infinite;
}

/* Scroll Triggered Animations */
.scroll-fade-in {
    opacity: 0;
    transform: translateY(20px);
    transition: all 0.6s ease;
}

.scroll-fade-in.visible {
    opacity: 1;
    transform: translateY(0);
}

.scroll-slide-right {
    opacity: 0;
    transform: translateX(-30px);
    transition: all 0.6s ease;
}

.scroll-slide-right.visible {
    opacity: 1;
    transform: translateX(0);
}

/* Stagger Delays */
.stagger-1 { animation-delay: 0.1s; }
.stagger-2 { animation-delay: 0.2s; }
.stagger-3 { animation-delay: 0.3s; }
.stagger-4 { animation-delay: 0.4s; }
.stagger-5 { animation-delay: 0.5s; }

/* Reduced Motion Styles */
.reduced-motion * {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
}

.reduced-motion .mandala,
.reduced-motion .cultural-pulse,
.reduced-motion .animated-gradient {
    animation: none !important;
}

/* Focus States for Accessibility */
.btn:focus,
.nav-link:focus {
    outline: 2px solid var(--saffron);
    outline-offset: 2px;
}

/* Print Styles */
@media print {
    .mandala,
    .cultural-pulse,
    .animated-gradient,
    [class*="float"],
    [class*="pulse"] {
        animation: none !important;
    }
}

/* Performance Optimizations */
.gpu-accelerated {
    transform: translateZ(0);
    will-change: transform;
}

.contain-layout {
    contain: layout;
}

.contain-paint {
    contain: paint;
}
`;

// Inject animation styles
const animationStyleSheet = document.createElement('style');
animationStyleSheet.textContent = animationStyles;
document.head.appendChild(animationStyleSheet);

// Performance monitoring
class PerformanceMonitor {
    constructor() {
        this.metrics = {
            fps: 0,
            frameCount: 0,
            lastTime: performance.now()
        };
        
        this.monitorFPS();
    }
    
    monitorFPS() {
        const calculateFPS = () => {
            const now = performance.now();
            this.metrics.frameCount++;
            
            if (now >= this.metrics.lastTime + 1000) {
                this.metrics.fps = Math.round((this.metrics.frameCount * 1000) / (now - this.metrics.lastTime));
                this.metrics.frameCount = 0;
                this.metrics.lastTime = now;
                
                // Adjust animations based on performance
                if (this.metrics.fps < 30) {
                    document.body.classList.add('low-performance');
                } else {
                    document.body.classList.remove('low-performance');
                }
            }
            
            requestAnimationFrame(calculateFPS);
        };
        
        requestAnimationFrame(calculateFPS);
    }
}

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    const animationController = new AnimationController();
    const performanceMonitor = new PerformanceMonitor();
    
    // Cleanup on page unload
    window.addEventListener('beforeunload', () => {
        animationController.cleanup();
    });
});

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { AnimationController, PerformanceMonitor };
}