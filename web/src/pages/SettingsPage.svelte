<script lang="ts">
  import { onMount } from "svelte";
  import { getSettings, updateSettings } from "../lib/api";
  import type { DeviceSettings } from "../lib/types";

  let loading = $state(true);
  let saving = $state(false);
  let error = $state("");
  let success = $state(false);

  function emptyForm(): DeviceSettings {
    return {
      messages: { sendIntervalMin: undefined, sendIntervalMax: undefined, limitPeriod: undefined, limitValue: undefined, simSelectionMode: undefined, logLifetimeDays: undefined, processingOrder: undefined },
      ping: { intervalSeconds: undefined },
      logs: { lifetimeDays: undefined },
      webhooks: { internetRequired: undefined, retryCount: undefined, signingKey: undefined },
      gateway: { cloudUrl: undefined, privateToken: undefined, notificationChannel: undefined },
      encryption: { passphrase: undefined },
    };
  }

  let form: DeviceSettings = $state(emptyForm());

  onMount(load);

  function mergeNested(target: any, source: any) {
    for (const key of Object.keys(source)) {
      const val = source[key];
      if (val !== null && val !== undefined && typeof val === "object" && !Array.isArray(val)) {
        if (!target[key] || typeof target[key] !== "object") target[key] = {};
        mergeNested(target[key], val);
      } else {
        target[key] = val;
      }
    }
  }

  async function load() {
    loading = true;
    try {
      const data = await getSettings();
      form = emptyForm();
      mergeNested(form, data);
    } catch {
      form = emptyForm();
    } finally {
      loading = false;
    }
  }

  async function handleSave() {
    error = "";
    success = false;
    saving = true;
    try {
      const result = await updateSettings(form);
      form = emptyForm();
      mergeNested(form, result);
      success = true;
      setTimeout(() => success = false, 3000);
    } catch {
      error = "Failed to save settings";
    } finally {
      saving = false;
    }
  }
</script>

