<script lang="ts">
	import { cn } from "$lib/utils";

	interface Props {
		class?: string;
		children?: import("svelte").Snippet;
	}

	let {
		class: className,
		children,
		...rest
	}: Props &
		import("svelte/elements").HTMLAttributes<HTMLDivElement> = $props();
</script>

<div
	class={cn(
		"inline-flex h-10 items-center justify-center rounded-md bg-muted p-1 text-muted-foreground",
		className,
	)}
	role="tablist"
	{...rest}
>
	{#if children}
		{@render children()}
	{/if}
</div>
