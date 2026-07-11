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
		throw new Error("TabsTrigger must be used within a Tabs component.");
	const isActive = $derived(ctx.value === tabValue);

	function handleKeydown(e: KeyboardEvent) {
		if (!['ArrowLeft', 'ArrowRight', 'Home', 'End'].includes(e.key)) return;
		const tablist = (e.target as HTMLElement).closest('[role="tablist"]');
		if (!tablist) return;
		const tabs = Array.from(tablist.querySelectorAll<HTMLElement>('[role="tab"]'));
		const idx = tabs.indexOf(e.target as HTMLElement);
		if (idx === -1) return;
		e.preventDefault();
		let next = idx;
		switch (e.key) {
			case 'ArrowRight':
				next = (idx + 1) % tabs.length;
				break;
			case 'ArrowLeft':
				next = (idx - 1 + tabs.length) % tabs.length;
				break;
			case 'Home':
				next = 0;
				break;
			case 'End':
				next = tabs.length - 1;
				break;
		}
		tabs[next].focus();
		tabs[next].click();
	}
</script>

<button
	class={cn(
		"inline-flex items-center justify-center whitespace-nowrap rounded-sm px-3 py-1.5 text-sm font-medium ring-offset-background transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
		isActive && "bg-background text-foreground shadow-sm",
		className,
	)}
	id="tab-trigger-{tabValue}"
	role="tab"
	aria-selected={isActive}
	aria-controls="tab-panel-{tabValue}"
	tabindex={isActive ? 0 : -1}
	onclick={() => ctx.setValue(tabValue)}
	onkeydown={handleKeydown}
>
	{#if children}
		{@render children()}
	{/if}
</button>
