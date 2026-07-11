<script lang="ts">
	import { cn } from '$lib/utils';

	interface Props {
		class?: string;
		children?: import('svelte').Snippet;
	}

	let { class: className, children }: Props = $props();
</script>

<tr class={cn('border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted', className)}>
	{@render children?.()}
</tr>
