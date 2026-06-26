<script lang="ts">
	import { fade, scale } from "svelte/transition";

	let {
		open,
		title,
		message,
		confirmLabel = "Delete",
		cancelLabel = "Cancel",
		variant = "danger",
		onConfirm,
		onCancel,
	}: {
		open: boolean;
		title: string;
		message: string;
		confirmLabel?: string;
		cancelLabel?: string;
		variant?: "danger" | "default";
		onConfirm: () => void;
		onCancel: () => void;
	} = $props();

	let titleId = $state("");

	$effect(() => {
		if (open) titleId = `confirm-title-${crypto.randomUUID()}`;
	});

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === "Escape") onCancel();
	}

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) onCancel();
	}
</script>

{#if open}
	<div class="backdrop" onclick={handleBackdropClick} role="presentation" transition:fade={{ duration: 150 }}>
		<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
		<div class="dialog" role="dialog" aria-modal="true" aria-labelledby={titleId} tabindex="-1" onkeydown={handleKeydown} transition:scale={{ start: 0.95, duration: 150 }}>
			<h3 id={titleId} class="dialog-title">{title}</h3>
			<p class="dialog-message">{message}</p>
			<div class="dialog-actions">
				<button class="btn-cancel" onclick={onCancel}>{cancelLabel}</button>
				<button
					class="btn-confirm"
					class:btn-danger={variant === "danger"}
					class:btn-primary={variant === "default"}
					onclick={onConfirm}
				>
					{confirmLabel}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.backdrop {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.4);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 999;
	}

	.dialog {
		background: white;
		border-radius: 10px;
		padding: 1.5rem;
		min-width: 320px;
		max-width: 440px;
		box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
	}

	.dialog-title {
		font-size: 1.0625rem;
		font-weight: 600;
		color: #111827;
		margin: 0 0 0.5rem;
	}

	.dialog-message {
		font-size: 0.875rem;
		color: #6b7280;
		margin: 0 0 1.5rem;
		line-height: 1.5;
	}

	.dialog-actions {
		display: flex;
		justify-content: flex-end;
		gap: 0.5rem;
	}

	.btn-cancel {
		padding: 0.5rem 1rem;
		background: white;
		color: #374151;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font-size: 0.875rem;
		cursor: pointer;
	}

	.btn-cancel:hover {
		background: #f9fafb;
	}

	.btn-confirm {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 6px;
		font-size: 0.875rem;
		cursor: pointer;
		font-weight: 500;
	}

	.btn-danger {
		background: #dc2626;
		color: white;
	}

	.btn-danger:hover {
		background: #b91c1c;
	}

	.btn-primary {
		background: #2563eb;
		color: white;
	}

	.btn-primary:hover {
		background: #1d4ed8;
	}
</style>
