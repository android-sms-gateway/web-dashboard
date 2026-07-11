<script lang="ts">
	import { getContext } from "svelte";
	import { cn } from "$lib/utils";
	import { TABS_KEY } from "./tabs.svelte";

	interface Props {
		value: string;
		class?: string;
		children?: import("svelte").Snippet;
	}

	let { value: tabValue, class: className, children }: Props = $props();

	const ctx = getContext<{
		get value(): string;
		setValue: (v: string) => void;
	}>(TABS_KEY);
	if (!ctx)
		throw new Error("TabsContent must be used within a Tabs component.");
	const isActive = $derived(ctx.value === tabValue);
</script>

{#if isActive}
	<div
		id="tab-panel-{tabValue}"
		class={cn(
			"mt-4 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus:ring-ring focus-visible:ring-offset-2",
			className,
		)}
		role="tabpanel"
		aria-labelledby="tab-trigger-{tabValue}"
		tabindex="0"
	>
		{#if children}
			{@render children()}
		{/if}
	</div>
{/if}
