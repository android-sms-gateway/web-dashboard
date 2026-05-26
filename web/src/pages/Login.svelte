<script lang="ts">
  let { onLogin }: { onLogin: () => void } = $props();

  let login = $state("");
  let password = $state("");
  let error = $state("");
  let loading = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = "";
    loading = true;

    try {
      const res = await fetch("/api/v1/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ login, password }),
      });

      if (!res.ok) {
        let data: { message?: string } | null = null;
        try {
          data = await res.json();
        } catch {
          // non-JSON error body
        }
        error = data?.message || "Login failed";
        return;
      }

      onLogin();
    } catch (e) {
      error = "Connection error";
    } finally {
      loading = false;
    }
  }
</script>

<div class="login-container">
  <form onsubmit={handleSubmit}>
    <h1>Web Dashboard</h1>
    <p class="subtitle">SMSGate</p>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    <label for="login">Login</label>
    <input
      id="login"
      type="text"
      bind:value={login}
      placeholder="Enter your login"
      disabled={loading}
    />

    <label for="password">Password</label>
    <input
      id="password"
      type="password"
      bind:value={password}
      placeholder="Enter your password"
      disabled={loading}
    />

    <button type="submit" disabled={loading}>
      {loading ? "Signing in..." : "Sign In"}
    </button>
  </form>
</div>

<style>
  .login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background: #f5f5f5;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
      sans-serif;
  }

  form {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 360px;
  }

  h1 {
    margin: 0 0 0.25rem;
    font-size: 1.5rem;
    color: #333;
  }

  .subtitle {
    color: #666;
    margin: 0 0 1.5rem;
    font-size: 0.9rem;
  }

  .error {
    background: #fef2f2;
    color: #dc2626;
    padding: 0.5rem 0.75rem;
    border-radius: 4px;
    margin-bottom: 1rem;
    font-size: 0.875rem;
  }

  label {
    display: block;
    margin-bottom: 0.25rem;
    font-size: 0.875rem;
    color: #555;
  }

  input {
    width: 100%;
    padding: 0.5rem 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    margin-bottom: 1rem;
    font-size: 0.9rem;
    box-sizing: border-box;
  }

  input:focus {
    outline: none;
    border-color: #2563eb;
    box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
  }

  button {
    width: 100%;
    padding: 0.625rem;
    background: #2563eb;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 0.9rem;
    cursor: pointer;
  }

  button:hover {
    background: #1d4ed8;
  }

  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
