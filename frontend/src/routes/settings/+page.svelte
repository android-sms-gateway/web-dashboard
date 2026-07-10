<script lang="ts">
	import { onMount } from 'svelte';
	import { getSettings, updateSettings } from '$lib/api';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/components/ui/tabs/index.js';
	import { Alert, AlertDescription } from '$lib/components/ui/alert/index.js';
	import { toast } from 'svelte-sonner';
	import type { DeviceSettings } from '$lib/types';

	let activeTab = $state('messages');
	let loading = $state(true);
	let saving = $state(false);
	let error = $state('');

	function emptyForm(): DeviceSettings {
		return {
			messages: {},
			ping: {},
			logs: {},
			webhooks: {},
			gateway: {},
			encryption: {},
		};
	}

	let form: DeviceSettings = $state(emptyForm());

	onMount(load);

	async function load() {
		loading = true;
		error = '';
		try {
			const data = await getSettings();
			form = emptyForm();
			mergeNested(form, data);
		} catch {
			error = 'Failed to load settings';
		} finally {
			loading = false;
		}
	}

	function mergeNested(target: Record<string, unknown>, source: Record<string, unknown>) {
		for (const key of Object.keys(source)) {
			const val = source[key];
			if (val !== null && val !== undefined && typeof val === 'object' && !Array.isArray(val)) {
				if (!target[key] || typeof target[key] !== 'object') target[key] = {};
				mergeNested(target[key] as Record<string, unknown>, val as Record<string, unknown>);
			} else {
				target[key] = val;
			}
		}
	}

	async function handleSave() {
		error = '';
		saving = true;
		try {
			const result = await updateSettings(form);
			form = emptyForm();
			mergeNested(form, result);
			toast.success('Settings saved.');
		} catch {
			error = 'Failed to save settings';
		} finally {
			saving = false;
		}
	}
</script>

