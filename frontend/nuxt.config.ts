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
    provider: 'server',
    customCollections: []
  },

  shadcn: {
    prefix: 'Ui',
    componentDir: './components/shadcn'
  },

  eslint: {
    config: {
      standalone: true
    }
  },

  pinia: {
    storesDirs: ['./stores/**'],
  },


})