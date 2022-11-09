// import {sassPlugin} from 'esbuild-sass-plugin'
import { sassPlugin } from 'esbuild-sass-plugin';
import { build } from 'esbuild';

build({
  format: 'esm',
  entryPoints: ['index.scss'],
  outdir: 'server/public',
  plugins: [sassPlugin({
    // loadPaths: ['node_modules/@primer/css', 'test.scss'],
    // type: "style",
  })]
})
