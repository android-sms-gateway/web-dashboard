<script lang="ts">
  import { createToken, revokeToken } from "../lib/api";
  import type { CreateTokenResponse } from "../lib/types";

  type View = "create" | "result";

  let view = $state<View>("create");
  let ttlInput = $state("");
  let selectedScopes = $state<string[]>([]);
  let saving = $state(false);
  let error = $state("");
  let result = $state<CreateTokenResponse | null>(null);
  let copied = $state(false);

  let revokeJti = $state("");
  let revoking = $state(false);
  let revokeError = $state("");
  let revokeSuccess = $state(false);

  const scopeList = [
    { value: "devices:list", label: "List Devices" },
    { value: "devices:delete", label: "Delete Devices" },
    { value: "inbox:list", label: "List Inbox Messages" },
    { value: "inbox:refresh", label: "Refresh Inbox" },
    { value: "logs:read", label: "Read Logs" },
    { value: "messages:send", label: "Send Messages" },
    { value: "messages:read", label: "Read Messages" },
    { value: "messages:list", label: "List Messages" },
    { value: "messages:export", label: "Export Messages" },
    { value: "settings:read", label: "Read Settings" },
    { value: "settings:write", label: "Write Settings" },
    { value: "tokens:manage", label: "Manage Tokens" },
    { value: "webhooks:list", label: "List Webhooks" },
    { value: "webhooks:write", label: "Write Webhooks" },
    { value: "webhooks:delete", label: "Delete Webhooks" },
  ];

  function toggleScope(scope: string) {
    if (selectedScopes.includes(scope)) {
      selectedScopes = selectedScopes.filter((s) => s !== scope);
    } else {
      selectedScopes = [...selectedScopes, scope];
    }
  }

  function resetForm() {
    view = "create";
    ttlInput = "";
    selectedScopes = [];
    error = "";
    result = null;
    copied = false;
  }

  async function handleGenerate() {
    error = "";

    if (selectedScopes.length === 0) {
      error = "Select at least one scope";
      return;
    }

    saving = true;
    try {
      result = await createToken({
        scopes: selectedScopes,
        ttl: ttlInput.trim() ? parseInt(ttlInput.trim(), 10) : undefined,
      });
      view = "result";
    } catch {
      error = "Failed to generate token";
    } finally {
      saving = false;
    }
  }

  async function copyToken() {
    if (!result) return;
    try {
      await navigator.clipboard.writeText(result.access_token);
      copied = true;
      setTimeout(() => (copied = false), 3000);
    } catch {
      // clipboard API may not be available
    }
  }

  async function handleRevoke() {
    revokeError = "";
    revokeSuccess = false;

    if (!revokeJti.trim()) {
      revokeError = "Enter a token JTI";
      return;
    }

    revoking = true;
    try {
      await revokeToken(revokeJti.trim());
      revokeSuccess = true;
      revokeJti = "";
      setTimeout(() => (revokeSuccess = false), 3000);
    } catch {
      revokeError = "Failed to revoke token";
    } finally {
      revoking = false;
    }
  }
</script>

