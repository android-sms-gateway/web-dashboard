<script lang="ts">
  import Sidebar from "../components/Sidebar.svelte";
  import DashboardPage from "./DashboardPage.svelte";
  import MessagesPage from "./MessagesPage.svelte";
  import DevicesPage from "./DevicesPage.svelte";
  import WebhooksPage from "./WebhooksPage.svelte";
  import SettingsPage from "./SettingsPage.svelte";
  import TokensPage from "./TokensPage.svelte";

  let { onLogout }: { onLogout: () => void } = $props();

  let currentPage = $state("dashboard");

  function onNavigate(page: string) {
    currentPage = page;
  }
</script>

<div class="layout">
  <Sidebar current={currentPage} {onNavigate} {onLogout} />
  <main class="content">
    {#if currentPage === "dashboard"}
      <DashboardPage />
    {:else if currentPage === "messages"}
      <MessagesPage />
    {:else if currentPage === "devices"}
      <DevicesPage />
    {:else if currentPage === "webhooks"}
      <WebhooksPage />
    {:else if currentPage === "settings"}
      <SettingsPage />
    {:else if currentPage === "tokens"}
      <TokensPage />
    {:else}
      <h2>{currentPage}</h2>
      <p>This page is coming soon.</p>
    {/if}
  </main>
</div>

<style>
  .layout {
    display: flex;
    min-height: 100vh;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
      sans-serif;
  }

  .content {
    flex: 1;
    padding: 2rem;
    background: #f8fafc;
  }
</style>
