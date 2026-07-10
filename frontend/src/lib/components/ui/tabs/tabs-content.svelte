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

{#if isActive}
	<div class={cn('mt-4 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2', className)} role="tabpanel">
		{#if children}
			{@render children()}
		{/if}
	</div>
{/if}
