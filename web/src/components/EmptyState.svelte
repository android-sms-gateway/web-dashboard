<script lang="ts">
	let {
		icon = "📭",
		title,
		description,
		actionLabel,
		onAction,
	}: {
		icon?: string;
		title: string;
		description?: string;
		actionLabel?: string;
		onAction?: () => void;
	} = $props();
</script>

<div class="empty-state">
	<div class="empty-icon">{icon}</div>
	<h3 class="empty-title">{title}</h3>
	{#if description}
		<p class="empty-desc">{description}</p>
	{/if}
	{#if actionLabel && onAction}
		<button class="btn-action" onclick={onAction}>{actionLabel}</button>
	{/if}
</div>

<style>
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 3rem 2rem;
		text-align: center;
		background: white;
		border-radius: 8px;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
	}

	.empty-icon {
		font-size: 2.5rem;
		margin-bottom: 0.75rem;
		line-height: 1;
	}

	.empty-title {
		font-size: 1rem;
		font-weight: 600;
		color: #374151;
		margin: 0 0 0.375rem;
	}

	.empty-desc {
		font-size: 0.875rem;
		color: #9ca3af;
		margin: 0 0 1.25rem;
		max-width: 320px;
		line-height: 1.4;
	}

	.btn-action {
		padding: 0.5rem 1.25rem;
		background: #2563eb;
		color: white;
		border: none;
		border-radius: 6px;
		font-size: 0.875rem;
		cursor: pointer;
		font-weight: 500;
	}

	.btn-action:hover {
		background: #1d4ed8;
	}
</style>
