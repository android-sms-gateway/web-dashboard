<script lang="ts">
  import Sidebar from "../components/Sidebar.svelte";
  import DashboardPage from "./DashboardPage.svelte";
  import MessagesPage from "./MessagesPage.svelte";
  import DevicesPage from "./DevicesPage.svelte";

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
