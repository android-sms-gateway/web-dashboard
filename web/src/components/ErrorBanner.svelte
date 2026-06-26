<script lang="ts">
	import { slide } from "svelte/transition";

	let {
		message,
		onRetry,
		onDismiss,
	}: {
		message: string;
		onRetry?: () => void;
		onDismiss?: () => void;
	} = $props();

	let dismissed = $state(false);

	$effect(() => {
		message;
		dismissed = false;
	});

	function handleDismiss() {
		dismissed = true;
		onDismiss?.();
	}
</script>

{#if !dismissed}
	<div class="error-banner" transition:slide={{ duration: 200 }}>
		<div class="error-icon">!</div>
		<span class="error-text">{message}</span>
		<div class="error-actions">
			{#if onRetry}
				<button class="btn-retry" onclick={onRetry}>Retry</button>
			{/if}
			<button class="btn-dismiss" onclick={handleDismiss} aria-label="Dismiss">×</button>
		</div>
	</div>
{/if}

<style>
	.error-banner {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		background: #fef2f2;
		color: #dc2626;
		padding: 0.625rem 0.875rem;
		border-radius: 8px;
		font-size: 0.875rem;
		margin-bottom: 1rem;
		border: 1px solid #fecaca;
	}

	.error-icon {
		flex-shrink: 0;
		width: 22px;
		height: 22px;
		border-radius: 50%;
		background: #dc2626;
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.75rem;
		font-weight: 700;
	}

	.error-text {
		flex: 1;
		line-height: 1.4;
	}

	.error-actions {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		flex-shrink: 0;
	}

	.btn-retry {
		padding: 0.25rem 0.75rem;
		background: #fee2e2;
		color: #991b1b;
		border: 1px solid #fecaca;
		border-radius: 4px;
		font-size: 0.8125rem;
		cursor: pointer;
		font-weight: 500;
	}

	.btn-retry:hover {
		background: #fecaca;
	}

	.btn-dismiss {
		background: none;
		border: none;
		color: #dc2626;
		font-size: 1.25rem;
		cursor: pointer;
		line-height: 1;
		padding: 0 0.25rem;
		opacity: 0.7;
	}

	.btn-dismiss:hover {
		opacity: 1;
	}
</style>
