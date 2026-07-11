<script lang="ts">
	import { goto } from '$app/navigation';
	import { login } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import {
		Card,
		CardContent,
		CardDescription,
		CardFooter,
		CardHeader,
		CardTitle,
	} from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Alert, AlertDescription } from '$lib/components/ui/alert/index.js';

	let username = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;

		try {
			await login({ login: username, password });
			goto('/');
		} catch (e: unknown) {
			error = e instanceof Error ? e.message : 'Login failed';
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center">
	<Card class="w-full max-w-sm">
		<CardHeader class="text-center">
			<CardTitle class="text-2xl">SMS Dashboard</CardTitle>
			<CardDescription>SMSGate</CardDescription>
		</CardHeader>
		<form onsubmit={handleSubmit}>
			<CardContent class="space-y-4">
				{#if error}
					<Alert variant="destructive">
						<AlertDescription>{error}</AlertDescription>
					</Alert>
				{/if}
				<div class="space-y-2">
					<label for="login" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
						Login
					</label>
					<Input
						id="login"
						type="text"
						placeholder="Enter your login"
						autocomplete="username"
						bind:value={username}
						disabled={loading}
					/>
				</div>
				<div class="space-y-2">
					<label for="password" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
						Password
					</label>
					<Input
						id="password"
						type="password"
						placeholder="Enter your password"
						autocomplete="current-password"
						bind:value={password}
						disabled={loading}
					/>
				</div>
			</CardContent>
			<CardFooter>
				<Button type="submit" class="w-full" disabled={loading}>
					{loading ? 'Signing in...' : 'Sign In'}
				</Button>
			</CardFooter>
		</form>
	</Card>
</div>
