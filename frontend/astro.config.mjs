import { defineConfig } from 'astro/config';

export default defineConfig({
  // Output to dist for Go embedding
  outDir: './dist',
  
  // Use static output
  output: 'static',
  
  // Vite config
  vite: {
    ssr: {
      external: ['node:events']
    }
  }
});
