<script lang="ts">
	import { cn } from "$lib/utils";
	import { formatTick, niceMax, tickValues, xLabelIndexes } from "./scale";

	const W = 720;
	const H = 240;
	const PAD_L = 40;
	const PAD_R = 8;
	const PAD_T = 10;
	const PAD_B = 26;

	interface Props {
		data: { label: string; values: Record<string, number> }[];
		series: { key: string; label: string; color: string }[];
		ariaLabel?: string;
		class?: string;
	}

	let {
		data,
		series,
		ariaLabel = "Line chart",
		class: className,
	}: Props = $props();

	let hovered = $state<number | null>(null);

	$effect(() => {
		data;
		hovered = null;
	});

	const plotW = W - PAD_L - PAD_R;
	const plotH = H - PAD_T - PAD_B;
	const max = $derived(
		niceMax(data.flatMap((d) => series.map((s) => d.values[s.key] ?? 0))),
	);
	const labels = $derived(xLabelIndexes(data.length));

	function y(v: number): number {
		return PAD_T + plotH * (1 - v / max);
	}

	function x(i: number): number {
		if (data.length <= 1) return PAD_L + plotW / 2;
		return PAD_L + (plotW * i) / (data.length - 1);
	}

	function linePath(s: { key: string }): string {
		return data
			.map(
				(d, i) =>
					`${i === 0 ? "M" : "L"} ${x(i)} ${y(d.values[s.key] ?? 0)}`,
			)
			.join(" ");
	}

	function areaPath(s: { key: string }): string {
		const baseline = y(0);
		return `${linePath(s)} L ${x(data.length - 1)} ${baseline} L ${x(0)} ${baseline} Z`;
	}

	function tooltipStyle(i: number): string {
		const left = `${(x(i) / W) * 100}%`;
		if (i <= 1) return `left: ${left}`;
		if (i >= data.length - 2)
			return `left: ${left}; transform: translateX(-100%)`;
		return `left: ${left}; transform: translateX(-50%)`;
	}
</script>

{#if data.length === 0}
	<div class={cn("flex h-56 items-center justify-center", className)}>
		<p class="text-sm text-muted-foreground">No data</p>
	</div>
{:else}
	<div
		class={cn("relative aspect-[3/1] w-full", className)}
		role="group"
		aria-label={ariaLabel}
		onpointerleave={() => (hovered = null)}
	>
		<svg
			viewBox={`0 0 ${W} ${H}`}
			preserveAspectRatio="xMidYMid meet"
			class="absolute inset-0 h-full w-full"
		>
			{#each tickValues(max) as tick}
				<line
					x1={PAD_L}
					x2={W - PAD_R}
					y1={y(tick)}
					y2={y(tick)}
					class="stroke-border"
					stroke-width="1"
				/>
				<text
					x={PAD_L - 6}
					y={y(tick)}
					class="fill-muted-foreground"
					font-size="10"
					text-anchor="end"
					dominant-baseline="middle"
				>
					{formatTick(tick)}
				</text>
			{/each}

			{#each data as d, i}
				{#each series as s}
					<circle
						cx={x(i)}
						cy={y(d.values[s.key] ?? 0)}
						r="2.5"
						style={`fill: ${s.color}`}
						aria-label={`${d.label}: ${s.label} ${d.values[s.key] ?? 0}`}
					/>
				{/each}
			{/each}

			{#each series as s}
				<path
					d={areaPath(s)}
					style={`fill: ${s.color}`}
					fill-opacity="0.12"
					stroke="none"
				/>
				<path
					d={linePath(s)}
					style={`stroke: ${s.color}`}
					fill="none"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				/>
			{/each}

			{#each data as d, i}
				<circle
					cx={x(i)}
					cy={PAD_T + plotH / 2}
					r="10"
					fill="transparent"
					role="button"
					tabindex="0"
					aria-label={`${d.label}: ${series.map((s) => `${s.label} ${d.values[s.key] ?? 0}`).join(", ")}`}
					onpointerenter={() => (hovered = i)}
					onfocus={() => (hovered = i)}
					onblur={() => (hovered = null)}
				/>
			{/each}

			{#each labels as i}
				<text
					x={x(i)}
					y={H - 8}
					class="fill-muted-foreground"
					font-size="10"
					text-anchor="middle"
				>
					{data[i].label}
				</text>
			{/each}
		</svg>

		{#if hovered !== null && data[hovered]}
			<div
				class="pointer-events-none absolute top-2 z-10 rounded-md border bg-card px-2 py-1 text-xs shadow-md"
				style={tooltipStyle(hovered)}
			>
				<p class="font-medium text-foreground">{data[hovered].label}</p>
				{#each series as s}
					<p class="flex items-center gap-1.5 text-muted-foreground">
						<span
							class="inline-block h-2 w-2 rounded-full"
							style={`background: ${s.color}`}
						></span>
						{s.label}
						{data[hovered].values[s.key] ?? 0}
					</p>
				{/each}
			</div>
		{/if}
	</div>
{/if}
