/** @type {import('tailwindcss').Config} */
export default {
  // Disable preflight to avoid conflicts with Starlight's base styles
  corePlugins: {
    preflight: false,
  },
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  darkMode: ['class', '[data-theme="dark"]'],
  theme: {
    extend: {
      // ─── Brand Colors ───
      colors: {
        bg:       '#071122',
        midnight: '#0E1F3A',
        navy:     '#1A3560',
        surface:  '#0B1828',
        teal: {
          DEFAULT: '#00C2A0',
          deep:    '#00A086',
          ghost:   'rgba(0, 194, 160, 0.10)',
          border:  'rgba(0, 194, 160, 0.25)',
        },
        slate: {
          DEFAULT: 'rgba(255, 255, 255, 0.45)',
          light:   'rgba(255, 255, 255, 0.65)',
          bright:  'rgba(255, 255, 255, 0.88)',
        },
        border: {
          DEFAULT: 'rgba(255, 255, 255, 0.07)',
          mid:     'rgba(255, 255, 255, 0.12)',
        },
        code: {
          blue:   '#61AFEF',
          green:  '#98C379',
          purple: '#C678DD',
          yellow: '#E5C07B',
        },
        error:   '#FF5757',
      },

      // ─── Typography ───
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
      fontSize: {
        'hero': ['clamp(40px, 5.5vw, 72px)', { lineHeight: '1.06', letterSpacing: '-0.04em' }],
        'section': ['clamp(26px, 2.8vw, 38px)', { lineHeight: '1.15', letterSpacing: '-0.03em' }],
      },

      // ─── Layout ───
      maxWidth: {
        content: '1100px',
      },
      height: {
        nav: '56px',
      },

      // ─── Borders ───
      borderRadius: {
        sm:   '4px',
        DEFAULT: '6px',
        md:   '6px',
        lg:   '10px',
        xl:   '16px',
        full: '999px',
      },

      // ─── Animations ───
      keyframes: {
        'fade-up': {
          from: { opacity: '0', transform: 'translateY(16px)' },
          to:   { opacity: '1', transform: 'translateY(0)' },
        },
        'pulse-teal': {
          '0%, 100%': { boxShadow: '0 0 0 0 rgba(0, 194, 160, 0.5)' },
          '50%':      { boxShadow: '0 0 0 5px rgba(0, 194, 160, 0)' },
        },
        'pixel-on': {
          from: { opacity: '0' },
          to:   { opacity: '1' },
        },
        blink: {
          '0%, 100%': { opacity: '1' },
          '50%':      { opacity: '0' },
        },
      },
      animation: {
        'fade-up':    'fade-up 0.55s ease forwards',
        'pulse-teal': 'pulse-teal 2s ease-in-out infinite',
        'pixel-on':   'pixel-on 0.3s ease forwards',
        'blink':      'blink 1s step-end infinite',
      },

      // ─── Box Shadows ───
      boxShadow: {
        card: '0 1px 3px rgba(0,0,0,0.4), 0 4px 16px rgba(0,0,0,0.2)',
        teal: '0 0 20px rgba(0, 194, 160, 0.15)',
        'teal-sm': '0 0 8px rgba(0, 194, 160, 0.3)',
      },
    },
  },
  plugins: [],
};
