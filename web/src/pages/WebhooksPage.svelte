<script lang="ts">
  import { onMount } from "svelte";
  import { listWebhooks, createWebhook, deleteWebhook } from "../lib/api";
  import DeviceSelect from "../components/DeviceSelect.svelte";
  import DeviceLabel from "../components/DeviceLabel.svelte";
  import type { Webhook, CreateWebhookRequest } from "../lib/types";

  type View = "list" | "create";

  let view = $state<View>("list");
  let webhooks = $state<Webhook[]>([]);
  let loading = $state(true);

  let eventInput = $state("sms:received");
  let urlInput = $state("");
  let deviceIdInput = $state("");
  let saving = $state(false);
  let error = $state("");

  const events = [
    { value: "sms:received", label: "SMS Received" },
    { value: "sms:data-received", label: "Data SMS Received" },
    { value: "sms:sent", label: "SMS Sent" },
    { value: "sms:delivered", label: "SMS Delivered" },
    { value: "sms:failed", label: "SMS Failed" },
    { value: "system:ping", label: "System Ping" },
    { value: "mms:received", label: "MMS Received" },
    { value: "mms:downloaded", label: "MMS Downloaded" },
  ];

  onMount(load);

  async function load() {
    loading = true;
    try {
      webhooks = await listWebhooks();
    } catch {
      webhooks = [];
    } finally {
      loading = false;
    }
  }

  function openCreate() {
    view = "create";
    eventInput = "sms:received";
    urlInput = "";
    deviceIdInput = "";
    error = "";
  }

  function backToList() {
    view = "list";
    load();
  }

  async function handleCreate() {
    error = "";

    if (!urlInput.trim()) {
      error = "URL is required";
      return;
    }

    const data: CreateWebhookRequest = {
      url: urlInput.trim(),
      event: eventInput,
    };
    if (deviceIdInput.trim()) {
      data.deviceId = deviceIdInput.trim();
    }

    saving = true;
    try {
      await createWebhook(data);
      backToList();
    } catch {
      error = "Failed to create webhook";
    } finally {
      saving = false;
    }
  }

  let deletingId = $state<string | null>(null);

  async function handleDelete(id: string) {
    if (!confirm("Delete this webhook?")) return;
    deletingId = id;
    try {
      await deleteWebhook(id);
      webhooks = webhooks.filter((w) => w.id !== id);
    } catch {
      // ignore
    } finally {
      deletingId = null;
    }
  }

  function eventLabel(value: string): string {
    const e = events.find((e) => e.value === value);
    return e ? e.label : value;
  }
</script>

<div class="page">
  {#if view === "create"}
    <div class="header">
      <button class="btn-text" onclick={backToList}>&larr; Back</button>
      <h2>Add Webhook</h2>
    </div>

    <form
      class="form"
      onsubmit={(e) => {
        e.preventDefault();
        handleCreate();
      }}
    >
      {#if error}
        <div class="error-msg">{error}</div>
      {/if}

      <label for="event">Event</label>
      <select id="event" bind:value={eventInput} disabled={saving}>
        {#each events as ev}
          <option value={ev.value}>{ev.label}</option>
        {/each}
      </select>

      <label for="url">URL</label>
      <input
        id="url"
        type="url"
        bind:value={urlInput}
        placeholder="https://example.com/webhook"
        disabled={saving}
      />

      <label for="deviceId">Device (optional)</label>
      <DeviceSelect
        id="deviceId"
        value={deviceIdInput}
        onValueChange={(v) => (deviceIdInput = v)}
        disabled={saving}
      />

      <button type="submit" class="btn-primary" disabled={saving}>
        {saving ? "Creating..." : "Create Webhook"}
      </button>
    </form>
  {:else}
    <div class="header">
      <h2>Webhooks</h2>
      <button class="btn-primary" onclick={openCreate}>+ Add</button>
    </div>

    {#if loading}
      <p>Loading...</p>
    {:else if webhooks.length === 0}
      <p class="empty">No webhooks configured.</p>
    {:else}
      <p class="total">
        {webhooks.length} webhook{webhooks.length !== 1 ? "s" : ""}
      </p>
      <table class="table">
        <thead>
          <tr>
            <th>Event</th>
            <th>URL</th>
            <th>Device</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {#each webhooks as w}
            <tr>
              <td
                ><span class="badge badge-blue">{eventLabel(w.event)}</span></td
              >
              <td class="url">{w.url}</td>
              <td class="mono"
                ><DeviceLabel deviceId={w.deviceId} emptyLabel="(all)" /></td
              >
              <td>
                <button
                  class="btn-danger-small"
                  onclick={() => handleDelete(w.id)}
                  disabled={deletingId === w.id}
                >
                  {deletingId === w.id ? "..." : "Delete"}
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}
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

  .url {
    max-width: 400px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: #374151;
  }

  .mono {
    font-family: "SF Mono", "Fira Code", monospace;
    font-size: 0.8125rem;
    color: #6b7280;
  }

  .badge {
    display: inline-block;
    padding: 0.125rem 0.5rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 600;
    white-space: nowrap;
  }

  .badge-blue {
    background: #dbeafe;
    color: #1e40af;
  }

  .btn-primary {
    padding: 0.5rem 1rem;
    background: #2563eb;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 0.875rem;
    cursor: pointer;
  }

  .btn-primary:hover {
    background: #1d4ed8;
  }

  .btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
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

  .btn-text {
    background: none;
    border: none;
    color: #2563eb;
    font-size: 0.875rem;
    cursor: pointer;
    padding: 0;
  }

  .btn-text:hover {
    text-decoration: underline;
  }

  .form {
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    max-width: 500px;
  }

  .form label {
    display: block;
    font-size: 0.875rem;
    font-weight: 600;
    color: #374151;
    margin-bottom: 0.375rem;
    margin-top: 1rem;
  }

  .form label:first-of-type {
    margin-top: 0;
  }

  .form input,
  .form select {
    width: 100%;
    padding: 0.5rem 0.75rem;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 0.875rem;
    box-sizing: border-box;
    font-family: inherit;
  }

  .form input:focus,
  .form select:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
  }

  .form button[type="submit"] {
    margin-top: 1.5rem;
    width: 100%;
  }

  .error-msg {
    background: #fef2f2;
    color: #dc2626;
    padding: 0.5rem 0.75rem;
    border-radius: 6px;
    font-size: 0.875rem;
    margin-bottom: 1rem;
  }
</style>
