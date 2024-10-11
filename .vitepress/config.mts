import { DefaultTheme, defineConfig } from 'vitepress'

/**
 * 侧边栏配置: https://vitepress.dev/zh/reference/default-theme-sidebar
 */
const sidebar: DefaultTheme.Sidebar = [
  {
    collapsed: false,
    text: '准备工作',
    base: '/starter/',
    items: [
      { text: '安装', link: '/install' },
    ]
  },
  {
    collapsed: false,
    text: 'MySQL',
    base: '/mysql/',
    items: [
      { text: 'MySQL的数据类型', link: 'mysql-data-type' },
      { text: 'MySQL的SQL', link: 'mysql-sql' },
      { text: 'MySQL的索引', link: 'mysql-index' },
    ]
  },
]

/**
 * 主题配置: https://vitepress.dev/zh/reference/default-theme-config
 */
const themeConfig: DefaultTheme.Config = {
  // logo: '/motorcycle.svg',
  sidebar: sidebar, // 侧边栏配置
  socialLinks: [
    { icon: 'github', link: 'https://github.com/henryzhuhr/hello-database/' }
  ],
  darkModeSwitchLabel: '外观',          // 深色模式开关标签
  lightModeSwitchTitle: '切换到浅色模式', // 悬停时显示的浅色模式开关标题
  darkModeSwitchTitle: '切换到深色模式',  // 悬停时显示的深色模式开关标题
  sidebarMenuLabel: '目录',             // 侧边栏菜单的标题
  returnToTopLabel: '返回顶部',          // 返回顶部按钮的标题
  langMenuLabel: '选择语言',             // 语言选择菜单的标题
  externalLinkIcon: true,
  docFooter: {
    prev: '⏪️ 上一页',
    next: '下一页 ⏩️'
  },
  footer: {
    message: 'Powered By <a href="https://vitepress.dev/">Vitepress</a>',
    copyright: `All rights reserved © 2024-${new Date().getFullYear()} <a href="https://github.com/henryzhuhr?tab=repositories">henryzhuhr</a>`
  },
  outline: {
    label: '页面导航'
  },
  lastUpdated: {
    text: '⏰ 内容最后更新于',
    formatOptions: {
      dateStyle: 'short',
      timeStyle: 'medium'
    }
  },
  search: {   // 本地搜索: https://vitepress.dev/zh/reference/default-theme-search#local-search
    provider: 'local',
  },
}
// https://vitepress.dev/reference/site-config
export default defineConfig({
  srcDir: 'docs',
  base: '/hello-database/',
  title: "Hello Database",
  description: "Learning Database",
  themeConfig: themeConfig,
  lastUpdated: true,
  markdown: {
    math: true
  },
  vite: {// Vite 配置选项
    publicDir: '../.vitepress/public', // 相对于 docs 目录
    css: {
      preprocessorOptions: {
        scss: {
          api: "modern-compiler" // or 'modern'
        }
      }
    },
  },
})

