<script lang="ts">
	import { createToken, revokeToken } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Alert, AlertDescription } from '$lib/components/ui/alert/index.js';
	import Dialog from '$lib/components/ui/dialog/dialog.svelte';
	import { toast } from 'svelte-sonner';
	import type { CreateTokenResponse } from '$lib/types';

	type View = 'create' | 'result';
	let view = $state<View>('create');
	let ttlInput = $state('');
	let selectedScopes = $state<string[]>([]);
	let saving = $state(false);
	let error = $state('');
	let result = $state<CreateTokenResponse | null>(null);
	let copied = $state(false);

	let revokeJti = $state('');
	let revoking = $state(false);
	let revokeError = $state('');
	let confirmRevoke = $state(false);

	const scopeList = [
		{ value: 'devices:list', label: 'List Devices' },
		{ value: 'devices:delete', label: 'Delete Devices' },
		{ value: 'inbox:list', label: 'List Inbox Messages' },
		{ value: 'inbox:refresh', label: 'Refresh Inbox' },
		{ value: 'logs:read', label: 'Read Logs' },
		{ value: 'messages:send', label: 'Send Messages' },
		{ value: 'messages:read', label: 'Read Messages' },
		{ value: 'messages:list', label: 'List Messages' },
		{ value: 'messages:export', label: 'Export Messages' },
		{ value: 'settings:read', label: 'Read Settings' },
		{ value: 'settings:write', label: 'Write Settings' },
		{ value: 'tokens:manage', label: 'Manage Tokens' },
		{ value: 'webhooks:list', label: 'List Webhooks' },
		{ value: 'webhooks:write', label: 'Write Webhooks' },
		{ value: 'webhooks:delete', label: 'Delete Webhooks' },
	];

	function toggleScope(scope: string) {
		if (selectedScopes.includes(scope)) {
			selectedScopes = selectedScopes.filter((s) => s !== scope);
		} else {
			selectedScopes = [...selectedScopes, scope];
		}
	}

	function resetForm() {
		view = 'create';
		ttlInput = '';
		selectedScopes = [];
		error = '';
		result = null;
		copied = false;
	}

	async function handleGenerate() {
		error = '';
		if (selectedScopes.length === 0) { error = 'Select at least one scope'; return; }
		const ttlValue = ttlInput.trim();
		let ttl: number | undefined;
		if (ttlValue) {
			const parsed = parseInt(ttlValue, 10);
			if (Number.isNaN(parsed) || parsed <= 0) {
				error = 'TTL must be a positive whole number';
				return;
			}
			ttl = parsed;
		}
		saving = true;
		try {
			result = await createToken({
				scopes: selectedScopes,
				ttl,
			});
			view = 'result';
		} catch {
			error = 'Failed to generate token';
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
		revokeError = '';
		if (!revokeJti.trim()) { revokeError = 'Enter a token JTI'; return; }
		confirmRevoke = true;
	}

	async function confirmRevokeToken() {
		confirmRevoke = false;
		revoking = true;
		try {
			await revokeToken(revokeJti.trim());
			toast.success('Token revoked.');
			revokeJti = '';
		} catch {
			revokeError = 'Failed to revoke token';
		} finally {
			revoking = false;
		}
	}
</script>

<div class="mx-auto max-w-2xl space-y-6">
	{#if view === 'result' && result}
		<h1 class="text-2xl font-bold tracking-tight">Token Generated</h1>

		<Card>
			<CardContent class="space-y-4 pt-6">
				<div class="rounded-md bg-yellow-50 p-3 text-center text-sm font-medium text-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-400">
					Save this token now. It will not be shown again.
				</div>

				<div class="space-y-2">
					<Label>Access Token</Label>
					<div class="flex gap-2">
						<code class="flex-1 rounded-md border bg-muted p-3 font-mono text-xs break-all">
							{result.access_token}
						</code>
						<Button variant="outline" onclick={copyToken}>
							{copied ? 'Copied!' : 'Copy'}
						</Button>
					</div>
				</div>

				<div class="grid grid-cols-[100px_1fr] gap-2 text-sm">
					<span class="font-medium text-muted-foreground">ID (JTI)</span>
					<span class="font-mono text-xs">{result.id}</span>

					<span class="font-medium text-muted-foreground">Type</span>
					<span>{result.token_type}</span>

					<span class="font-medium text-muted-foreground">Expires At</span>
					<span>{new Date(result.expires_at).toLocaleString()}</span>

					{#if result.refresh_token}
						<span class="font-medium text-muted-foreground">Refresh Token</span>
						<span class="font-mono text-xs break-all">{result.refresh_token}</span>
					{/if}
				</div>

				<Button class="w-full" onclick={resetForm}>Done</Button>
			</CardContent>
		</Card>
	{:else}
		<h1 class="text-2xl font-bold tracking-tight">Generate API Token</h1>

		{#if error}
			<Alert variant="destructive">
				<AlertDescription>{error}</AlertDescription>
			</Alert>
		{/if}

		<Card>
			<CardContent class="space-y-4 pt-6">
				<form onsubmit={(e) => { e.preventDefault(); handleGenerate(); }}>
					<div class="space-y-2">
						<Label for="ttl">TTL (seconds, optional)</Label>
						<Input id="ttl" type="number" min="1" bind:value={ttlInput} placeholder="Leave empty for default" disabled={saving} />
					</div>

					<div class="space-y-2">
						<Label>Scopes</Label>
						<div class="grid grid-cols-1 gap-1 sm:grid-cols-2">
							{#each scopeList as scope}
								<label class="flex items-center gap-2 text-sm cursor-pointer py-1">
									<input
										type="checkbox"
										class="h-4 w-4 rounded border-input"
										checked={selectedScopes.includes(scope.value)}
										onchange={() => toggleScope(scope.value)}
										disabled={saving}
									/>
									{scope.label}
								</label>
							{/each}
						</div>
					</div>

					<Button type="submit" class="w-full" disabled={saving || selectedScopes.length === 0}>
						{saving ? 'Generating...' : 'Generate Token'}
					</Button>
				</form>
			</CardContent>
		</Card>

		<Card>
			<details class="group">
				<summary class="flex cursor-pointer items-center gap-2 p-4 text-sm font-semibold text-muted-foreground hover:bg-accent [&::-webkit-details-marker]:hidden">
					Revoke an existing token
				</summary>
				<CardContent class="space-y-3 pt-0">
					{#if revokeError}
						<Alert variant="destructive">
							<AlertDescription>{revokeError}</AlertDescription>
						</Alert>
					{/if}
					<div class="space-y-2">
						<Label for="revoke-jti">Token JTI</Label>
						<Input id="revoke-jti" type="text" bind:value={revokeJti} placeholder="Paste the token JTI" disabled={revoking} />
					</div>
					<Button variant="destructive" onclick={handleRevoke} disabled={revoking}>
						{revoking ? 'Revoking...' : 'Revoke'}
					</Button>
				</CardContent>
			</details>
		</Card>
	{/if}
</div>

<Dialog
	open={confirmRevoke}
	title="Revoke Token"
	message="Are you sure you want to revoke this token? This action cannot be undone."
	confirmLabel="Revoke"
	variant="destructive"
	onConfirm={confirmRevokeToken}
	onCancel={() => (confirmRevoke = false)}
/>
