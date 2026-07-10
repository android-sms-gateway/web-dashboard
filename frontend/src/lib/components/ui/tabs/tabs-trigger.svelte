<script lang="ts">
	import { getContext } from 'svelte';
	import { cn } from '$lib/utils';

	interface Props {
		value: string;
		class?: string;
		children?: import('svelte').Snippet;
	}

	let { value: tabValue, class: className, children }: Props = $props();

	const ctx = getContext<{ get value(): string; setValue: (v: string) => void }>('tabs');
	const isActive = $derived(ctx.value === tabValue);
</script>

<button
	class={cn(
		'inline-flex items-center justify-center whitespace-nowrap rounded-sm px-3 py-1.5 text-sm font-medium ring-offset-background transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
		isActive && 'bg-background text-foreground shadow-sm',
		className,
	)}
	role="tab"
	aria-selected={isActive}
	onclick={() => ctx.setValue(tabValue)}
>
	{#if children}
		{@render children()}
	{/if}
</button>
