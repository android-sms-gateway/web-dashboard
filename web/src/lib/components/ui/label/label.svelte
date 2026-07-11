<script lang="ts">
	import { cn } from '$lib/utils';

	interface Props {
		class?: string;
		for?: string;
		children?: import('svelte').Snippet;
	}

	let {
		class: className,
		for: htmlFor,
		children,
	}: Props = $props();
</script>

<label
	class={cn(
		'text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70',
		className,
	)}
	for={htmlFor}
>
	{@render children?.()}
</label>
