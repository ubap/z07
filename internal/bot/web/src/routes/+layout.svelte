<script>
	import './layout.css';
	import "../app.css";
	import { onMount } from 'svelte';
	import { connect } from '$lib/socket.js';
	import { page } from '$app/stores';

	onMount(() => connect());

	const navItems = [
		{ name: 'Tools', href: '/tools', icon: '⚙️' },
		{ name: 'Character Stats', href: '/stats', icon: '📊' },
		{ name: 'Waypoints', href: '/waypoints', icon: '' }
	];
</script>

<div class="flex h-screen w-full">
	<!-- Sidebar -->
	<aside
		class="w-64 bg-slate-900 border-r border-slate-800 flex flex-col"
	>
		<div class="p-6">
			<h1
				class="text-xl font-bold text-orange-500 tracking-tight"
			>
				z07 Proxy
			</h1>
			<p
				class="text-xs text-slate-500 uppercase tracking-widest mt-1"
			>
				Internal v1.0
			</p>
		</div>
		<nav class="flex-1 px-4 space-y-2">
			{#each navItems as item}
				<a
					href={item.href}
					class="flex items-center gap-3 px-4 py-3 rounded-lg transition-all
          {$page.url.pathname === item.href
						? 'bg-orange-600 text-white'
						: 'text-slate-400 hover:bg-slate-800 hover:text-white'}"
				>
					<span>{item.icon}</span>
					<span class="font-medium">{item.name}</span>
				</a>
			{/each}
		</nav>
	</aside>
	<!-- Main Content -->
	<main class="flex-1 overflow-auto p-8"><div class="max-w-4xl mx-auto"><slot></slot></div></main>
</div>
