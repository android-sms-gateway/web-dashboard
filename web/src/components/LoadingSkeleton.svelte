<script lang="ts">
	let {
		variant = "text",
		rows = 1,
		height,
	}: {
		variant?: "card" | "table" | "text" | "detail";
		rows?: number;
		height?: string;
	} = $props();

	let minHeight = $derived(height ?? (variant === "card" ? "120px" : variant === "table" ? "40px" : variant === "detail" ? "60px" : "20px"));
</script>

{#if variant === "card"}
	<div class="skeleton-grid" style="--rows: {rows}">
		{#each Array(rows) as _}
			<div class="skeleton card-shape">
				<div class="shimmer card-line card-line-wide"></div>
				<div class="shimmer card-line"></div>
				<div class="shimmer card-line card-line-narrow"></div>
			</div>
		{/each}
	</div>
{:else if variant === "table"}
	<div class="skeleton-table">
		<div class="skeleton table-header">
			<div class="shimmer" style="width: 100%; height: 16px"></div>
		</div>
		{#each Array(rows) as _}
			<div class="skeleton table-row">
				<div class="shimmer" style="width: {30 + Math.random() * 40}%; height: 14px"></div>
				<div class="shimmer" style="width: {20 + Math.random() * 30}%; height: 14px"></div>
				<div class="shimmer" style="width: 60px; height: 14px"></div>
			</div>
		{/each}
	</div>
{:else if variant === "detail"}
	{#each Array(rows) as _}
		<div class="skeleton detail-row">
			<div class="shimmer detail-label-bar"></div>
			<div class="shimmer" style="width: {40 + Math.random() * 40}%; height: 14px"></div>
		</div>
	{/each}
{:else}
	{#each Array(rows) as _}
		<div class="skeleton text-block" style="min-height: {minHeight}">
			<div class="shimmer" style="width: {60 + Math.random() * 30}%; height: 14px"></div>
		</div>
	{/each}
{/if}

<style>
	.skeleton {
		background: #f3f4f6;
		border-radius: 8px;
		overflow: hidden;
		position: relative;
	}

	.shimmer {
		background: linear-gradient(90deg, #f3f4f6 25%, #e5e7eb 50%, #f3f4f6 75%);
		background-size: 200% 100%;
		border-radius: 4px;
		animation: shimmer 1.5s ease-in-out infinite;
	}

	@keyframes shimmer {
		0% { background-position: 200% 0; }
		100% { background-position: -200% 0; }
	}

	.skeleton-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
		gap: 1rem;
	}

	.card-shape {
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.75rem;
	}

	.card-line {
		height: 14px;
		width: 60%;
	}

	.card-line-wide {
		width: 80%;
		height: 28px;
	}

	.card-line-narrow {
		width: 40%;
	}

	.skeleton-table {
		display: flex;
		flex-direction: column;
		gap: 0;
		border-radius: 8px;
		overflow: hidden;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
	}

	.table-header {
		padding: 0.625rem 1rem;
		background: #f9fafb;
		border-radius: 0;
	}

	.table-row {
		padding: 0.625rem 1rem;
		display: flex;
		gap: 1rem;
		align-items: center;
		border-radius: 0;
		border-top: 1px solid #f3f4f6;
	}

	.detail-row {
		padding: 0.75rem 1rem;
		display: flex;
		gap: 1rem;
		align-items: center;
		background: white;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
		margin-bottom: 1px;
		border-radius: 0;
	}

	.detail-row:first-child {
		border-radius: 8px 8px 0 0;
	}

	.detail-row:last-child {
		border-radius: 0 0 8px 8px;
	}

	.detail-label-bar {
		width: 100px;
		height: 12px;
		flex-shrink: 0;
	}

	.text-block {
		padding: 0.75rem 1rem;
		background: white;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
	}
</style>
