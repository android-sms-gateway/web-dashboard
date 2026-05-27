<script lang="ts">
  import { onMount } from "svelte";
  import { listDevices, deleteDevice } from "../lib/api";
  import type { Device } from "../lib/types";

  let devices = $state<Device[]>([]);
  let loading = $state(true);

  onMount(load);

  async function load() {
    loading = true;
    try {
      devices = await listDevices();
    } catch {
      devices = [];
    } finally {
      loading = false;
    }
  }

  let deleting = $state<string | null>(null);

  async function handleDelete(id: string) {
    if (!confirm("Delete this device?")) return;
    deleting = id;
    try {
      await deleteDevice(id);
      devices = devices.filter(d => d.id !== id);
    } catch {
      // ignore
    } finally {
      deleting = null;
    }
  }

  function timeAgo(iso: string): string {
    const diff = Date.now() - new Date(iso).getTime();
    const minutes = Math.floor(diff / 60000);
    if (minutes < 1) return "just now";
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`;
    const days = Math.floor(hours / 24);
    return `${days}d ago`;
  }
</script>

<div class="page">
  <div class="header">
    <h2>Devices</h2>
  </div>

  {#if loading}
    <p>Loading...</p>
  {:else if devices.length === 0}
    <p class="empty">No devices registered.</p>
  {:else}
    <p class="total">{devices.length} device{devices.length !== 1 ? "s" : ""}</p>
    <table class="table">
      <thead>
        <tr>
          <th>Name</th>
          <th>ID</th>
          <th>Status</th>
          <th>Last Seen</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#each devices as d}
          <tr>
            <td class="name">{d.name || "(unnamed)"}</td>
            <td class="mono id">{d.id}</td>
            <td>
              {#if d.isOnline}
                <span class="badge badge-green">Online</span>
              {:else}
                <span class="badge badge-gray">Offline</span>
              {/if}
            </td>
            <td class="time">{timeAgo(d.lastSeen)}</td>
            <td>
              <button
                class="btn-danger-small"
                onclick={() => handleDelete(d.id)}
                disabled={deleting === d.id}
              >
                {deleting === d.id ? "..." : "Delete"}
              </button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>

<style>
  .page {
    max-width: 960px;
  }

  .header {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  .header h2 {
    margin: 0;
    font-size: 1.25rem;
    color: #111827;
  }

  .total {
    margin-bottom: 1rem;
    font-size: 0.875rem;
    color: #6b7280;
  }

  .empty {
    color: #9ca3af;
    font-style: italic;
  }

  .table {
    width: 100%;
    border-collapse: collapse;
    background: white;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  }

  .table th,
  .table td {
    padding: 0.625rem 1rem;
    text-align: left;
    font-size: 0.875rem;
  }

  .table th {
    background: #f9fafb;
    color: #6b7280;
    font-weight: 600;
    text-transform: uppercase;
    font-size: 0.75rem;
    letter-spacing: 0.05em;
  }

  .table td {
    border-top: 1px solid #f3f4f6;
  }

  .table tbody tr:hover {
    background: #f9fafb;
  }

  .name {
    font-weight: 500;
    color: #111827;
  }

  .id {
    font-size: 0.8125rem;
    color: #6b7280;
  }

  .time {
    white-space: nowrap;
    color: #6b7280;
    font-size: 0.8125rem;
  }

  .mono {
    font-family: "SF Mono", "Fira Code", monospace;
  }

  .badge {
    display: inline-block;
    padding: 0.125rem 0.5rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 600;
  }

  .badge-green {
    background: #d1fae5;
    color: #065f46;
  }

  .badge-gray {
    background: #f3f4f6;
    color: #6b7280;
  }

  .btn-danger-small {
    padding: 0.25rem 0.625rem;
    background: #fef2f2;
    color: #dc2626;
    border: 1px solid #fecaca;
    border-radius: 4px;
    font-size: 0.8125rem;
    cursor: pointer;
  }

  .btn-danger-small:hover {
    background: #fee2e2;
  }

  .btn-danger-small:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
