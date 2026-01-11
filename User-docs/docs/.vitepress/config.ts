import { defineConfig } from 'vitepress'

// https://vitepress.vuejs.org/config/app-configs
export default defineConfig({
  title: 'Load-Pulse',
  description: 'Load testing tool built in Go which works based on the Raft consensus algorithm',
  base: '/Load-Pulse/',
  themeConfig: {
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Installation', link: '/install' },
      { text: 'Commands', link: '/commands' },
      { text: 'Results', link: '/results' },
      { text: 'Config Reference', link: '/config-reference' },
      { text: 'Contributors', link: '/contributors' }
    ],
    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Home', link: '/' },
          { text: 'Installation', link: '/install' }
        ]
      },
      {
        text: 'Usage',
        items: [
          { text: 'Commands', link: '/commands' },
          { text: 'Config Reference', link: '/config-reference' },
          { text: 'Results', link: '/results' }
        ]
      },
      {
        text: 'Community',
        items: [
          { text: 'Contributors', link: '/contributors' }
        ]
      }
    ]
  }
})