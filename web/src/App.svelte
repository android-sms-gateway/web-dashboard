<script lang="ts">
  import { onMount } from "svelte";
  import Login from "./pages/Login.svelte";
  import Dashboard from "./pages/Dashboard.svelte";
  import { me, logout } from "./lib/api";

  let currentPage = $state("loading");

  onMount(async () => {
    try {
      await me();
      currentPage = "dashboard";
    } catch {
      currentPage = "login";
    }
  });

  function handleLogin() {
    currentPage = "dashboard";
  }

  async function handleLogout() {
    try {
      await logout();
    } catch {
      // Server session may already be gone
    }
    currentPage = "login";
  }
</script>

{#if currentPage === "loading"}
  <div class="loading-screen">Loading...</div>
{:else if currentPage === "login"}
  <Login onLogin={handleLogin} />
{:else}
  <Dashboard onLogout={handleLogout} />
{/if}

<style>
  .loading-screen {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
      sans-serif;
    color: #6b7280;
    font-size: 0.9rem;
  }
</style>
