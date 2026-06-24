<script lang="ts">
	let {
		total,
		offset,
		limit,
		onPageChange,
		onLimitChange,
	}: {
		total: number;
		offset: number;
		limit: number;
		onPageChange: (offset: number) => void;
		onLimitChange: (limit: number) => void;
	} = $props();

	const PAGE_SIZE_OPTIONS = [25, 50, 100, 200];

	let page = $derived(Math.floor(offset / limit) + 1);
	let totalPages = $derived(Math.max(1, Math.ceil(total / limit)));

	function go(delta: number) {
		const newOffset = Math.max(
			0,
			Math.min(offset + delta * limit, (totalPages - 1) * limit),
		);
		if (newOffset !== offset) onPageChange(newOffset);
	}

	function setLimit(e: Event) {
		const newLimit = parseInt((e.target as HTMLSelectElement).value, 10);
		if (!Number.isFinite(newLimit) || newLimit === limit) return;
		onLimitChange(newLimit);
	}
</script>

{#if total > 0}
	<div class="pagination">
		<button class="btn-page" disabled={page <= 1} onclick={() => go(-1)}
			>&larr; Prev</button
		>
		<span class="page-info">
			Page {page} of {totalPages}
			<span class="total-msg">({total} total)</span>
		</span>
		<button class="btn-page" disabled={page >= totalPages} onclick={() => go(1)}
			>Next &rarr;</button
		>

		<span class="size-label">Show:</span>
		<select value={limit} onchange={setLimit}>
			{#each PAGE_SIZE_OPTIONS as n}
				<option value={n}>{n}</option>
			{/each}
		</select>
	</div>
{/if}

<style>
	.pagination {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		margin-top: 1rem;
		padding-top: 1rem;
		font-size: 0.875rem;
		color: #374151;
	}

	.btn-page {
		padding: 0.375rem 0.75rem;
		background: white;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		cursor: pointer;
	}

	.btn-page:hover:not(:disabled) {
		background: #f9fafb;
	}

	.btn-page:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.page-info {
		color: #6b7280;
	}

	.total-msg {
		margin-left: 0.25rem;
		color: #9ca3af;
	}

	.size-label {
		margin-left: auto;
		color: #6b7280;
	}

	select {
		padding: 0.25rem 0.5rem;
		border: 1px solid #d1d5db;
		border-radius: 4px;
		font-size: 0.875rem;
	}
</style>