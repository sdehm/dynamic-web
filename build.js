import { build } from 'esbuild';

build({
  format: 'esm',
  bundle: true,
  entryPoints: ['templates/index.js'],
  outdir: 'server/public',
  plugins: []
})
