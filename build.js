import { build } from 'esbuild';

build({
  format: "esm",
  bundle: true,
  entryPoints: ["templates/index.ts"],
  outdir: "server/public",
  plugins: [],
});