<div class="page">
  <div class="header">
    <h2>Settings</h2>
    <button class="btn-primary" onclick={handleSave} disabled={loading || saving}>
      {saving ? "Saving..." : "Save"}
    </button>
  </div>

  {#if success}
    <div class="success-msg">Settings saved.</div>
  {/if}

  {#if error}
    <div class="error-msg">{error}</div>
  {/if}

  {#if loading}
    <p>Loading settings...</p>
  {:else}
    <div class="settings-sections">
      <details class="section" open>
        <summary>Messages</summary>
        <div class="section-body">
          <label for="msg-send-interval-min">Send Interval Min (seconds)</label>
          <input id="msg-send-interval-min" type="number" min="1" placeholder="Leave empty for default" bind:value={form.messages.sendIntervalMin} />

          <label for="msg-send-interval-max">Send Interval Max (seconds)</label>
          <input id="msg-send-interval-max" type="number" min="1" placeholder="Leave empty for default" bind:value={form.messages.sendIntervalMax} />

          <label for="msg-limit-period">Limit Period</label>
          <select id="msg-limit-period" bind:value={form.messages.limitPeriod}>
            <option value={undefined}>Default</option>
            <option value="Disabled">Disabled</option>
            <option value="PerMinute">Per Minute</option>
            <option value="PerHour">Per Hour</option>
            <option value="PerDay">Per Day</option>
          </select>

          <label for="msg-limit-value">Limit Value</label>
          <input id="msg-limit-value" type="number" min="1" placeholder="Leave empty for default" bind:value={form.messages.limitValue} />

          <label for="msg-sim-mode">SIM Selection Mode</label>
          <select id="msg-sim-mode" bind:value={form.messages.simSelectionMode}>
            <option value={undefined}>Default</option>
            <option value="OSDefault">OS Default</option>
            <option value="RoundRobin">Round Robin</option>
            <option value="Random">Random</option>
          </select>

          <label for="msg-log-lifetime">Log Lifetime (days)</label>
          <input id="msg-log-lifetime" type="number" min="1" placeholder="Leave empty for default" bind:value={form.messages.logLifetimeDays} />

          <label for="msg-processing-order">Processing Order</label>
          <select id="msg-processing-order" bind:value={form.messages.processingOrder}>
            <option value={undefined}>Default</option>
            <option value="LIFO">LIFO (last in, first out)</option>
            <option value="FIFO">FIFO (first in, first out)</option>
          </select>
        </div>
      </details>

      <details class="section" open>
        <summary>Ping</summary>
        <div class="section-body">
          <label for="ping-interval">Interval (seconds)</label>
          <input id="ping-interval" type="number" min="1" placeholder="Leave empty for default" bind:value={form.ping.intervalSeconds} />
        </div>
      </details>

      <details class="section" open>
        <summary>Logs</summary>
        <div class="section-body">
          <label for="logs-lifetime">Lifetime (days)</label>
          <input id="logs-lifetime" type="number" min="1" placeholder="Leave empty for default" bind:value={form.logs.lifetimeDays} />
        </div>
      </details>

      <details class="section" open>
        <summary>Webhooks</summary>
        <div class="section-body">
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={form.webhooks.internetRequired} />
            Internet Required
          </label>

          <label for="webhook-retry-count">Retry Count</label>
          <input id="webhook-retry-count" type="number" min="1" placeholder="Leave empty for default" bind:value={form.webhooks.retryCount} />

          <label for="webhook-signing-key">Signing Key</label>
          <input id="webhook-signing-key" type="password" placeholder="Leave empty to keep current" bind:value={form.webhooks.signingKey} />
        </div>
      </details>

      <details class="section" open>
        <summary>Gateway</summary>
        <div class="section-body">
          <label for="gw-cloud-url">Cloud URL</label>
          <input id="gw-cloud-url" type="url" placeholder="Leave empty for default" bind:value={form.gateway.cloudUrl} />

          <label for="gw-private-token">Private Token</label>
          <input id="gw-private-token" type="password" placeholder="Leave empty to keep current" bind:value={form.gateway.privateToken} />

          <label for="gw-notif-channel">Notification Channel</label>
          <select id="gw-notif-channel" bind:value={form.gateway.notificationChannel}>
            <option value={undefined}>Default</option>
            <option value="AUTO">Auto</option>
            <option value="SSE_ONLY">SSE Only</option>
          </select>
        </div>
      </details>

      <details class="section" open>
        <summary>Encryption</summary>
        <div class="section-body">
          <label for="enc-passphrase">Passphrase</label>
          <input id="enc-passphrase" type="password" placeholder="Leave empty to keep current" bind:value={form.encryption.passphrase} />
        </div>
      </details>
    </div>

    <div class="save-bar">
      <button class="btn-primary" onclick={handleSave} disabled={saving}>
        {saving ? "Saving..." : "Save Settings"}
      </button>
    </div>
  {/if}
</div>

<style>
  .page {
    max-width: 720px;
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

  .success-msg {
    background: #d1fae5;
    color: #065f46;
    padding: 0.5rem 0.75rem;
    border-radius: 6px;
    font-size: 0.875rem;
    margin-bottom: 1rem;
  }

  .error-msg {
    background: #fef2f2;
    color: #dc2626;
    padding: 0.5rem 0.75rem;
    border-radius: 6px;
    font-size: 0.875rem;
    margin-bottom: 1rem;
  }

  .settings-sections {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .section {
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    overflow: hidden;
  }

  .section summary {
    padding: 0.875rem 1.25rem;
    font-size: 0.9375rem;
    font-weight: 600;
    color: #111827;
    cursor: pointer;
    user-select: none;
  }

  .section summary:hover {
    background: #f9fafb;
  }

  .section-body {
    padding: 0 1.25rem 1.25rem;
    display: grid;
    grid-template-columns: 1fr;
    gap: 0.25rem;
  }

  .section-body label {
    display: block;
    font-size: 0.8125rem;
    font-weight: 600;
    color: #6b7280;
    margin-top: 0.75rem;
  }

  .section-body label:first-child {
    margin-top: 0;
  }

  .section-body input,
  .section-body select {
    width: 100%;
    padding: 0.5rem 0.75rem;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 0.875rem;
    box-sizing: border-box;
    font-family: inherit;
  }

  .section-body input:focus,
  .section-body select:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
  }

  .checkbox-label {
    display: flex !important;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    margin-top: 0.75rem;
  }

  .checkbox-label input[type="checkbox"] {
    width: auto;
    cursor: pointer;
  }

  .save-bar {
    margin-top: 1.5rem;
  }

  .btn-primary {
    padding: 0.5rem 1.5rem;
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
</style>
