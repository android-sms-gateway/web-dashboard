<script lang="ts">
	import { onMount } from "svelte";
	import LoadingSkeleton from "../components/LoadingSkeleton.svelte";
	import ErrorBanner from "../components/ErrorBanner.svelte";

	let stats: {
		messagesSent: number;
		messagesPending: number;
		messagesFailed: number;
		devicesOnline: number;
		devicesActive: number;
		devicesTotal: number;
	} | null = $state(null);
	let userName = $state("");
	let loading = $state(true);
	let error = $state("");

	onMount(load);

	async function load() {
		loading = true;
		error = "";
		try {
			const [meRes, statsRes] = await Promise.all([
				fetch("/api/v1/auth/me"),
				fetch("/api/v1/stats"),
			]);

			if (meRes.ok) {
				const me = await meRes.json();
				userName = me.login;
			}

			if (statsRes.ok) {
				stats = await statsRes.json();
			} else {
				throw new Error("Failed to load stats");
			}
		} catch {
			error = "Failed to load dashboard data";
		} finally {
			loading = false;
		}
	}
</script>

<div class="page">
	{#if loading}
		<LoadingSkeleton variant="card" rows={5} />
	{:else if error}
		<ErrorBanner message={error} onRetry={load} onDismiss={() => error = ""} />
		<div class="cards">
			{#each ["Devices Online", "Devices Active", "Total Devices", "Messages Sent", "Pending"] as label}
				<div class="card">
					<div class="card-value">?</div>
					<div class="card-label">{label}</div>
				</div>
			{/each}
		</div>
	{:else}
		<div class="welcome">
			Welcome, <strong>{userName || "User"}</strong>
		</div>

		<div class="cards">
			<div class="card">
				<div class="card-value">{stats?.devicesOnline ?? "?"}</div>
				<div class="card-label">Devices Online</div>
			</div>
			<div class="card">
				<div class="card-value">{stats?.devicesActive ?? "?"}</div>
				<div class="card-label">Devices Active</div>
			</div>
			<div class="card">
				<div class="card-value">{stats?.devicesTotal ?? "?"}</div>
				<div class="card-label">Total Devices</div>
			</div>
			<div class="card">
				<div class="card-value">{stats?.messagesSent ?? "?"}</div>
				<div class="card-label">Messages Sent</div>
			</div>
			<div class="card">
				<div class="card-value">{stats?.messagesPending ?? "?"}</div>
				<div class="card-label">Pending</div>
			</div>
		</div>
	{/if}
</div>

<style>
	.welcome {
		margin-bottom: 1.5rem;
		font-size: 1.1rem;
		color: #374151;
	}

	.cards {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
		gap: 1rem;
	}

	.card {
		background: white;
		padding: 1.5rem;
		border-radius: 8px;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
		text-align: center;
		transition: transform 0.15s ease, box-shadow 0.15s ease;
	}

	.card:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
	}

	.card-value {
		font-size: 2rem;
		font-weight: 700;
		color: #2563eb;
	}

	.card-label {
		margin-top: 0.25rem;
		font-size: 0.875rem;
		color: #6b7280;
	}
</style>
