<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';

	let {
		open,
		title,
		message,
		confirmLabel = 'Delete',
		cancelLabel = 'Cancel',
		variant = 'destructive' as 'destructive' | 'default',
		onConfirm,
		onCancel,
	}: {
		open: boolean;
		title: string;
		message: string;
		confirmLabel?: string;
		cancelLabel?: string;
		variant?: 'destructive' | 'default';
		onConfirm: () => void;
		onCancel: () => void;
	} = $props();
</script>

{#if open}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		onclick={(e) => { if (e.target === e.currentTarget) onCancel(); }}
		role="presentation"
	>
		<div
			class="w-full max-w-sm rounded-lg border bg-card p-6 shadow-lg"
			role="dialog"
			aria-modal="true"
		>
			<h3 class="text-lg font-semibold">{title}</h3>
			<p class="mt-2 text-sm text-muted-foreground">{message}</p>
			<div class="mt-6 flex justify-end gap-2">
				<Button variant="outline" onclick={onCancel}>{cancelLabel}</Button>
				<Button variant={variant} onclick={onConfirm}>{confirmLabel}</Button>
			</div>
		</div>
	</div>
{/if}
