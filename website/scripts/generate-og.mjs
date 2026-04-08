import sharp from 'sharp';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __dirname = dirname(fileURLToPath(import.meta.url));
const src = join(__dirname, '../public/images/og-image.svg');
const out = join(__dirname, '../public/images/og-image.png');

await sharp(src).png().toFile(out);
console.log('og-image.png generated');
