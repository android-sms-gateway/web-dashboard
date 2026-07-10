<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { listMessages } from '$lib/api';
	import { loadDevices } from '$lib/device-cache.svelte';
	import { dc } from '$lib/device-cache.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Table, Thead, Tbody, Tr, Th, Td } from '$lib/components/ui/table/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import type { MessageListItem } from '$lib/types';

	let messages = $state<MessageListItem[]>([]);
	let total = $state(0);
	let loading = $state(true);
	let error = $state('');

	let offset = $state(0);
	let limit = $state(50);
	let loadSeq = 0;

	let stateFilter = $state('');
	let deviceFilter = $state('');
	let fromFilter = $state('');
	let toFilter = $state('');

	const stateOptions = [
		{ value: '', label: 'All states' },
		{ value: 'Pending', label: 'Pending' },
		{ value: 'Processed', label: 'Processed' },
		{ value: 'Sent', label: 'Sent' },
		{ value: 'Delivered', label: 'Delivered' },
		{ value: 'Failed', label: 'Failed' },
	];

	const pageSizes = [25, 50, 100, 200];

	let page = $derived(Math.floor(offset / limit) + 1);
	let totalPages = $derived(Math.max(1, Math.ceil(total / limit)));

	onMount(() => {
		loadMessages();
		loadDevices();
	});

	async function loadMessages() {
		const seq = ++loadSeq;
		loading = true;
		error = '';
		try {
			const res = await listMessages({
				limit,
				offset,
				state: stateFilter || undefined,
				deviceId: deviceFilter || undefined,
				from: fromFilter ? new Date(`${fromFilter}T00:00:00`).toISOString() : undefined,
				to: toFilter ? new Date(`${toFilter}T23:59:59`).toISOString() : undefined,
			});
			if (seq !== loadSeq) return;
			messages = res.items;
			total = res.total;
		} catch {
			if (seq !== loadSeq) return;
			error = 'Failed to load messages';
		} finally {
			if (seq === loadSeq) loading = false;
		}
	}

	function applyFilters() {
		offset = 0;
		loadMessages();
	}

	function goPage(delta: number) {
		const newOffset = Math.max(0, Math.min(offset + delta * limit, (totalPages - 1) * limit));
		if (newOffset !== offset) {
			offset = newOffset;
			loadMessages();
		}
	}

	function stateBadgeVariant(state: string) {
		switch (state) {
			case 'Delivered': return 'default';
			case 'Failed': return 'destructive';
			case 'Sent': return 'secondary';
			case 'Processed': return 'secondary';
			default: return 'outline';
		}
	}

	function hasFilters() {
		return stateFilter || deviceFilter || fromFilter || toFilter;
	}

	function resetFilters() {
		stateFilter = '';
		deviceFilter = '';
		fromFilter = '';
		toFilter = '';
		applyFilters();
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h1 class="text-2xl font-bold tracking-tight">Messages</h1>
		<Button onclick={() => goto('/messages/compose')}>+ Send</Button>
	</div>

	<div class="flex flex-wrap items-end gap-3 rounded-lg border p-3">
		<div class="space-y-1">
			<Label for="filter-state">State</Label>
			<select
				id="filter-state"
				class="flex h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
				value={stateFilter}
				onchange={(e) => (stateFilter = (e.target as HTMLSelectElement).value)}
			>
				{#each stateOptions as opt}
					<option value={opt.value}>{opt.label}</option>
				{/each}
			</select>
		</div>

		<div class="space-y-1">
			<Label for="filter-device">Device</Label>
			<select
				id="filter-device"
				class="flex h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
				value={deviceFilter}
				onchange={(e) => (deviceFilter = (e.target as HTMLSelectElement).value)}
			>
				<option value="">All devices</option>
				{#if dc.loading}
					<option value="" disabled>Loading...</option>
				{:else}
					{#each dc.devices as device}
						<option value={device.id}>{device.name}</option>
					{/each}
				{/if}
				{#if deviceFilter && !dc.devices.some((d) => d.id === deviceFilter)}
					<option value={deviceFilter} disabled>(previous selection)</option>
				{/if}
			</select>
		</div>

		<div class="space-y-1">
			<Label for="filter-from">From</Label>
			<Input id="filter-from" type="date" bind:value={fromFilter} />
		</div>

		<div class="space-y-1">
			<Label for="filter-to">To</Label>
			<Input id="filter-to" type="date" bind:value={toFilter} />
		</div>

		<div class="flex gap-2">
			<Button variant="secondary" onclick={applyFilters}>Apply</Button>
			{#if hasFilters()}
				<Button variant="ghost" onclick={resetFilters}>Clear</Button>
			{/if}
		</div>
	</div>

	{#if loading}
		<div class="space-y-2">
			<Skeleton class="h-10 w-full rounded-lg" />
			<Skeleton class="h-10 w-full rounded-lg" />
			<Skeleton class="h-10 w-full rounded-lg" />
			<Skeleton class="h-10 w-full rounded-lg" />
			<Skeleton class="h-10 w-full rounded-lg" />
		</div>
	{:else if error}
		<div class="rounded-lg border p-6 text-center">
			<p class="text-destructive">{error}</p>
			<Button variant="outline" class="mt-2" onclick={loadMessages}>Retry</Button>
		</div>
	{:else if messages.length === 0}
		<div class="rounded-lg border p-12 text-center">
			<p class="text-lg font-medium">No messages yet</p>
			<p class="mt-1 text-sm text-muted-foreground">Send your first message to get started.</p>
			<Button class="mt-4" onclick={() => goto('/messages/compose')}>Send Message</Button>
		</div>
	{:else}
		<p class="text-sm text-muted-foreground">{total} message{total !== 1 ? 's' : ''}</p>

		<Table>
			<Thead>
				<Tr>
					<Th>Preview</Th>
					<Th>Recipients</Th>
					<Th>State</Th>
					<Th>Time</Th>
					<Th></Th>
				</Tr>
			</Thead>
			<Tbody>
				{#each messages as msg}
					<Tr class="cursor-pointer" onclick={() => goto(`/messages/${msg.id}`)}>
						<Td class="max-w-[280px] truncate font-medium">
							{msg.textPreview || '(binary)'}
						</Td>
						<Td>{msg.recipients.length}</Td>
						<Td>
							<Badge variant={stateBadgeVariant(msg.state)}>
								{msg.state}
							</Badge>
						</Td>
						<Td class="whitespace-nowrap text-muted-foreground">
							{msg.createdAt ? new Date(msg.createdAt).toLocaleString() : ''}
						</Td>
						<Td>
							<Button variant="outline" size="sm" onclick={(e) => { e.stopPropagation(); goto(`/messages/${msg.id}`); }}>
								View
							</Button>
						</Td>
					</Tr>
				{/each}
			</Tbody>
		</Table>

		<div class="flex items-center gap-3 pt-3 text-sm">
			<Button variant="outline" size="sm" disabled={page <= 1} onclick={() => goPage(-1)}>
				&larr; Prev
			</Button>
			<span class="text-muted-foreground">
				Page {page} of {totalPages}
				<span class="ml-1 text-muted-foreground/60">({total} total)</span>
			</span>
			<Button variant="outline" size="sm" disabled={page >= totalPages} onclick={() => goPage(1)}>
				Next &rarr;
			</Button>
			<span class="ml-auto text-muted-foreground">Show:</span>
			<select
				class="flex h-8 rounded-md border border-input bg-background px-2 text-sm"
				value={limit}
				onchange={(e) => {
					const v = parseInt((e.target as HTMLSelectElement).value, 10);
					if (v !== limit) { limit = v; offset = 0; loadMessages(); }
				}}
			>
				{#each pageSizes as n}
					<option value={n}>{n}</option>
				{/each}
			</select>
		</div>
	{/if}
</div>
