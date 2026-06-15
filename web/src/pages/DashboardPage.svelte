<script lang="ts">
  import { onMount } from "svelte";

  let stats: {
    messagesSent: number;
    messagesPending: number;
    messagesFailed: number;
    devicesOnline: number;
    devicesActive: number;
    devicesTotal: number;
  } | null = $state(null);
  let userName = $state("");
  let loading = $state(true);

  onMount(async () => {
    try {
      const [meRes, statsRes] = await Promise.all([
        fetch("/api/v1/auth/me"),
        fetch("/api/v1/stats"),
      ]);

      if (meRes.ok) {
        const me = await meRes.json();
        userName = me.login;
      }

      if (statsRes.ok) {
        stats = await statsRes.json();
      }
    } finally {
      loading = false;
    }
  });
</script>

{#if loading}
  <p>Loading...</p>
{:else}
  <div class="welcome">
    Welcome, <strong>{userName || "User"}</strong>
  </div>

  <div class="cards">
    <div class="card">
      <div class="card-value">{stats?.devicesOnline ?? "?"}</div>
      <div class="card-label">Devices Online</div>
    </div>
    <div class="card">
      <div class="card-value">{stats?.devicesActive ?? "?"}</div>
      <div class="card-label">Devices Active</div>
    </div>
    <div class="card">
      <div class="card-value">{stats?.devicesTotal ?? "?"}</div>
      <div class="card-label">Total Devices</div>
    </div>
    <div class="card">
      <div class="card-value">{stats?.messagesSent ?? "?"}</div>
      <div class="card-label">Messages Sent</div>
    </div>
    <div class="card">
      <div class="card-value">{stats?.messagesPending ?? "?"}</div>
      <div class="card-label">Pending</div>
    </div>
  </div>
{/if}

<style>
  .welcome {
    margin-bottom: 1.5rem;
    font-size: 1.1rem;
    color: #374151;
  }

  .cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 1rem;
  }

  .card {
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    text-align: center;
  }

  .card-value {
    font-size: 2rem;
    font-weight: 700;
    color: #2563eb;
  }

  .card-label {
    margin-top: 0.25rem;
    font-size: 0.875rem;
    color: #6b7280;
  }
</style>
