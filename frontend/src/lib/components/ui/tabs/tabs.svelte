<script lang="ts">
	import { setContext } from 'svelte';
	import { cn } from '$lib/utils';

	interface Props {
		value?: string;
		class?: string;
		children?: import('svelte').Snippet;
	}

	let { value = $bindable(), class: className, children }: Props = $props();

	setContext('tabs', {
		get value() { return value; },
		setValue: (v: string) => value = v,
	});
</script>

<div class={cn('w-full', className)}>
	{#if children}
		{@render children()}
	{/if}
</div>
