<script lang="ts">
	import { events, dismissToast } from "../lib/events.svelte";

	function variant(type: string): string {
		if (type === "device.status_changed") return "warning";
		return "info";
	}
</script>

<div class="toast-container" role="status" aria-live="polite">
	{#each events.toasts as t (t.id)}
		<div class="toast toast-{variant(t.type)}">
			<div class="toast-type">{t.type}</div>
			<div class="toast-message">{t.message}</div>
			<button class="toast-close" aria-label="Dismiss" onclick={() => dismissToast(t.id)}
				>×</button
			>
		</div>
	{/each}
</div>

<style>
	.toast-container {
		position: fixed;
		top: 1rem;
		right: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		z-index: 1000;
		pointer-events: none;
		max-height: calc(100vh - 2rem);
		overflow-y: auto;
		max-width: 360px;
	}

	.toast {
		background: white;
		border-left: 4px solid #2563eb;
		border-radius: 6px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
		padding: 0.75rem 1rem;
		min-width: 280px;
		display: grid;
		grid-template-columns: 1fr auto;
		gap: 0.25rem 0.75rem;
		pointer-events: auto;
		font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
			sans-serif;
	}

	.toast-info {
		border-left-color: #2563eb;
	}

	.toast-warning {
		border-left-color: #f59e0b;
	}

	.toast-error {
		border-left-color: #dc2626;
	}

	.toast-type {
		font-size: 0.7rem;
		text-transform: uppercase;
		color: #6b7280;
		letter-spacing: 0.05em;
		grid-column: 1;
	}

	.toast-message {
		font-size: 0.875rem;
		color: #111827;
		grid-column: 1;
		word-break: break-word;
	}

	.toast-close {
		grid-column: 2;
		grid-row: 1 / 3;
		background: none;
		border: none;
		font-size: 1.25rem;
		color: #9ca3af;
		cursor: pointer;
		line-height: 1;
		padding: 0 0.25rem;
	}

	.toast-close:hover {
		color: #374151;
	}
</style>
