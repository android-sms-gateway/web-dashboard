<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { sendMessage } from '$lib/api';
	import { dc, loadDevices } from '$lib/device-cache.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { Alert, AlertDescription } from '$lib/components/ui/alert/index.js';

	let phones = $state('');
	let text = $state('');
	let deviceId = $state('');
	let simNumber = $state('');
	let sending = $state(false);
	let error = $state('');

	onMount(() => loadDevices());

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';

		const phoneList = phones
			.split(/[\n,]+/)
			.map((p) => p.trim())
			.filter((p) => p.length > 0);

		if (phoneList.length === 0) {
			error = 'Enter at least one phone number';
			return;
		}

		if (!text.trim()) {
			error = 'Message text is required';
			return;
		}

		sending = true;
		try {
			await sendMessage({
				phoneNumbers: phoneList,
				text: text.trim(),
				...(deviceId ? { deviceId } : {}),
				...(simNumber.trim()
					? (() => { const n = parseInt(simNumber.trim(), 10); return Number.isFinite(n) ? { simNumber: Math.min(Math.max(n, 1), 3) } : {}; })()
					: {}),
			});
			goto('/messages');
		} catch {
			error = 'Failed to send message';
		} finally {
			sending = false;
		}
	}
</script>

<div class="mx-auto max-w-lg space-y-6">
	<div class="flex items-center gap-3">
		<Button variant="ghost" onclick={() => goto('/messages')}>&larr; Back</Button>
		<h1 class="text-2xl font-bold tracking-tight">Send Message</h1>
	</div>

	<Card>
		<CardHeader>
			<CardTitle>New Message</CardTitle>
		</CardHeader>
		<CardContent>
			<form onsubmit={handleSubmit} class="space-y-4">
				{#if error}
					<Alert variant="destructive">
						<AlertDescription>{error}</AlertDescription>
					</Alert>
				{/if}

				<div class="space-y-2">
					<Label for="phones">Phone Numbers (one per line or comma-separated)</Label>
					<Textarea
						id="phones"
						bind:value={phones}
						rows={4}
						placeholder={"+79990001234\n+79990005678"}
						disabled={sending}
					/>
				</div>

				<div class="space-y-2">
					<Label for="text">Message Text</Label>
					<Textarea
						id="text"
						bind:value={text}
						rows={5}
						placeholder="Enter your message..."
						disabled={sending}
					/>
				</div>

				<div class="space-y-2">
					<Label for="device">Device</Label>
					<select
						id="device"
						class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						value={deviceId}
						onchange={(e) => (deviceId = (e.target as HTMLSelectElement).value)}
						disabled={sending}
					>
						<option value="">Any device</option>
						{#if dc.loading}
							<option value="" disabled>Loading devices...</option>
						{:else if dc.error}
							<option value="" disabled>Failed to load</option>
						{:else}
							{#each dc.devices as d}
								<option value={d.id}>{d.name}</option>
							{/each}
						{/if}
					</select>
				</div>

				<div class="space-y-2">
					<Label for="sim">SIM Number (optional, 1-3)</Label>
					<Input
						id="sim"
						type="text"
						inputmode="numeric"
						bind:value={simNumber}
						placeholder="Leave empty for default"
						disabled={sending}
					/>
				</div>

				<Button type="submit" class="w-full" disabled={sending}>
					{sending ? 'Sending...' : 'Send Message'}
				</Button>
			</form>
		</CardContent>
	</Card>
</div>
