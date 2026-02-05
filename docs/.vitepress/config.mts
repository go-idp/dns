import { defineConfig } from 'vitepress'

export default defineConfig({
  base: '/dns/',
  title: 'DNS CLI',
  description: 'Simple DNS Client and Server CLI tool written in Go',
  
  head: [
    ['link', { rel: 'icon', href: '/dns/favicon.ico' }]
  ],

  themeConfig: {
    // logo: '/logo.svg', // Uncomment and add logo.svg to public/ directory
    
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/' },
      { text: 'Examples', link: '/examples/' },
      { 
        text: 'GitHub', 
        link: 'https://github.com/go-idp/dns',
        target: '_blank'
      }
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Getting Started',
          items: [
            { text: 'Introduction', link: '/guide/' },
            { text: 'Installation', link: '/guide/installation' },
            { text: 'Quick Start', link: '/guide/quick-start' }
          ]
        },
        {
          text: 'DNS Client',
          items: [
            { text: 'Client Usage', link: '/guide/client' },
            { text: 'Supported Protocols', link: '/guide/client-protocols' }
          ]
        },
        {
          text: 'DNS Server',
          items: [
            { text: 'Server Usage', link: '/guide/server' },
            { text: 'Configuration', link: '/guide/configuration' },
            { text: 'DNS-over-TLS', link: '/guide/dot' },
            { text: 'DNS-over-HTTPS', link: '/guide/doh' },
            { text: 'DNS-over-QUIC', link: '/guide/doq' }
          ]
        }
      ],
      '/examples/': [
        {
          text: 'Examples',
          items: [
            { text: 'Basic Server', link: '/examples/basic-server' },
            { text: 'DoT Server', link: '/examples/dot-server' },
            { text: 'DoH and DoQ Server', link: '/examples/doh-doq-server' },
            { text: 'DoH and DoQ Client', link: '/examples/doh-doq-client' },
            { text: 'Configuration File', link: '/examples/config-file' }
          ]
        }
      ]
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/go-idp/dns' }
    ],

    search: {
      provider: 'local'
    },

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024 go-idp'
    }
  }
})
