/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './src/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        // DeshChain Brand Colors
        primary: {
          50: '#fff7ed',
          100: '#ffedd5',
          200: '#fed7aa',
          300: '#fdba74',
          400: '#fb923c',
          500: '#f97316', // Saffron
          600: '#ea580c',
          700: '#c2410c',
          800: '#9a3412',
          900: '#7c2d12',
        },
        cultural: {
          saffron: '#ff9933',
          white: '#ffffff',
          green: '#138808',
          navy: '#000080',
          gold: '#ffd700',
        },
        lending: {
          krishi: '#22c55e',    // Green for agriculture
          vyavasaya: '#3b82f6', // Blue for business
          shiksha: '#8b5cf6',   // Purple for education
        },
        // Dark mode support
        dark: {
          bg: '#0f172a',
          surface: '#1e293b',
          border: '#334155',
          text: '#e2e8f0',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'Menlo', 'Monaco', 'monospace'],
        hindi: ['Noto Sans Devanagari', 'sans-serif'],
        bengali: ['Noto Sans Bengali', 'sans-serif'],
        tamil: ['Noto Sans Tamil', 'sans-serif'],
      },
      animation: {
        'fade-in': 'fadeIn 0.5s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'bounce-soft': 'bounceSoft 2s infinite',
        'pulse-slow': 'pulse 3s infinite',
        'shimmer': 'shimmer 2s infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(20px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        bounceSoft: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-5px)' },
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' },
        },
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-conic': 'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
        'cultural-pattern': "url('/images/cultural-pattern.svg')",
        'festival-fireworks': "url('/images/festival-bg.svg')",
      },
      backdropBlur: {
        xs: '2px',
      },
      boxShadow: {
        'cultural': '0 4px 20px rgba(255, 153, 51, 0.15)',
        'lending': '0 4px 20px rgba(59, 130, 246, 0.15)',
        'glow': '0 0 20px rgba(249, 115, 22, 0.3)',
      },
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
        '128': '32rem',
      },
      maxWidth: {
        '8xl': '88rem',
        '9xl': '96rem',
      },
      zIndex: {
        '60': '60',
        '70': '70',
        '80': '80',
        '90': '90',
        '100': '100',
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('@tailwindcss/aspect-ratio'),
    // Custom plugin for cultural utilities
    function({ addUtilities }) {
      const newUtilities = {
        '.text-gradient': {
          background: 'linear-gradient(45deg, #ff9933, #138808)',
          '-webkit-background-clip': 'text',
          '-webkit-text-fill-color': 'transparent',
          'background-clip': 'text',
        },
        '.border-gradient': {
          border: '2px solid',
          'border-image': 'linear-gradient(45deg, #ff9933, #138808) 1',
        },
        '.glass': {
          background: 'rgba(255, 255, 255, 0.1)',
          'backdrop-filter': 'blur(10px)',
          border: '1px solid rgba(255, 255, 255, 0.2)',
        },
        '.cultural-card': {
          background: 'linear-gradient(135deg, rgba(255, 153, 51, 0.05) 0%, rgba(19, 136, 8, 0.05) 100%)',
          border: '1px solid rgba(255, 153, 51, 0.2)',
          'backdrop-filter': 'blur(5px)',
        },
      }
      addUtilities(newUtilities)
    }
  ],
  darkMode: 'class',
}