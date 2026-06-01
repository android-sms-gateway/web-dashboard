<script lang="ts">
  import { onMount } from "svelte";
  import { dc, loadDevices } from "../lib/device-cache.svelte";

  let {
    deviceId,
    emptyLabel = "(all devices)",
  }: { deviceId?: string; emptyLabel?: string } = $props();

  const resolvedName = $derived(
    dc.devices.find((d) => d.id === deviceId)?.name,
  );

  onMount(() => loadDevices());
</script>

{#if !deviceId}
  {emptyLabel}
{:else if resolvedName}
  {resolvedName} ({deviceId})
{:else}
  {deviceId}
{/if}
