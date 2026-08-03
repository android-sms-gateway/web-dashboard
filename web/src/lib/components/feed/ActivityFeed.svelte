<script lang="ts">
	import { onMount } from 'svelte';
	import { listMessages } from '$lib/api';
	import { on as onEvent } from '$lib/events.svelte';
	import { stateBadgeVariant } from '$lib/utils';
	import type { MessageListItem, MessageReceivedPayload, MessageStateChangedPayload } from '$lib/types';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';

	const MAX_ITEMS = 50;
	const SEED_LIMIT = 20;

	type FeedItem = {
		id: string;
		kind: 'received' | 'state';
		title: string;
		preview?: string;
		state?: string;
		ts: number;
	};

	let items = $state<FeedItem[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let now = $state(Date.now());

	let seq = 0;
	let timer: ReturnType<typeof setInterval> | null = null;

	onMount(() => {
		const unsubReceived = onEvent('message.received', (e) => onReceived(e.payload));
		const unsubState = onEvent('message.state_changed', (e) => onStateChanged(e.payload));

		seed();

		timer = setInterval(() => {
			now = Date.now();
		}, 30_000);

		return () => {
			unsubReceived();
			unsubState();
			if (timer) clearInterval(timer);
		};
	});

	function onReceived(payload: MessageReceivedPayload) {
		seq += 1;
		prepend({
			id: `received-${seq}`,
			kind: 'received',
			title: `From ${payload.sender}`,
			preview: payload.message,
			ts: Date.now(),
		});
	}

	function onStateChanged(payload: MessageStateChangedPayload) {
		const idx = items.findIndex((i) => i.id === payload.messageId);
		if (idx >= 0) {
			items[idx].state = payload.state;
			items[idx].ts = Date.now();
			items = [items[idx], ...items.filter((_, j) => j !== idx)].slice(0, MAX_ITEMS);
			return;
		}

		prepend({
			id: payload.messageId,
			kind: 'state',
			title: `Message ${payload.messageId}`,
			state: payload.state,
			ts: Date.now(),
		});
	}

	function prepend(item: FeedItem) {
		items = [item, ...items].slice(0, MAX_ITEMS);
	}

	function seedItems(list: MessageListItem[]) {
		const seeded: FeedItem[] = list.map((m) => {
			const ts = Date.parse(m.createdAt);
			return {
				id: m.id,
				kind: 'state',
				title: m.recipients.map((r) => r.phoneNumber).join(', ') || 'Unknown recipient',
				preview: m.textPreview,
				state: m.state,
				ts: Number.isNaN(ts) ? Date.now() : ts,
			};
		});
		seeded.sort((a, b) => b.ts - a.ts);

		const live = new Set(items.map((i) => i.id));
		items = [...items, ...seeded.filter((i) => !live.has(i.id))]
			.sort((a, b) => b.ts - a.ts)
			.slice(0, MAX_ITEMS);
	}

	async function seed() {
		loading = true;
		error = null;
		try {
			const res = await listMessages({ limit: SEED_LIMIT });
			seedItems(res.items);
		} catch (e: unknown) {
			error = e instanceof Error ? e.message : 'Failed to load activity';
		} finally {
			loading = false;
		}
	}

	function timeAgo(ts: number): string {
		const diff = Math.max(0, now - ts);
		const s = Math.floor(diff / 1000);
		if (s < 60) return 'just now';
		const m = Math.floor(s / 60);
		if (m < 60) return `${m}m ago`;
		const h = Math.floor(m / 60);
		if (h < 24) return `${h}h ago`;
		const d = Math.floor(h / 24);
		return `${d}d ago`;
	}
</script>

{#if loading}
	<div class="space-y-3 px-4 py-4">
		{#each Array(4) as _}
			<Skeleton class="h-14 w-full rounded-lg" />
		{/each}
	</div>
{:else if items.length === 0 && error}
	<div class="p-6 text-center">
		<p class="text-destructive">{error}</p>
		<button class="mt-2 text-sm text-primary hover:underline" onclick={seed}>Try again</button>
	</div>
{:else if items.length === 0}
	<div class="flex h-56 items-center justify-center rounded-lg border border-dashed p-6">
		<p class="text-sm text-muted-foreground">No activity yet</p>
	</div>
{:else}
	<ul aria-label="Recent activity" class="divide-y divide-border lg:max-h-[calc(100vh-16rem)] lg:overflow-y-auto">
		{#each items as item (item.id)}
			<li class="flex items-start gap-3 px-4 py-3">
				<span
					class={`mt-1.5 h-2 w-2 shrink-0 rounded-full ${item.kind === 'received' ? 'bg-primary' : 'bg-muted-foreground'}`}
				></span>
				<div class="min-w-0 flex-1">
					<p class="truncate text-sm font-medium text-foreground">{item.title}</p>
					{#if item.preview}
						<p class="truncate text-xs text-muted-foreground">{item.preview}</p>
					{/if}
				</div>
				<div class="flex shrink-0 flex-col items-end gap-1">
					{#if item.state}
						<Badge variant={stateBadgeVariant(item.state)}>{item.state}</Badge>
					{/if}
					<time class="text-xs text-muted-foreground" datetime={new Date(item.ts).toISOString()}>
						{timeAgo(item.ts)}
					</time>
				</div>
			</li>
		{/each}
	</ul>
{/if}
