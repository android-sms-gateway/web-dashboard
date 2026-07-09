<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { logout } from '$lib/api';

	onMount(async () => {
		try {
			await logout();
		} catch {
			// Server session may already be gone
		}
		goto('/login');
	});
</script>

<div class="flex min-h-screen items-center justify-center">
	<p class="text-muted-foreground">Signing out...</p>
</div>
