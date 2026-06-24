<script lang="ts">
  import { onMount } from "svelte";
  import { listMessages, getMessage, sendMessage } from "../lib/api";
  import DeviceSelect from "../components/DeviceSelect.svelte";
  import DeviceLabel from "../components/DeviceLabel.svelte";
  import FilterBar from "../components/FilterBar.svelte";
  import Pagination from "../components/Pagination.svelte";
  import type { FilterValue } from "../components/FilterBar.svelte";
  import type {
    MessageListItem,
    MessageDetail,
    SendMessageRequest,
  } from "../lib/types";

  type View = "list" | "compose" | "detail";

  let view = $state<View>("list");
  let messages = $state<MessageListItem[]>([]);
  let total = $state(0);
  let loading = $state(true);
  let loadError = $state("");
  let selectedMessage = $state<MessageDetail | null>(null);
  let loadingDetail = $state(false);
  let detailError = $state("");

  let phoneInput = $state("");
  let textInput = $state("");
  let deviceIdFilter: string = $state("");
  let simInput: string = $state("");
  let sending = $state(false);
  let sendError = $state("");

  let offset = $state(0);
  let limit = $state(50);
  let loadSeq = 0;

  let filterValue = $state<FilterValue>({
    state: "",
    deviceId: "",
    from: "",
    to: "",
  });

  const stateOptions = [
    { value: "", label: "All states" },
    { value: "Pending", label: "Pending" },
    { value: "Processed", label: "Processed" },
    { value: "Sent", label: "Sent" },
    { value: "Delivered", label: "Delivered" },
    { value: "Failed", label: "Failed" },
  ];

  onMount(loadMessages);

  async function loadMessages() {
    const seq = ++loadSeq;
    loading = true;
    loadError = "";
    try {
      const res = await listMessages({
        limit,
        offset,
        state: filterValue.state || undefined,
        deviceId: filterValue.deviceId || undefined,
        from: filterValue.from
          ? new Date(`${filterValue.from}T00:00:00`).toISOString()
          : undefined,
        to: filterValue.to
          ? new Date(`${filterValue.to}T23:59:59`).toISOString()
          : undefined,
      });
      if (seq !== loadSeq) return;
      messages = res.items;
      total = res.total;
    } catch {
      if (seq !== loadSeq) return;
      loadError = "Failed to load messages";
      messages = [];
      total = 0;
    } finally {
      if (seq === loadSeq) loading = false;
    }
  }

  function onFilterChange(v: FilterValue) {
    filterValue = v;
    offset = 0;
    loadMessages();
  }

  function onPageChange(newOffset: number) {
    offset = newOffset;
    loadMessages();
  }

  function onLimitChange(newLimit: number) {
    limit = newLimit;
    offset = 0;
    loadMessages();
  }

  async function openDetail(id: string) {
    view = "detail";
    loadingDetail = true;
    selectedMessage = null;
    detailError = "";
    try {
      selectedMessage = await getMessage(id);
    } catch {
      selectedMessage = null;
      detailError = "Failed to load message details";
    } finally {
      loadingDetail = false;
    }
  }

  function openCompose() {
    view = "compose";
    phoneInput = "";
    textInput = "";
    deviceIdFilter = "";
    simInput = "";
    sendError = "";
  }

  function backToList() {
    view = "list";
    selectedMessage = null;
    loadMessages();
  }

  async function handleSend() {
    sendError = "";
    const phones = phoneInput
      .split(/[\n,]+/)
      .map((p) => p.trim())
      .filter((p) => p.length > 0);

    if (phones.length === 0) {
      sendError = "Enter at least one phone number";
      return;
    }

    if (!textInput.trim()) {
      sendError = "Message text is required";
      return;
    }

    const data: SendMessageRequest = {
      phoneNumbers: phones,
      text: textInput.trim(),
    };

    if (deviceIdFilter) {
      data.deviceId = deviceIdFilter;
    }
    if (simInput.trim()) {
      const n = parseInt(simInput.trim(), 10);
      if (n >= 1 && n <= 3) data.simNumber = n;
    }

    sending = true;
    try {
      await sendMessage(data);
      backToList();
    } catch {
      sendError = "Failed to send message";
    } finally {
      sending = false;
    }
  }

  function stateLabel(state: string): string {
    switch (state) {
      case "Pending":
        return "Pending";
      case "Processed":
        return "Processed";
      case "Sent":
        return "Sent";
      case "Delivered":
        return "Delivered";
      case "Failed":
        return "Failed";
      default:
        return state;
    }
  }

  function stateClass(state: string): string {
    switch (state) {
      case "Delivered":
        return "badge-green";
      case "Failed":
        return "badge-red";
      case "Sent":
        return "badge-blue";
      case "Processed":
        return "badge-yellow";
      default:
        return "badge-gray";
    }
  }
