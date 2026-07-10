<script lang="ts">
	import { onMount } from 'svelte';
	import { listWebhooks, createWebhook, deleteWebhook } from '$lib/api';
	import { loadDevices, dc } from '$lib/device-cache.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Table, Thead, Tbody, Tr, Th, Td } from '$lib/components/ui/table/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Alert, AlertDescription } from '$lib/components/ui/alert/index.js';
	import Dialog from '$lib/components/ui/dialog/dialog.svelte';
	import { toast } from 'svelte-sonner';
	import type { Webhook } from '$lib/types';

	type View = 'list' | 'create';
	let view = $state<View>('list');
	let webhooks = $state<Webhook[]>([]);
	let loading = $state(true);
	let error = $state('');

	let eventInput = $state('sms:received');
	let urlInput = $state('');
	let deviceIdInput = $state('');
	let saving = $state(false);
	let formError = $state('');

	const events = [
		{ value: 'sms:received', label: 'SMS Received' },
		{ value: 'sms:data-received', label: 'Data SMS Received' },
		{ value: 'sms:sent', label: 'SMS Sent' },
		{ value: 'sms:delivered', label: 'SMS Delivered' },
		{ value: 'sms:failed', label: 'SMS Failed' },
		{ value: 'system:ping', label: 'System Ping' },
		{ value: 'mms:received', label: 'MMS Received' },
		{ value: 'mms:downloaded', label: 'MMS Downloaded' },
	];

	onMount(() => { load(); loadDevices(); });

	async function load() {
		loading = true;
		error = '';
		try {
			webhooks = await listWebhooks();
		} catch {
			error = 'Failed to load webhooks';
		} finally {
			loading = false;
		}
	}

	function openCreate() {
		view = 'create';
		eventInput = 'sms:received';
		urlInput = '';
		deviceIdInput = '';
		formError = '';
	}

	function backToList() {
		view = 'list';
		load();
	}

	async function handleCreate() {
		formError = '';
		if (!urlInput.trim()) { formError = 'URL is required'; return; }
		saving = true;
		try {
			await createWebhook({
				url: urlInput.trim(),
				event: eventInput,
				...(deviceIdInput.trim() ? { deviceId: deviceIdInput.trim() } : {}),
			});
			backToList();
		} catch {
			formError = 'Failed to create webhook';
		} finally {
			saving = false;
		}
	}

	let deletingId = $state<string | null>(null);
	let confirmDelete = $state<string | null>(null);

	async function handleDelete() {
		if (!confirmDelete) return;
		const id = confirmDelete;
		deletingId = id;
		confirmDelete = null;
		try {
			await deleteWebhook(id);
			webhooks = webhooks.filter((w) => w.id !== id);
		} catch {
			toast.error('Failed to delete webhook');
		} finally {
			deletingId = null;
		}
	}

	function eventLabel(value: string): string {
		return events.find((e) => e.value === value)?.label ?? value;
	}

	function deviceName(id: string | undefined): string {
		if (!id) return '(all)';
		return dc.devices.find((d) => d.id === id)?.name ?? id;
	}
</script>

<div class="space-y-4">
	{#if view === 'create'}
		<div class="flex items-center gap-3">
			<Button variant="ghost" onclick={backToList}>&larr; Back</Button>
			<h1 class="text-2xl font-bold tracking-tight">Add Webhook</h1>
		</div>

		<Card class="mx-auto max-w-lg">
			<CardHeader>
				<CardTitle>New Webhook</CardTitle>
			</CardHeader>
			<CardContent>
				<form onsubmit={(e) => { e.preventDefault(); handleCreate(); }} class="space-y-4">
					{#if formError}
						<Alert variant="destructive">
							<AlertDescription>{formError}</AlertDescription>
						</Alert>
					{/if}

					<div class="space-y-2">
						<Label for="event">Event</Label>
						<select
							id="event"
							class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
							bind:value={eventInput}
							disabled={saving}
						>
							{#each events as ev}
								<option value={ev.value}>{ev.label}</option>
							{/each}
						</select>
					</div>

					<div class="space-y-2">
						<Label for="url">URL</Label>
						<Input id="url" type="url" bind:value={urlInput} placeholder="https://example.com/webhook" disabled={saving} />
					</div>

					<div class="space-y-2">
						<Label for="device">Device (optional)</Label>
						<select
							id="device"
							class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
							bind:value={deviceIdInput}
							disabled={saving}
						>
							<option value="">All devices</option>
							{#if dc.loading}
								<option value="" disabled>Loading...</option>
							{:else}
								{#each dc.devices as d}
									<option value={d.id}>{d.name}</option>
								{/each}
							{/if}
						</select>
					</div>

					<Button type="submit" class="w-full" disabled={saving}>
						{saving ? 'Creating...' : 'Create Webhook'}
					</Button>
				</form>
			</CardContent>
		</Card>
	{:else}
		<div class="flex items-center justify-between">
			<h1 class="text-2xl font-bold tracking-tight">Webhooks</h1>
			<Button onclick={openCreate}>+ Add</Button>
		</div>

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
		{:else if webhooks.length === 0}
			<div class="rounded-lg border p-12 text-center">
				<p class="text-lg font-medium">No webhooks configured</p>
				<p class="mt-1 text-sm text-muted-foreground">Create a webhook to receive real-time SMS events.</p>
				<Button class="mt-4" onclick={openCreate}>Add Webhook</Button>
			</div>
		{:else}
			<p class="text-sm text-muted-foreground">{webhooks.length} webhook{webhooks.length !== 1 ? 's' : ''}</p>

			<Table>
				<Thead>
					<Tr>
						<Th>Event</Th>
						<Th>URL</Th>
						<Th>Device</Th>
						<Th></Th>
					</Tr>
				</Thead>
				<Tbody>
					{#each webhooks as w}
						<Tr>
							<Td><Badge variant="secondary">{eventLabel(w.event)}</Badge></Td>
							<Td class="max-w-[400px] truncate">{w.url}</Td>
							<Td class="font-mono text-xs text-muted-foreground">{deviceName(w.deviceId)}</Td>
							<Td>
								<Button
									variant="destructive"
									size="sm"
									onclick={() => (confirmDelete = w.id)}
									disabled={deletingId === w.id}
								>
									{deletingId === w.id ? '...' : 'Delete'}
								</Button>
							</Td>
						</Tr>
					{/each}
				</Tbody>
			</Table>
		{/if}
	{/if}
</div>

<Dialog
	open={confirmDelete !== null}
	title="Delete Webhook"
	message="Are you sure you want to delete this webhook?"
	confirmLabel="Delete"
	variant="destructive"
	onConfirm={handleDelete}
	onCancel={() => (confirmDelete = null)}
/>
