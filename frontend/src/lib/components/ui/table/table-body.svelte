<script lang="ts">
	import { cn } from '$lib/utils';

	interface Props {
		class?: string;
		children?: import('svelte').Snippet;
	}

	let { class: className, children }: Props = $props();
</script>

<tbody class={cn('[&_tr:last-child]:border-0', className)}>
	{@render children?.()}
</tbody>