<div class="page">
  {#if view === "result" && result}
    <div class="header">
      <h2>Token Generated</h2>
    </div>

    <div class="result-card">
      <div class="warning-banner">
        Save this token now. It will not be shown again.
      </div>

      <label for="token-result">Access Token</label>
      <div class="token-box">
        <code id="token-result" class="token-text">{result.access_token}</code>
        <button class="btn-small" onclick={copyToken}>
          {copied ? "Copied!" : "Copy"}
        </button>
      </div>

      <div class="detail-row">
        <span class="detail-label">ID (JTI)</span>
        <span class="detail-value mono">{result.id}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">Type</span>
        <span class="detail-value">{result.token_type}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">Expires At</span>
        <span class="detail-value"
          >{new Date(result.expires_at).toLocaleString()}</span
        >
      </div>
      {#if result.refresh_token}
        <div class="detail-row">
          <span class="detail-label">Refresh Token</span>
          <span class="detail-value mono">{result.refresh_token}</span>
        </div>
      {/if}

      <button
        class="btn-primary"
        onclick={resetForm}
        style="margin-top:1.5rem;width:100%;"
      >
        Done
      </button>
    </div>
  {:else}
    <div class="header">
      <h2>Generate API Token</h2>
    </div>

    {#if error}
      <div class="error-msg">{error}</div>
    {/if}

    <div class="form-card">
      <label for="ttl">TTL (seconds, optional)</label>
      <input
        id="ttl"
        type="number"
        min="1"
        bind:value={ttlInput}
        placeholder="Leave empty for default"
        disabled={saving}
      />

      <fieldset class="scopes-fieldset">
        <legend class="scopes-legend">Scopes</legend>
        <div class="scope-grid">
          {#each scopeList as scope}
            <label class="scope-item">
              <input
                type="checkbox"
                checked={selectedScopes.includes(scope.value)}
                onchange={() => toggleScope(scope.value)}
                disabled={saving}
              />
              {scope.label}
            </label>
          {/each}
        </div>
      </fieldset>

      <button
        class="btn-primary"
        onclick={handleGenerate}
        disabled={saving || selectedScopes.length === 0}
      >
        {saving ? "Generating..." : "Generate Token"}
      </button>
    </div>
  {/if}

  <details class="revoke-section">
    <summary>Revoke an existing token</summary>
    <div class="revoke-body">
      {#if revokeSuccess}
        <div class="success-msg">Token revoked.</div>
      {/if}
      {#if revokeError}
        <div class="error-msg">{revokeError}</div>
      {/if}

      <label for="revoke-jti">Token JTI</label>
      <input
        id="revoke-jti"
        type="text"
        bind:value={revokeJti}
        placeholder="Paste the token JTI"
        disabled={revoking}
      />
      <button class="btn-danger" onclick={handleRevoke} disabled={revoking}>
        {revoking ? "Revoking..." : "Revoke"}
      </button>
    </div>
  </details>
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

  .error-msg {
    background: #fef2f2;
    color: #dc2626;
    padding: 0.5rem 0.75rem;
    border-radius: 6px;
    font-size: 0.875rem;
    margin-bottom: 1rem;
  }

  .success-msg {
    background: #d1fae5;
    color: #065f46;
    padding: 0.5rem 0.75rem;
    border-radius: 6px;
    font-size: 0.875rem;
    margin-bottom: 0.75rem;
  }

  .form-card {
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  }

  .form-card label {
    display: block;
    font-size: 0.875rem;
    font-weight: 600;
    color: #374151;
    margin-bottom: 0.375rem;
    margin-top: 1rem;
  }

  .form-card label:first-of-type {
    margin-top: 0;
  }

  .form-card input[type="number"] {
    width: 100%;
    padding: 0.5rem 0.75rem;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 0.875rem;
    box-sizing: border-box;
  }

  .form-card input:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
  }

  .scopes-fieldset {
    border: none;
    padding: 0;
    margin: 0;
  }

  .scopes-legend {
    display: block;
    font-size: 0.875rem;
    font-weight: 600;
    color: #374151;
    margin-bottom: 0.375rem;
    margin-top: 1rem;
  }

  .scope-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 0.375rem;
    margin-top: 0.5rem;
  }

  .scope-item {
    display: flex !important;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.8125rem !important;
    font-weight: 400 !important;
    color: #374151 !important;
    cursor: pointer;
    margin: 0 !important;
    padding: 0.25rem 0;
  }

  .scope-item input[type="checkbox"] {
    width: auto;
    cursor: pointer;
  }

  .btn-primary {
    padding: 0.5rem 1.5rem;
    background: #2563eb;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 0.875rem;
    cursor: pointer;
    margin-top: 1rem;
  }

  .btn-primary:hover {
    background: #1d4ed8;
  }

  .btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .result-card {
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  }

  .warning-banner {
    background: #fef3c7;
    color: #92400e;
    padding: 0.75rem 1rem;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 600;
    margin-bottom: 1.25rem;
    text-align: center;
  }

  .result-card label {
    display: block;
    font-size: 0.8125rem;
    font-weight: 600;
    color: #6b7280;
    margin-bottom: 0.375rem;
  }

  .token-box {
    display: flex;
    align-items: stretch;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }

  .token-text {
    flex: 1;
    padding: 0.625rem 0.75rem;
    background: #f3f4f6;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 0.75rem;
    word-break: break-all;
    line-height: 1.4;
    font-family: "SF Mono", "Fira Code", monospace;
  }

  .btn-small {
    padding: 0.375rem 0.75rem;
    background: #e5e7eb;
    color: #374151;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 0.8125rem;
    cursor: pointer;
    white-space: nowrap;
  }

  .btn-small:hover {
    background: #d1d5db;
  }

  .detail-row {
    display: flex;
    gap: 1rem;
    padding: 0.5rem 0;
    border-bottom: 1px solid #f3f4f6;
  }

  .detail-row:last-of-type {
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

  .mono {
    font-family: "SF Mono", "Fira Code", monospace;
    font-size: 0.8125rem;
  }

  .revoke-section {
    margin-top: 2rem;
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    overflow: hidden;
  }

  .revoke-section summary {
    padding: 0.875rem 1.25rem;
    font-size: 0.9375rem;
    font-weight: 600;
    color: #6b7280;
    cursor: pointer;
    user-select: none;
  }

  .revoke-section summary:hover {
    background: #f9fafb;
  }

  .revoke-body {
    padding: 0 1.25rem 1.25rem;
  }

  .revoke-body label {
    display: block;
    font-size: 0.8125rem;
    font-weight: 600;
    color: #6b7280;
    margin-bottom: 0.375rem;
  }

  .revoke-body input {
    width: 100%;
    padding: 0.5rem 0.75rem;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 0.875rem;
    box-sizing: border-box;
    margin-bottom: 0.75rem;
  }

  .revoke-body input:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
  }

  .btn-danger {
    padding: 0.5rem 1rem;
    background: #dc2626;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 0.875rem;
    cursor: pointer;
  }

  .btn-danger:hover {
    background: #b91c1c;
  }

  .btn-danger:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
