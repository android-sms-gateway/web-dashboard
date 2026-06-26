<script lang="ts">
	import Sidebar from "../components/Sidebar.svelte";
	import DashboardPage from "./DashboardPage.svelte";
	import MessagesPage from "./MessagesPage.svelte";
	import DevicesPage from "./DevicesPage.svelte";
	import WebhooksPage from "./WebhooksPage.svelte";
	import SettingsPage from "./SettingsPage.svelte";
	import TokensPage from "./TokensPage.svelte";

	let { onLogout }: { onLogout: () => void } = $props();

	let currentPage = $state("dashboard");
	let sidebarOpen = $state(false);

	function onNavigate(page: string) {
		currentPage = page;
	}

	function toggleSidebar() {
		sidebarOpen = !sidebarOpen;
	}

	function closeSidebar() {
		sidebarOpen = false;
	}
</script>

<div class="layout">
	<Sidebar
		current={currentPage}
		{onNavigate}
		{onLogout}
		isOpen={sidebarOpen}
		onToggle={toggleSidebar}
		onClose={closeSidebar}
	/>
	<main class="content">
		<button class="hamburger" onclick={toggleSidebar} aria-label="Open sidebar">
			<span class="hamburger-line"></span>
			<span class="hamburger-line"></span>
			<span class="hamburger-line"></span>
		</button>

		{#if currentPage === "dashboard"}
			<DashboardPage />
		{:else if currentPage === "messages"}
			<MessagesPage />
		{:else if currentPage === "devices"}
			<DevicesPage />
		{:else if currentPage === "webhooks"}
			<WebhooksPage />
		{:else if currentPage === "settings"}
			<SettingsPage />
		{:else if currentPage === "tokens"}
			<TokensPage />
		{:else}
			<h2>{currentPage}</h2>
			<p>This page is coming soon.</p>
		{/if}
	</main>
</div>

<style>
	.layout {
		display: flex;
		min-height: 100vh;
		font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
			sans-serif;
	}

	.content {
		flex: 1;
		padding: 2rem;
		background: #f8fafc;
	}

	.hamburger {
		display: none;
		flex-direction: column;
		gap: 4px;
		background: none;
		border: none;
		cursor: pointer;
		padding: 0.5rem;
		margin-bottom: 1rem;
	}

	.hamburger-line {
		display: block;
		width: 20px;
		height: 2px;
		background: #6b7280;
		border-radius: 1px;
	}

	@media (max-width: 768px) {
		.content {
			padding: 1rem;
		}

		.hamburger {
			display: flex;
		}
	}
</style>
