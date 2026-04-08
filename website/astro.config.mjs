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
      title: 'Docs',
      description: 'Build self-contained interactive workshops as a single Go binary.',
      logo: {
        src: './public/logo.svg',
        alt: 'Worky',
        replacesTitle: true,
      },
      components: {
        SiteTitle: './src/components/starlight/SiteTitle.astro',
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
        { label: 'Getting Started', link: '/docs/getting-started/' },
        { label: 'CLI',             link: '/docs/cli/' },
        { label: 'Runtime',         link: '/docs/runtime/' },
        { label: 'Customization',   link: '/docs/customization/' },
        {
          label: 'Reference',
          collapsed: false,
          items: [
            { label: 'Configuration', link: '/docs/reference/configuration/' },
            { label: 'Chapters',      link: '/docs/reference/chapters/' },
            {
              label: 'Checks',
              items: [
                { label: 'Built-in', link: '/docs/reference/checks/built-in/' },
                { label: 'Patterns', link: '/docs/reference/checks/patterns/' },
              ],
            },
          ],
        },
        { label: 'Troubleshooting', link: '/docs/troubleshooting/' },
        { label: 'Showcase',        link: '/docs/showcase/' },
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
