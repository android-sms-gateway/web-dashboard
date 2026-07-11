<script module lang="ts">
	let counter = 0;
</script>

<script lang="ts">
	import { tick } from 'svelte';
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

	let dialogEl: HTMLDivElement | undefined = $state();
	const titleId = `dialog-title-${++counter}`;
	const descId = `dialog-desc-${counter}`;

	$effect(() => {
		if (open) {
			const prev = document.activeElement as HTMLElement;
			tick().then(() => dialogEl?.focus());
			return () => {
				if (prev && prev.isConnected) prev.focus();
			};
		}
	});

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			onCancel();
			return;
		}
		if (e.key === 'Tab' && dialogEl) {
			const nodes = dialogEl.querySelectorAll<HTMLElement>(
				'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
			);
			if (nodes.length === 0) return;
			const first = nodes[0];
			const last = nodes[nodes.length - 1];
			if (e.shiftKey) {
				if (document.activeElement === first) {
					e.preventDefault();
					last.focus();
				}
			} else {
				if (document.activeElement === last) {
					e.preventDefault();
					first.focus();
				}
			}
		}
	}
</script>

{#if open}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		onclick={(e) => { if (e.target === e.currentTarget) onCancel(); }}
		onkeydown={handleKeydown}
		role="presentation"
	>
		<div
			bind:this={dialogEl}
			class="w-full max-w-sm rounded-lg border bg-card p-6 shadow-lg"
			role="dialog"
			aria-modal="true"
			aria-labelledby={titleId}
			aria-describedby={descId}
			tabindex="-1"
		>
			<h3 id={titleId} class="text-lg font-semibold">{title}</h3>
			<p id={descId} class="mt-2 text-sm text-muted-foreground">{message}</p>
			<div class="mt-6 flex justify-end gap-2">
				<Button variant="outline" onclick={onCancel}>{cancelLabel}</Button>
				<Button variant={variant} onclick={onConfirm}>{confirmLabel}</Button>
			</div>
		</div>
	</div>
{/if}
