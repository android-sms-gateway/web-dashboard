<script lang="ts">
  import { onMount } from "svelte";
  import { dc, loadDevices } from "../lib/device-cache.svelte";

  let {
    value = "",
    onValueChange = (_v: string) => {},
    disabled = false,
    allowAll = true,
    allLabel = "All devices",
    id,
  }: {
    value: string;
    onValueChange?: (v: string) => void;
    disabled?: boolean;
    allowAll?: boolean;
    allLabel?: string;
    id?: string;
  } = $props();

  function handleChange(e: Event) {
    const newVal = (e.target as HTMLSelectElement).value;
    value = newVal;
    onValueChange(newVal);
  }

  onMount(() => loadDevices());
</script>

<select {id} {value} onchange={handleChange} disabled={disabled || dc.loading}>
  {#if dc.loading}
    <option value="" disabled>Loading devices...</option>
  {:else if dc.error}
    <option value="" disabled>Failed to load devices</option>
  {:else if dc.devices.length === 0}
    <option value="" disabled>No devices found</option>
  {:else}
    {#if allowAll}
      <option value="">{allLabel}</option>
    {/if}
    {#each dc.devices as device}
      <option value={device.id}>{device.name} ({device.id})</option>
    {/each}
    {#if value && !dc.devices.some((d) => d.id === value)}
      <option {value} disabled>(previous selection — device gone)</option>
    {/if}
  {/if}
</select>
