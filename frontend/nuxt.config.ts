import tailwindcss from "@tailwindcss/vite"

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  modules: [
    '@nuxt/icon',
    'shadcn-nuxt',
    '@nuxt/devtools',
    '@nuxt/eslint',
    '@pinia/nuxt'
  ],


  icon: {
    //Load from client
    provider: 'iconify',
    serverBundle: false,
  },

  shadcn: {
    /**
     * Prefix for all the imported component.
     * @default "Ui"
     */
    prefix: '',
    /**
     * Directory that the component lives in.
     * Will respect the Nuxt aliases.
     * @link https://nuxt.com/docs/api/nuxt-config#alias
     * @default "@/components/ui"
     */
    componentDir: '@/components/ui'
  },

  eslint: {
    config: {
      standalone: true
    }
  },

  pinia: {
    storesDirs: ['./stores/**'],
  },

  css: ['~/assets/css/tailwind.css'],

  vite: {
    plugins: [
      tailwindcss(),
    ],
  },

  //Activate SPA
  ssr: false,
  //Activate Loading Screen when SPA starts, can be changed under app/spa-loading-template.html
  spaLoadingTemplate: true,

  runtimeConfig: {

  },

  

})