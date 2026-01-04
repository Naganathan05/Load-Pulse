import { defineConfig } from 'vitepress'

// https://vitepress.vuejs.org/config/app-configs
export default defineConfig({
	title: 'Load-Pulse',
	description: 'Distributed load testing tool built in Go',
	themeConfig: {
		nav: [
			{ text: 'Home', link: '/' },
			{ text: 'Contributors', link: '/contributors' }
		]
	}
})
