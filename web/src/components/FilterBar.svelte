<script lang="ts">
	import DeviceSelect from "./DeviceSelect.svelte";

	export type FilterValue = {
		state: string;
		deviceId: string;
		from: string;
		to: string;
	};

	let {
		value = $bindable(),
		onChange,
		stateOptions,
	}: {
		value: FilterValue;
		onChange: (v: FilterValue) => void;
		stateOptions: { value: string; label: string }[];
	} = $props();

	function setField<K extends keyof FilterValue>(key: K, val: FilterValue[K]) {
		const next = { ...value, [key]: val };
		value = next;
		onChange(next);
	}

	function reset() {
		const cleared: FilterValue = { state: "", deviceId: "", from: "", to: "" };
		value = cleared;
		onChange(cleared);
	}

	let isActive = $derived(
		value.state !== "" || value.deviceId !== "" || value.from !== "" || value.to !== "",
	);
</script>

<div class="filter-bar">
	<div class="field">
		<label for="filter-state">State</label>
		<select
			id="filter-state"
			value={value.state}
			onchange={(e) => setField('state', (e.target as HTMLSelectElement).value)}
		>
			{#each stateOptions as opt}
				<option value={opt.value}>{opt.label}</option>
			{/each}
		</select>
	</div>

	<div class="field">
		<label for="filter-device">Device</label>
		<DeviceSelect
			id="filter-device"
			value={value.deviceId}
			onValueChange={(v) => setField('deviceId', v)}
			allLabel="All devices"
		/>
	</div>

	<div class="field">
		<label for="filter-from">From</label>
		<input
			id="filter-from"
			type="date"
			value={value.from}
			onchange={(e) => setField('from', (e.target as HTMLInputElement).value)}
		/>
	</div>

	<div class="field">
		<label for="filter-to">To</label>
		<input
			id="filter-to"
			type="date"
			value={value.to}
			onchange={(e) => setField('to', (e.target as HTMLInputElement).value)}
		/>
	</div>

	{#if isActive}
		<button class="btn-reset" onclick={reset}>Reset filters</button>
	{/if}
</div>

<style>
	.filter-bar {
		display: flex;
		align-items: flex-end;
		gap: 1rem;
		margin-bottom: 1rem;
		padding: 0.75rem 1rem;
		background: white;
		border-radius: 8px;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
		flex-wrap: wrap;
	}

	.field {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	label {
		font-size: 0.75rem;
		color: #6b7280;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	select,
	input[type="date"] {
		padding: 0.375rem 0.625rem;
		border: 1px solid #d1d5db;
		border-radius: 4px;
		font-size: 0.875rem;
		font-family: inherit;
	}

	select:focus,
	input[type="date"]:focus {
		outline: none;
		border-color: #2563eb;
		box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
	}

	.btn-reset {
		align-self: flex-end;
		padding: 0.375rem 0.75rem;
		background: #f3f4f6;
		border: 1px solid #d1d5db;
		border-radius: 4px;
		font-size: 0.8125rem;
		color: #374151;
		cursor: pointer;
	}

	.btn-reset:hover {
		background: #e5e7eb;
	}
</style>