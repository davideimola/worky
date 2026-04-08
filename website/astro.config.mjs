import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import react from '@astrojs/react';
import vercel from '@astrojs/vercel';

export default defineConfig({
  site: 'https://worky.davideimola.dev',
  output: 'static',
  adapter: vercel(),
  integrations: [
    starlight({
      title: 'Worky',
      description: 'Build self-contained interactive workshops as a single Go binary.',
      logo: {
        src: './public/logo.svg',
        alt: 'Worky',
      },
      social: [
        { icon: 'github', label: 'GitHub', href: 'https://github.com/davideimola/worky' },
      ],
      expressiveCode: {
        themes: ['one-dark-pro'],
      },
      customCss: [
        './src/styles/fonts.css',
        './src/styles/base.css',
        './src/styles/starlight-theme.css',
      ],
      sidebar: [
        { label: 'Getting Started', link: '/getting-started/' },
        { label: 'CLI',             link: '/cli/' },
        { label: 'Runtime',         link: '/runtime/' },
        { label: 'Customization',   link: '/customization/' },
        {
          label: 'Reference',
          collapsed: false,
          items: [
            { label: 'Configuration', link: '/reference/configuration/' },
            { label: 'Chapters',      link: '/reference/chapters/' },
            {
              label: 'Checks',
              items: [
                { label: 'Built-in', link: '/reference/checks/built-in/' },
                { label: 'Patterns', link: '/reference/checks/patterns/' },
              ],
            },
          ],
        },
        { label: 'Troubleshooting', link: '/troubleshooting/' },
        { label: 'Showcase',        link: '/showcase/' },
      ],
      head: [
        {
          tag: 'link',
          attrs: { rel: 'icon', type: 'image/svg+xml', href: '/icon.svg' },
        },
        {
          tag: 'meta',
          attrs: { property: 'og:image', content: 'https://worky.davideimola.dev/images/og-image.png' },
        },
        {
          tag: 'meta',
          attrs: { name: 'twitter:card', content: 'summary_large_image' },
        },
      ],
    }),
    react(),
  ],
});