</script>

<div class="page">
  {#if view === "compose"}
    <div class="header">
      <button class="btn-text" onclick={backToList}>&larr; Back</button>
      <h2>Send Message</h2>
    </div>

    <form
      class="compose-form"
      onsubmit={(e) => {
        e.preventDefault();
        handleSend();
      }}
    >
      {#if sendError}
        <div class="error-msg">{sendError}</div>
      {/if}

      <label for="phones">Phone Numbers (one per line or comma-separated)</label
      >
      <textarea
        id="phones"
        bind:value={phoneInput}
        rows={4}
        placeholder="+79990001234&#10;+79990005678"
        disabled={sending}
      ></textarea>

      <label for="text">Message Text</label>
      <textarea
        id="text"
        bind:value={textInput}
        rows={5}
        placeholder="Enter your message..."
        disabled={sending}
      ></textarea>

      <label for="device-select">Device</label>
      <DeviceSelect
        id="device-select"
        value={deviceIdFilter}
        onValueChange={(v) => (deviceIdFilter = v)}
        disabled={sending}
        allLabel="Any device"
      />

      <label for="sim">SIM Number (optional, 1-3)</label>
      <input
        id="sim"
        type="text"
        inputmode="numeric"
        bind:value={simInput}
        placeholder="Leave empty for default"
        disabled={sending}
      />

      <button type="submit" class="btn-primary" disabled={sending}>
        {sending ? "Sending..." : "Send Message"}
      </button>
    </form>
  {:else if view === "detail"}
    <div class="header">
      <button class="btn-text" onclick={backToList}>&larr; Back to list</button>
      <h2>Message Detail</h2>
    </div>

    {#if loadingDetail}
      <p>Loading...</p>
    {:else if detailError}
      <p class="empty">{detailError}</p>
    {:else if selectedMessage}
      <div class="detail-card">
        <div class="detail-row">
          <span class="detail-label">ID</span>
          <span class="detail-value mono">{selectedMessage.id}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Device</span>
          <span class="detail-value"
            ><DeviceLabel deviceId={selectedMessage.deviceId} /></span
          >
        </div>
        <div class="detail-row">
          <span class="detail-label">State</span>
          <span class="detail-value">
            <span class="badge {stateClass(selectedMessage.state)}"
              >{stateLabel(selectedMessage.state)}</span
            >
          </span>
        </div>

        {#if selectedMessage.textMessage}
          <div class="detail-row">
            <span class="detail-label">Message</span>
            <span class="detail-value message-text"
              >{selectedMessage.textMessage.text}</span
            >
          </div>
        {/if}
        {#if selectedMessage.dataMessage}
          <div class="detail-row">
            <span class="detail-label">Data</span>
            <span class="detail-value mono"
              >port={selectedMessage.dataMessage.port} data={selectedMessage
                .dataMessage.data}</span
            >
          </div>
        {/if}
        {#if selectedMessage.hashedMessage}
          <div class="detail-row">
            <span class="detail-label">Hash</span>
            <span class="detail-value mono"
              >{selectedMessage.hashedMessage.hash}</span
            >
          </div>
        {/if}
      </div>

      <h3>Recipients</h3>
      <table class="table">
        <thead>
          <tr>
            <th>Phone Number</th>
            <th>State</th>
            <th>Error</th>
          </tr>
        </thead>
        <tbody>
          {#each selectedMessage.recipients as r}
            <tr>
              <td class="mono">{r.phoneNumber}</td>
              <td
                ><span class="badge {stateClass(r.state)}"
                  >{stateLabel(r.state)}</span
                ></td
              >
              <td class="error-cell">{r.error || ""}</td>
            </tr>
          {/each}
        </tbody>
      </table>

      <h3>Timeline</h3>
      <table class="table">
        <thead>
          <tr>
            <th>Event</th>
            <th>Time</th>
          </tr>
        </thead>
        <tbody>
          {#each Object.entries(selectedMessage.states) as [event, time]}
            <tr>
              <td>{event}</td>
              <td>{time}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}
  {:else}
    <div class="header">
      <h2>Messages</h2>
      <button class="btn-primary" onclick={openCompose}>+ Send</button>
    </div>

    <FilterBar value={filterValue} onChange={onFilterChange} {stateOptions} />

    {#if loading}
      <p>Loading...</p>
    {:else if loadError}
      <p class="empty">{loadError}</p>
    {:else if messages.length === 0}
      <p class="empty">No messages yet.</p>
    {:else}
      <p class="total">{total} message{total !== 1 ? "s" : ""}</p>
      <table class="table">
        <thead>
          <tr>
            <th>Preview</th>
            <th>Recipients</th>
            <th>State</th>
            <th>Time</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {#each messages as msg}
            <tr>
              <td class="preview">{msg.textPreview || "(binary)"}</td>
              <td>{msg.recipients.length}</td>
              <td
                ><span class="badge {stateClass(msg.state)}"
                  >{stateLabel(msg.state)}</span
                ></td
              >
              <td class="time"
                >{msg.createdAt
                  ? new Date(msg.createdAt).toLocaleString()
                  : ""}</td
              >
              <td>
                <button class="btn-small" onclick={() => openDetail(msg.id)}
                  >View</button
                >
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
      <Pagination {total} {offset} {limit} {onPageChange} {onLimitChange} />
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

  .preview {
    max-width: 280px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: #374151;
  }

  .time {
    white-space: nowrap;
    color: #6b7280;
    font-size: 0.8125rem;
  }

  .mono {
    font-family: "SF Mono", "Fira Code", monospace;
    font-size: 0.8125rem;
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

  .badge-red {
    background: #fee2e2;
    color: #991b1b;
  }

  .badge-blue {
    background: #dbeafe;
    color: #1e40af;
  }

  .badge-yellow {
    background: #fef3c7;
    color: #92400e;
  }

  .badge-gray {
    background: #f3f4f6;
    color: #6b7280;
  }

  .error-cell {
    color: #dc2626;
    font-size: 0.8125rem;
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

  .btn-small {
    padding: 0.25rem 0.625rem;
    background: #f3f4f6;
    color: #374151;
    border: 1px solid #d1d5db;
    border-radius: 4px;
    font-size: 0.8125rem;
    cursor: pointer;
  }

  .btn-small:hover {
    background: #e5e7eb;
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

  .compose-form {
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    max-width: 600px;
  }

  .compose-form label {
    display: block;
    font-size: 0.875rem;
    font-weight: 600;
    color: #374151;
    margin-bottom: 0.375rem;
    margin-top: 1rem;
  }

  .compose-form label:first-child {
    margin-top: 0;
  }

  .compose-form textarea,
  .compose-form input {
    width: 100%;
    padding: 0.5rem 0.75rem;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 0.875rem;
    box-sizing: border-box;
    font-family: inherit;
  }

  .compose-form textarea:focus,
  .compose-form input:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
  }

  .compose-form button[type="submit"] {
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

  .detail-card {
    background: white;
    padding: 1.25rem;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    margin-bottom: 1.5rem;
  }

  .detail-row {
    display: flex;
    gap: 1rem;
    padding: 0.5rem 0;
    border-bottom: 1px solid #f3f4f6;
  }

  .detail-row:last-child {
    border-bottom: none;
  }

  .detail-label {
    flex: 0 0 100px;
    font-size: 0.8125rem;
    font-weight: 600;
    color: #6b7280;
  }

  .detail-value {
    flex: 1;
    font-size: 0.875rem;
    color: #111827;
    word-break: break-all;
  }

  .message-text {
    white-space: pre-wrap;
  }

  h3 {
    font-size: 1rem;
    color: #111827;
    margin: 1.5rem 0 0.75rem;
  }
</style>
