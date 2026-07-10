<script lang="ts">
	import { onMount } from 'svelte';
	import { listDevices, deleteDevice } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Table, Thead, Tbody, Tr, Th, Td } from '$lib/components/ui/table/index.js';
	import Dialog from '$lib/components/ui/dialog/dialog.svelte';
	import type { Device } from '$lib/types';

	let devices = $state<Device[]>([]);
	let loading = $state(true);
	let error = $state('');

	onMount(load);

	async function load() {
		loading = true;
		error = '';
		try {
			devices = await listDevices();
		} catch {
			error = 'Failed to load devices';
		} finally {
			loading = false;
		}
	}

	let deleting = $state<string | null>(null);
	let confirmDelete = $state<string | null>(null);

	async function handleDelete() {
		if (!confirmDelete) return;
		const id = confirmDelete;
		deleting = id;
		confirmDelete = null;
		try {
			await deleteDevice(id);
			devices = devices.filter((d) => d.id !== id);
		} catch {
			// ignore
		} finally {
			deleting = null;
		}
	}

	function timeAgo(iso: string): string {
		const ts = new Date(iso).getTime();
		if (!Number.isFinite(ts)) return 'unknown';
		const diff = Date.now() - ts;
		const minutes = Math.floor(diff / 60000);
		if (minutes < 1) return 'just now';
		if (minutes < 60) return `${minutes}m ago`;
		const hours = Math.floor(minutes / 60);
		if (hours < 24) return `${hours}h ago`;
		const days = Math.floor(hours / 24);
		return `${days}d ago`;
	}
</script>

<div class="space-y-4">
	<h1 class="text-2xl font-bold tracking-tight">Devices</h1>

	{#if loading}
		<div class="space-y-2">
			<Skeleton class="h-10 w-full rounded-lg" />
			<Skeleton class="h-10 w-full rounded-lg" />
			<Skeleton class="h-10 w-full rounded-lg" />
		</div>
	{:else if error}
		<div class="rounded-lg border p-6 text-center">
			<p class="text-destructive">{error}</p>
			<Button variant="outline" class="mt-2" onclick={load}>Retry</Button>
		</div>
	{:else if devices.length === 0}
		<div class="rounded-lg border p-12 text-center">
			<p class="text-lg font-medium">No devices registered</p>
			<p class="mt-1 text-sm text-muted-foreground">Register a device on the SMS Gateway app to get started.</p>
		</div>
	{:else}
		<p class="text-sm text-muted-foreground">{devices.length} device{devices.length !== 1 ? 's' : ''}</p>

		<Table>
			<Thead>
				<Tr>
					<Th>Name</Th>
					<Th>ID</Th>
					<Th>Status</Th>
					<Th>Last Seen</Th>
					<Th></Th>
				</Tr>
			</Thead>
			<Tbody>
				{#each devices as d (d.id)}
					<Tr>
						<Td class="font-medium">{d.name || '(unnamed)'}</Td>
						<Td class="font-mono text-xs text-muted-foreground">{d.id}</Td>
						<Td>
							{#if d.isOnline}
								<Badge variant="default">Online</Badge>
							{:else}
								<Badge variant="outline">Offline</Badge>
							{/if}
						</Td>
						<Td class="whitespace-nowrap text-muted-foreground">{timeAgo(d.lastSeen)}</Td>
						<Td>
							<Button
								variant="destructive"
								size="sm"
								onclick={() => (confirmDelete = d.id)}
								disabled={deleting === d.id}
							>
								{deleting === d.id ? '...' : 'Delete'}
							</Button>
						</Td>
					</Tr>
				{/each}
			</Tbody>
		</Table>
	{/if}
</div>

<Dialog
	open={confirmDelete !== null}
	title="Delete Device"
	message="Are you sure you want to delete this device? This action cannot be undone."
	confirmLabel="Delete"
	variant="destructive"
	onConfirm={handleDelete}
	onCancel={() => (confirmDelete = null)}
/>
