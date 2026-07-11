<script lang="ts">
	import { cn } from '$lib/utils';

	interface Props {
		class?: string;
		children?: import('svelte').Snippet;
	}

	let { class: className, children }: Props = $props();
</script>

<thead class={cn('[&_tr]:border-b', className)}>
	{@render children?.()}
</thead>