<div class="mx-auto max-w-2xl space-y-6">
	<div class="flex items-center justify-between">
		<h1 class="text-2xl font-bold tracking-tight">Settings</h1>
		<Button onclick={handleSave} disabled={loading || saving}>
			{saving ? 'Saving...' : 'Save'}
		</Button>
	</div>

	{#if error}
		<Alert variant="destructive">
			<AlertDescription>{error}</AlertDescription>
		</Alert>
	{/if}

	{#if loading}
		<div class="space-y-3">
			<Skeleton class="h-48 w-full rounded-lg" />
			<Skeleton class="h-24 w-full rounded-lg" />
		</div>
	{:else}
		<Tabs bind:value={activeTab}>
			<TabsList class="w-full justify-start">
				<TabsTrigger value="messages">Messages</TabsTrigger>
				<TabsTrigger value="ping">Ping</TabsTrigger>
				<TabsTrigger value="logs">Logs</TabsTrigger>
				<TabsTrigger value="webhooks">Webhooks</TabsTrigger>
				<TabsTrigger value="gateway">Gateway</TabsTrigger>
				<TabsTrigger value="encryption">Encryption</TabsTrigger>
			</TabsList>

			<TabsContent value="messages">
				<Card>
					<CardHeader><CardTitle>Messages</CardTitle></CardHeader>
					<CardContent class="grid gap-3">
						<div class="space-y-1">
							<Label for="si-min">Send Interval Min (seconds)</Label>
							<Input id="si-min" type="number" min="1" placeholder="Leave empty for default" bind:value={form.messages.sendIntervalMin} />
						</div>
						<div class="space-y-1">
							<Label for="si-max">Send Interval Max (seconds)</Label>
							<Input id="si-max" type="number" min="1" placeholder="Leave empty for default" bind:value={form.messages.sendIntervalMax} />
						</div>
						<div class="space-y-1">
							<Label for="limit-period">Limit Period</Label>
							<select id="limit-period" class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2" bind:value={form.messages.limitPeriod}>
								<option value={undefined}>Default</option>
								<option value="Disabled">Disabled</option>
								<option value="PerMinute">Per Minute</option>
								<option value="PerHour">Per Hour</option>
								<option value="PerDay">Per Day</option>
							</select>
						</div>
						<div class="space-y-1">
							<Label for="limit-value">Limit Value</Label>
							<Input id="limit-value" type="number" min="1" placeholder="Leave empty for default" bind:value={form.messages.limitValue} />
						</div>
						<div class="space-y-1">
							<Label for="sim-mode">SIM Selection Mode</Label>
							<select id="sim-mode" class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2" bind:value={form.messages.simSelectionMode}>
								<option value={undefined}>Default</option>
								<option value="OSDefault">OS Default</option>
								<option value="RoundRobin">Round Robin</option>
								<option value="Random">Random</option>
							</select>
						</div>
						<div class="space-y-1">
							<Label for="log-lifetime">Log Lifetime (days)</Label>
							<Input id="log-lifetime" type="number" min="1" placeholder="Leave empty for default" bind:value={form.messages.logLifetimeDays} />
						</div>
						<div class="space-y-1">
							<Label for="processing-order">Processing Order</Label>
							<select id="processing-order" class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2" bind:value={form.messages.processingOrder}>
								<option value={undefined}>Default</option>
								<option value="LIFO">LIFO (last in, first out)</option>
								<option value="FIFO">FIFO (first in, first out)</option>
							</select>
						</div>
					</CardContent>
				</Card>
			</TabsContent>

			<TabsContent value="ping">
				<Card>
					<CardHeader><CardTitle>Ping</CardTitle></CardHeader>
					<CardContent class="grid gap-3">
						<div class="space-y-1">
							<Label for="ping-int">Interval (seconds)</Label>
							<Input id="ping-int" type="number" min="1" placeholder="Leave empty for default" bind:value={form.ping.intervalSeconds} />
						</div>
					</CardContent>
				</Card>
			</TabsContent>

			<TabsContent value="logs">
				<Card>
					<CardHeader><CardTitle>Logs</CardTitle></CardHeader>
					<CardContent class="grid gap-3">
						<div class="space-y-1">
							<Label for="logs-lt">Lifetime (days)</Label>
							<Input id="logs-lt" type="number" min="1" placeholder="Leave empty for default" bind:value={form.logs.lifetimeDays} />
						</div>
					</CardContent>
				</Card>
			</TabsContent>

			<TabsContent value="webhooks">
				<Card>
					<CardHeader><CardTitle>Webhooks</CardTitle></CardHeader>
					<CardContent class="grid gap-3">
						<label class="flex items-center gap-2 text-sm">
							<input type="checkbox" class="h-4 w-4 rounded border-input" bind:checked={form.webhooks.internetRequired} />
							Internet Required
						</label>
						<div class="space-y-1">
							<Label for="wh-retry">Retry Count</Label>
							<Input id="wh-retry" type="number" min="1" placeholder="Leave empty for default" bind:value={form.webhooks.retryCount} />
						</div>
						<div class="space-y-1">
							<Label for="wh-signing">Signing Key</Label>
							<Input id="wh-signing" type="password" placeholder="Leave empty to keep current" bind:value={form.webhooks.signingKey} />
						</div>
					</CardContent>
				</Card>
			</TabsContent>

			<TabsContent value="gateway">
				<Card>
					<CardHeader><CardTitle>Gateway</CardTitle></CardHeader>
					<CardContent class="grid gap-3">
						<div class="space-y-1">
							<Label for="gw-url">Cloud URL</Label>
							<Input id="gw-url" type="url" placeholder="Leave empty for default" bind:value={form.gateway.cloudUrl} />
						</div>
						<div class="space-y-1">
							<Label for="gw-token">Private Token</Label>
							<Input id="gw-token" type="password" placeholder="Leave empty to keep current" bind:value={form.gateway.privateToken} />
						</div>
						<div class="space-y-1">
							<Label for="gw-notif">Notification Channel</Label>
							<select id="gw-notif" class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2" bind:value={form.gateway.notificationChannel}>
								<option value={undefined}>Default</option>
								<option value="AUTO">Auto</option>
								<option value="SSE_ONLY">SSE Only</option>
							</select>
						</div>
					</CardContent>
				</Card>
			</TabsContent>

			<TabsContent value="encryption">
				<Card>
					<CardHeader><CardTitle>Encryption</CardTitle></CardHeader>
					<CardContent class="grid gap-3">
						<div class="space-y-1">
							<Label for="enc-pass">Passphrase</Label>
							<Input id="enc-pass" type="password" placeholder="Leave empty to keep current" bind:value={form.encryption.passphrase} />
						</div>
					</CardContent>
				</Card>
			</TabsContent>
		</Tabs>

		<div class="pb-6">
			<Button onclick={handleSave} disabled={saving} class="w-full sm:w-auto">
				{saving ? 'Saving...' : 'Save Settings'}
			</Button>
		</div>
	{/if}
</div>
