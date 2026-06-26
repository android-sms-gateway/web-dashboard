<script lang="ts">
	import { onMount } from "svelte";
	import Login from "./pages/Login.svelte";
	import Dashboard from "./pages/Dashboard.svelte";
	import Toast from "./components/Toast.svelte";
	import { me, logout } from "./lib/api";
	import {
		connect as connectEvents,
		disconnect as disconnectEvents,
	} from "./lib/events.svelte";

	let currentPage = $state("loading");

	onMount(async () => {
		try {
			await me();
			currentPage = "dashboard";
			connectEvents();
		} catch {
			currentPage = "login";
		}
	});

	function handleLogin() {
		currentPage = "dashboard";
		connectEvents();
	}

	async function handleLogout() {
		disconnectEvents();
		try {
			await logout();
		} catch {
			// Server session may already be gone
		}
		currentPage = "login";
	}
</script>

{#if currentPage === "loading"}
	<div class="loading-screen">
		<div class="spinner"></div>
		<span>Loading...</span>
	</div>
{:else if currentPage === "login"}
	<Login onLogin={handleLogin} />
{:else}
	<Dashboard onLogout={handleLogout} />
{/if}

<Toast />

<style>
	.loading-screen {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		justify-content: center;
		align-items: center;
		min-height: 100vh;
		font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
			sans-serif;
		color: #6b7280;
		font-size: 0.9rem;
	}

	.spinner {
		width: 28px;
		height: 28px;
		border: 3px solid #e5e7eb;
		border-top-color: #2563eb;
		border-radius: 50%;
		animation: spin 0.7s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
