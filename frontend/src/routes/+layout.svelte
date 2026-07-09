<script lang="ts">
	import { page } from '$app/stores';
	import '../app.css';

	let { children }: { children: import('svelte').Snippet } = $props();

	let sidebarOpen = $state(false);

	const navItems = [
		{ href: '/', label: 'Dashboard', icon: 'LayoutDashboard' },
		{ href: '/messages', label: 'Messages', icon: 'MessageSquare' },
		{ href: '/devices', label: 'Devices', icon: 'Smartphone' },
		{ href: '/webhooks', label: 'Webhooks', icon: 'Link' },
		{ href: '/settings', label: 'Settings', icon: 'Settings' },
		{ href: '/tokens', label: 'API Tokens', icon: 'Key' },
	];

	function navClass(href: string) {
		const active =
			href === '/'
				? $page.url.pathname === '/'
				: $page.url.pathname.startsWith(href);
		const base =
			'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors';
		if (active) return `${base} bg-primary/10 text-primary`;
		return `${base} text-muted-foreground hover:bg-accent`;
	}
</script>

<svelte:head>
	<title>SMS Dashboard</title>
</svelte:head>

{#if $page.url.pathname.startsWith('/login')}
	{@render children()}
{:else}
	<div class="flex min-h-screen">
		{#if sidebarOpen}
			<div
				class="fixed inset-0 z-40 bg-black/50 lg:hidden"
				role="presentation"
				onclick={() => (sidebarOpen = false)}
			></div>
		{/if}

		<aside
			class="fixed inset-y-0 left-0 z-50 w-60 -translate-x-full border-r bg-background transition-transform duration-200 lg:static lg:translate-x-0"
			class:translate-x-0={sidebarOpen}
		>
			<div class="flex h-14 items-center border-b px-4 font-bold text-lg">
				SMS Dashboard
			</div>

			<nav class="flex flex-col gap-1 p-2">
				{#each navItems as item}
					<a
						href={item.href}
						class={navClass(item.href)}
						onclick={() => (sidebarOpen = false)}
					>
						<span class="h-4 w-4">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="h-4 w-4"
							>
								{#if item.icon === 'LayoutDashboard'}
									<rect width="7" height="9" x="3" y="3" rx="1" />
									<rect width="7" height="5" x="14" y="3" rx="1" />
									<rect width="7" height="9" x="14" y="12" rx="1" />
									<rect width="7" height="5" x="3" y="16" rx="1" />
								{:else if item.icon === 'MessageSquare'}
									<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
								{:else if item.icon === 'Smartphone'}
									<rect width="14" height="20" x="5" y="2" rx="2" ry="2" />
									<path d="M12 18h.01" />
								{:else if item.icon === 'Link'}
									<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
									<path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
								{:else if item.icon === 'Settings'}
									<path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z" />
									<circle cx="12" cy="12" r="3" />
								{:else if item.icon === 'Key'}
									<circle cx="8" cy="21" r="1" />
									<circle cx="16" cy="21" r="1" />
									<path d="M10 17V7a2 2 0 0 1 2-2h4" />
									<path d="M6 17V4a2 2 0 0 1 2-2h8l4 4v11" />
								{/if}
							</svg>
						</span>
						{item.label}
					</a>
				{/each}
			</nav>

			<div class="absolute bottom-0 left-0 right-0 border-t p-2">
				<a
					href="/logout"
					class="flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-accent"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						class="h-4 w-4"
					>
						<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
						<polyline points="16 17 21 12 16 7" />
						<line x1="21" x2="9" y1="12" y2="12" />
					</svg>
					Logout
				</a>
			</div>
		</aside>

		<main class="flex-1 overflow-auto p-4 lg:p-6">
			<button
				class="mb-4 flex items-center gap-2 rounded-md p-2 text-muted-foreground hover:bg-accent lg:hidden"
				onclick={() => (sidebarOpen = true)}
				aria-label="Open sidebar"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="h-5 w-5"
				>
					<line x1="4" x2="20" y1="12" y2="12" />
					<line x1="4" x2="20" y1="6" y2="6" />
					<line x1="4" x2="20" y1="18" y2="18" />
				</svg>
				<span class="text-sm font-medium">Menu</span>
			</button>

			{@render children()}
		</main>
	</div>
{/if}
