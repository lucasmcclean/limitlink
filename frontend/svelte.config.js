import { vitePreprocess } from '@sveltejs/vite-plugin-svelte'
import sveltePreprocess from 'svelte-preprocess';

export default {
  preprocess: [
    vitePreprocess(),
    sveltePreprocess({
      typescript: {
        transpileOnly: true,
      },
    })
  ],
}
