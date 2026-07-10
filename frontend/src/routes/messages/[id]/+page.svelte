<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { getMessage } from '$lib/api';
	import { loadDevices } from '$lib/device-cache.svelte';
	import { dc } from '$lib/device-cache.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Table, Thead, Tbody, Tr, Th, Td } from '$lib/components/ui/table/index.js';
	import type { MessageDetail } from '$lib/types';

	let message = $state<MessageDetail | null>(null);
	let loading = $state(true);
	let error = $state('');

	onMount(() => {
		loadDetail();
		loadDevices();
	});

	async function loadDetail() {
		const id = $page.params.id;
		if (!id) return;
		loading = true;
		error = '';
		try {
			message = await getMessage(id);
		} catch {
			error = 'Failed to load message details';
		} finally {
			loading = false;
		}
	}

	function deviceName(id: string) {
		return dc.devices.find((d) => d.id === id)?.name ?? id;
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
</script>

<div class="mx-auto max-w-3xl space-y-6">
	<div class="flex items-center gap-3">
		<Button variant="ghost" onclick={() => goto('/messages')}>&larr; Back to list</Button>
		<h1 class="text-2xl font-bold tracking-tight">Message Detail</h1>
	</div>

	{#if loading}
		<div class="space-y-3">
			<Skeleton class="h-40 w-full rounded-lg" />
			<Skeleton class="h-32 w-full rounded-lg" />
		</div>
	{:else if error}
		<div class="rounded-lg border p-6 text-center">
			<p class="text-destructive">{error}</p>
			<Button variant="outline" class="mt-2" onclick={loadDetail}>Retry</Button>
		</div>
	{:else if message}
		<Card>
			<CardHeader>
				<CardTitle>Message</CardTitle>
			</CardHeader>
			<CardContent class="space-y-3">
				<div class="grid grid-cols-[100px_1fr] gap-2 text-sm">
					<span class="font-medium text-muted-foreground">ID</span>
					<span class="font-mono text-xs">{message.id}</span>

					<span class="font-medium text-muted-foreground">Device</span>
					<span>{deviceName(message.deviceId)}</span>

					<span class="font-medium text-muted-foreground">State</span>
					<span>
						<Badge variant={stateBadgeVariant(message.state)}>{message.state}</Badge>
					</span>

					{#if message.textMessage}
						<span class="font-medium text-muted-foreground">Message</span>
						<span class="whitespace-pre-wrap break-all">{message.textMessage.text}</span>
					{/if}

					{#if message.dataMessage}
						<span class="font-medium text-muted-foreground">Data</span>
						<span class="font-mono text-xs break-all">
							port={message.dataMessage.port} data={message.dataMessage.data}
						</span>
					{/if}

					{#if message.hashedMessage}
						<span class="font-medium text-muted-foreground">Hash</span>
						<span class="font-mono text-xs break-all">{message.hashedMessage.hash}</span>
					{/if}
				</div>
			</CardContent>
		</Card>

		<div>
			<h2 class="mb-3 text-lg font-semibold">Recipients</h2>
			<Table>
				<Thead>
					<Tr>
						<Th>Phone Number</Th>
						<Th>State</Th>
						<Th>Error</Th>
					</Tr>
				</Thead>
				<Tbody>
					{#each message.recipients as r}
						<Tr>
							<Td class="font-mono text-xs">{r.phoneNumber}</Td>
							<Td>
								<Badge variant={stateBadgeVariant(r.state)}>{r.state}</Badge>
							</Td>
							<Td class="text-destructive text-xs">{r.error ?? ''}</Td>
						</Tr>
					{/each}
				</Tbody>
			</Table>
		</div>

		<div>
			<h2 class="mb-3 text-lg font-semibold">Timeline</h2>
			<Table>
				<Thead>
					<Tr>
						<Th>Event</Th>
						<Th>Time</Th>
					</Tr>
				</Thead>
				<Tbody>
					{#each Object.entries(message.states) as [event, time]}
						<Tr>
							<Td class="font-medium">{event}</Td>
							<Td class="text-muted-foreground">{time}</Td>
						</Tr>
					{/each}
				</Tbody>
			</Table>
		</div>
	{/if}
</div>
