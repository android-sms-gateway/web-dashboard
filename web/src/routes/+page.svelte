<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { stats } from '$lib/api';
	import { on as onEvent } from '$lib/events.svelte';
	import { toast } from 'svelte-sonner';
	import { Card, CardHeader, CardTitle, CardContent } from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import type { Stats } from '$lib/types';

	let data = $state<Stats | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);

	onMount(() => {
		loadStats();

		const unsubStats = onEvent('stats.updated', () => {
			loadStats();
		});

		const unsubMessage = onEvent('message.received', (e) => {
			toast.info(`SMS from ${e.payload.sender}: ${e.payload.message.slice(0, 50)}`);
		});

		const unsubState = onEvent('message.state_changed', (e) => {
			toast.info(`Message ${e.payload.messageId} → ${e.payload.state}`);
		});

		const unsubDevice = onEvent('device.status_changed', (e) => {
			if (e.payload.isOnline) {
				toast.success(`Device ${e.payload.deviceId} online`);
			} else {
				toast.error(`Device ${e.payload.deviceId} offline`);
			}
		});

		return () => {
			unsubStats();
			unsubMessage();
			unsubState();
			unsubDevice();
		};
	});

	async function loadStats() {
		try {
			data = await stats();
			error = null;
		} catch (e: unknown) {
			error = e instanceof Error ? e.message : 'Failed to load stats';
		} finally {
			loading = false;
		}
	}
</script>

<div class="mx-auto max-w-4xl space-y-6">
	<div>
		<h1 class="text-3xl font-bold tracking-tight">Dashboard</h1>
		<p class="text-muted-foreground">
			Welcome back, {$page.data.user?.login ?? 'User'}
		</p>
	</div>

	<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
		{#if loading}
			<Skeleton class="h-32 rounded-lg" />
			<Skeleton class="h-32 rounded-lg" />
			<Skeleton class="h-32 rounded-lg" />
		{:else if error}
			<Card class="md:col-span-2 lg:col-span-3">
				<CardContent class="p-6 text-center">
					<p class="text-destructive">{error}</p>
					<button
						class="mt-2 text-sm text-primary hover:underline"
						onclick={() => { loading = true; loadStats(); }}
					>
						Try again
					</button>
				</CardContent>
			</Card>
		{:else if data}
			<Card>
				<CardHeader class="pb-2">
					<CardTitle class="text-sm font-medium text-muted-foreground">
						Devices Online
					</CardTitle>
				</CardHeader>
				<CardContent>
					<p class="text-2xl font-bold">{data.devicesOnline}</p>
					<p class="text-xs text-muted-foreground">
						{data.devicesActive} active · {data.devicesTotal} total
					</p>
				</CardContent>
			</Card>
			<Card>
				<CardHeader class="pb-2">
					<CardTitle class="text-sm font-medium text-muted-foreground">
						Messages Sent
					</CardTitle>
				</CardHeader>
				<CardContent>
					<p class="text-2xl font-bold">{data.messagesSent}</p>
					<p class="text-xs text-muted-foreground">Total all time</p>
				</CardContent>
			</Card>
			<Card>
				<CardHeader class="pb-2">
					<CardTitle class="text-sm font-medium text-muted-foreground">
						Pending / Failed
					</CardTitle>
				</CardHeader>
				<CardContent>
					<div class="flex items-center gap-2">
						<p class="text-2xl font-bold">{data.messagesPending}</p>
						<Badge variant="secondary">pending</Badge>
					</div>
					<div class="flex items-center gap-2 mt-1">
						<p class="text-2xl font-bold">{data.messagesFailed}</p>
						<Badge variant="destructive">failed</Badge>
					</div>
				</CardContent>
			</Card>
		{/if}
	</div>
</div>
