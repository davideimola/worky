import { defineCollection } from 'astro:content';
import { docsLoader } from '@astrojs/starlight/loaders';

export const collections = {
  docs: defineCollection({ loader: docsLoader() }),
};
