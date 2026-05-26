<script lang="ts">
  let { current, onNavigate, onLogout }: {
    current: string;
    onNavigate: (page: string) => void;
    onLogout: () => void;
  } = $props();

  const items = [
    { id: 'dashboard', label: 'Dashboard', icon: '📊' },
    { id: 'messages', label: 'Messages', icon: '✉️' },
    { id: 'devices', label: 'Devices', icon: '📱' },
    { id: 'webhooks', label: 'Webhooks', icon: '🔗' },
    { id: 'settings', label: 'Settings', icon: '⚙️' },
    { id: 'tokens', label: 'API Tokens', icon: '🔑' },
  ];
</script>

<aside class="sidebar">
  <div class="brand">SMS Dashboard</div>

  <nav>
    {#each items as item}
      <button
        class="nav-item"
        class:active={current === item.id}
        onclick={() => onNavigate(item.id)}
      >
        <span class="icon">{item.icon}</span>
        {item.label}
      </button>
    {/each}
  </nav>

  <div class="spacer"></div>

  <button class="nav-item logout" onclick={onLogout}>
    <span class="icon">🚪</span>
    Logout
  </button>
</aside>

<style>
  .sidebar {
    width: 220px;
    background: #1e293b;
    color: #e2e8f0;
    display: flex;
    flex-direction: column;
    padding: 0;
  }

  .brand {
    padding: 1.25rem 1rem;
    font-size: 1rem;
    font-weight: 700;
    color: #f8fafc;
    border-bottom: 1px solid #334155;
  }

  nav {
    padding: 0.5rem 0;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    width: 100%;
    padding: 0.625rem 1rem;
    background: none;
    border: none;
    color: #94a3b8;
    font-size: 0.875rem;
    cursor: pointer;
    text-align: left;
  }

  .nav-item:hover {
    background: #334155;
    color: #e2e8f0;
  }

  .nav-item.active {
    background: #2563eb;
    color: white;
  }

  .icon {
    font-size: 1rem;
    width: 1.25rem;
    text-align: center;
  }

  .spacer {
    flex: 1;
  }

  .logout {
    border-top: 1px solid #334155;
    margin-top: auto;
  }
</style>
