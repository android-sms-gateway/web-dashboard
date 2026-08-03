<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { stats, trends } from '$lib/api';
	import { on as onEvent } from '$lib/events.svelte';
	import { toast } from 'svelte-sonner';
	import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import BarChart from '$lib/components/charts/BarChart.svelte';
	import LineChart from '$lib/components/charts/LineChart.svelte';
	import ActivityFeed from '$lib/components/feed/ActivityFeed.svelte';
	import type { Stats, TrendsResponse } from '$lib/types';

	type Range = 7 | 14 | 30;
	const RANGES: { value: Range; label: string }[] = [
		{ value: 7, label: '7 days' },
		{ value: 14, label: '14 days' },
		{ value: 30, label: '30 days' },
	];

	const STATS_REFRESH_DELAY = 1500;

	let data = $state<Stats | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let statsTimer: ReturnType<typeof setTimeout> | null = null;

	let range = $state<Range>(7);
	let trendsData = $state<TrendsResponse | null>(null);
	let trendsLoading = $state(true);
	let trendsError = $state<string | null>(null);
	let trendsRequestId = 0;

	const volumeData = $derived(
		trendsData?.messageVolume.map((d) => ({
			label: fmtDate(d.date),
			values: { sent: d.sent, failed: d.failed },
		})) ?? [],
	);
	const activityData = $derived(
		trendsData?.deviceActivity.map((d) => ({ label: fmtDate(d.date), values: { active: d.active } })) ?? [],
	);
	const volumeSeries = $derived([
		{ key: 'sent', label: 'Sent', color: 'hsl(var(--primary))' },
		{ key: 'failed', label: 'Failed', color: 'hsl(var(--destructive))' },
	]);
	const activitySeries = $derived([{ key: 'active', label: 'Active', color: 'hsl(var(--primary))' }]);
	const hasMessageData = $derived(volumeData.some((d) => d.values.sent > 0 || d.values.failed > 0));
	const hasActivityData = $derived(activityData.some((d) => d.values.active > 0));

	onMount(() => {
		loadStats();
		loadTrends();

		const unsubStats = onEvent('stats.updated', () => {
			if (statsTimer) clearTimeout(statsTimer);
			statsTimer = setTimeout(loadStats, STATS_REFRESH_DELAY);
		});

		const unsubMessage = onEvent('message.received', (e) => {
			toast.info(`SMS from ${e.payload.sender}: ${e.payload.message.slice(0, 50)}`);
		});

		const unsubState = onEvent('message.state_changed', (e) => {
			toast.info(`Message ${e.payload.messageId} → ${e.payload.state}`);
		});

		return () => {
			if (statsTimer) clearTimeout(statsTimer);
			unsubStats();
			unsubMessage();
			unsubState();
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

	async function loadTrends() {
		const requestId = ++trendsRequestId;
		trendsLoading = true;
		trendsError = null;
		try {
			const result = await trends(range);
			if (requestId !== trendsRequestId) return;
			trendsData = result;
		} catch (e: unknown) {
			if (requestId !== trendsRequestId) return;
			trendsError = e instanceof Error ? e.message : 'Failed to load trends';
		} finally {
			if (requestId === trendsRequestId) trendsLoading = false;
		}
	}

	function changeRange(value: Range) {
		if (value === range) return;
		range = value;
		loadTrends();
	}

	function fmtDate(date: string): string {
		const [y, m, d] = date.split('-').map(Number);
		return new Date(Date.UTC(y, m - 1, d)).toLocaleDateString(undefined, {
			month: 'short',
			day: 'numeric',
			timeZone: 'UTC',
		});
	}
</script>

<div class="mx-auto max-w-6xl space-y-6">
	<div>
		<h1 class="text-3xl font-bold tracking-tight">Dashboard</h1>
		<p class="text-muted-foreground">
			Welcome back, {$page.data.user?.login ?? 'User'}
		</p>
	</div>

	<div class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
		<div class="min-w-0 space-y-6">
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

			<div class="flex flex-wrap items-center justify-between gap-2">
				<h2 class="text-xl font-semibold tracking-tight">Trends</h2>
				<div class="flex gap-1" role="group" aria-label="Trend range">
					{#each RANGES as r}
						<Button
							size="sm"
							variant={range === r.value ? 'default' : 'outline'}
							aria-pressed={range === r.value}
							onclick={() => changeRange(r.value)}
						>
							{r.label}
						</Button>
					{/each}
				</div>
			</div>

			{#if trendsLoading}
				<div class="grid gap-4">
					<Card>
						<CardContent class="p-6">
							<Skeleton class="h-56 w-full rounded-lg" />
						</CardContent>
					</Card>
					<Card>
						<CardContent class="p-6">
							<Skeleton class="h-56 w-full rounded-lg" />
						</CardContent>
					</Card>
				</div>
			{:else if trendsError}
				<Card>
					<CardContent class="p-6 text-center">
						<p class="text-destructive">{trendsError}</p>
						<Button variant="link" class="mt-2 text-sm" onclick={loadTrends}>
							Try again
						</Button>
					</CardContent>
				</Card>
			{:else if trendsData}
				<div class="grid gap-4">
					<Card>
						<CardHeader>
							<CardTitle>Message volume</CardTitle>
							<CardDescription>Sent vs failed per day</CardDescription>
						</CardHeader>
						<CardContent class="p-6 pt-0">
							{#if hasMessageData}
								<BarChart
									data={volumeData}
									series={volumeSeries}
									ariaLabel="Message volume per day"
								/>
							{:else}
								<div class="flex h-56 items-center justify-center rounded-lg border border-dashed">
									<p class="text-sm text-muted-foreground">No messages yet</p>
								</div>
							{/if}
						</CardContent>
					</Card>
					<Card>
						<CardHeader>
							<CardTitle>Device activity</CardTitle>
							<CardDescription>Devices last active per day</CardDescription>
						</CardHeader>
						<CardContent class="p-6 pt-0">
							{#if hasActivityData}
								<LineChart
									data={activityData}
									series={activitySeries}
									ariaLabel="Devices last active per day"
								/>
							{:else}
								<div class="flex h-56 items-center justify-center rounded-lg border border-dashed">
									<p class="text-sm text-muted-foreground">No device activity yet</p>
								</div>
							{/if}
						</CardContent>
					</Card>
				</div>
			{/if}
		</div>

		<div class="min-w-0">
			<Card>
				<CardHeader>
					<CardTitle>Activity</CardTitle>
					<CardDescription>Live message events</CardDescription>
				</CardHeader>
				<CardContent class="p-0">
					<ActivityFeed />
				</CardContent>
			</Card>
		</div>
	</div>
</div>
