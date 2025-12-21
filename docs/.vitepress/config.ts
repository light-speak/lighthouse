import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Lighthouse',
  description: 'A Go GraphQL Framework',
  base: '/lighthouse/',

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/lighthouse/logo.svg' }],
  ],

  themeConfig: {
    logo: '/logo.svg',

    nav: [
      { text: '指南', link: '/guide/getting-started' },
      { text: 'Schema', link: '/schema/basics' },
      { text: '功能', link: '/features/database' },
      { text: 'GitHub', link: 'https://github.com/light-speak/lighthouse' }
    ],

    sidebar: {
      '/guide/': [
        {
          text: '入门',
          items: [
            { text: '快速开始', link: '/guide/getting-started' },
            { text: 'CLI 命令', link: '/guide/cli' },
            { text: '项目结构', link: '/guide/project-structure' },
          ]
        }
      ],
      '/schema/': [
        {
          text: 'GraphQL Schema',
          items: [
            { text: '基础语法', link: '/schema/basics' },
            { text: '指令 Directives', link: '/schema/directives' },
            { text: 'Resolver 编写', link: '/schema/resolver' },
            { text: 'DataLoader', link: '/schema/dataloader' },
          ]
        }
      ],
      '/features/': [
        {
          text: '核心功能',
          items: [
            { text: '数据库', link: '/features/database' },
            { text: '数据库迁移', link: '/features/migration' },
            { text: '中间件与认证', link: '/features/auth' },
            { text: '健康检查', link: '/features/health' },
          ]
        },
        {
          text: '扩展功能',
          items: [
            { text: '异步任务队列', link: '/features/queue' },
            { text: '消息系统', link: '/features/messaging' },
            { text: '文件存储', link: '/features/storage' },
            { text: '实时推送', link: '/features/subscription' },
            { text: '监控与指标', link: '/features/metrics' },
          ]
        }
      ]
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/light-speak/lighthouse' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024 Light Speak'
    },

    search: {
      provider: 'local'
    },

    outline: {
      level: [2, 3],
      label: '目录'
    },

    docFooter: {
      prev: '上一页',
      next: '下一页'
    },

    lastUpdated: {
      text: '最后更新于'
    }
  }
})
