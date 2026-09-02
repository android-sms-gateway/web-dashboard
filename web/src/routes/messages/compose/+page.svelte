<script lang="ts">
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";
	import { sendMessage } from "$lib/api";
	import { dc, loadDevices } from "$lib/device-cache.svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import { Textarea } from "$lib/components/ui/textarea/index.js";
	import {
		Card,
		CardContent,
		CardHeader,
		CardTitle,
	} from "$lib/components/ui/card/index.js";
	import { Alert, AlertDescription } from "$lib/components/ui/alert/index.js";
	import {
		Tabs,
		TabsList,
		TabsTrigger,
	} from "$lib/components/ui/tabs/index.js";
	import AttachmentUploader from "$lib/components/messages/AttachmentUploader.svelte";
	import type { AttachmentDraft, SendMessageRequest } from "$lib/types";

	let phones = $state("");
	let text = $state("");
	let subject = $state("");
	let deviceId = $state("");
	let simNumber = $state("");
	let messageType = $state<"sms" | "mms">("sms");
	let attachments = $state<AttachmentDraft[]>([]);
	let sending = $state(false);
	let error = $state("");
	let attempted = $state(false);

	function parseSimNumber(): { error?: string; simNumber?: number } {
		const raw = simNumber.trim();
		if (!raw) return {};
		const n = parseInt(raw, 10);
		if (!Number.isFinite(n) || String(n) !== raw || n < 1 || n > 3) {
			return { error: "SIM number must be 1, 2, or 3" };
		}
		return { simNumber: n };
	}

	onMount(() => loadDevices());

	let readyAttachments = $derived(
		attachments.filter((a) => a.status === "ready"),
	);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = "";
		attempted = true;

		const phoneList = phones
			.split(/[\n,]+/)
			.map((p) => p.trim())
			.filter((p) => p.length > 0);

		if (phoneList.length === 0) {
			error = "Enter at least one phone number";
			return;
		}

		if (messageType === "sms") {
			if (!text.trim()) {
				error = "Message text is required";
				return;
			}
		} else if (!text.trim() && readyAttachments.length === 0) {
			error = "Add message text or at least one attachment";
			return;
		}

		const { error: simError, ...simPayload } = parseSimNumber();
		if (simError) {
			error = simError;
			return;
		}

		sending = true;
		try {
			const payload: SendMessageRequest =
				messageType === "mms"
					? {
							phoneNumbers: phoneList,
							...(deviceId ? { deviceId } : {}),
							...simPayload,
							mmsMessage: {
								...(subject.trim()
									? { subject: subject.trim() }
									: {}),
								...(text.trim() ? { text: text.trim() } : {}),
								...(readyAttachments.length > 0
									? {
											attachments: readyAttachments.map(
												(a) => ({
													contentType: a.contentType,
													name: a.name,
													data: a.data,
												}),
											),
										}
									: {}),
							},
						}
					: {
							phoneNumbers: phoneList,
							text: text.trim(),
							...(deviceId ? { deviceId } : {}),
							...simPayload,
						};
			await sendMessage(payload);
			goto("/messages");
		} catch (e) {
			error =
				e instanceof Error && e.message
					? e.message
					: "Failed to send message";
		} finally {
			sending = false;
		}
	}
</script>

<div class="mx-auto max-w-lg space-y-6">
	<div class="flex items-center gap-3">
		<Button variant="ghost" onclick={() => goto("/messages")}
			>&larr; Back</Button
		>
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
					<Label for="phones"
						>Phone Numbers (one per line or comma-separated)</Label
					>
					<Textarea
						id="phones"
						bind:value={phones}
						rows={4}
						placeholder={"+79990001234\n+79990005678"}
						disabled={sending}
					/>
				</div>

				<div class="space-y-2">
					<span
						id="message-type-label"
						class="text-sm font-medium leading-none"
						>Message Type</span
					>
					<Tabs bind:value={messageType}>
						<TabsList
							aria-labelledby="message-type-label"
							class="grid w-full grid-cols-2"
						>
							<TabsTrigger value="sms">SMS</TabsTrigger>
							<TabsTrigger value="mms">MMS</TabsTrigger>
						</TabsList>
					</Tabs>
				</div>

				{#if messageType === "mms"}
					<div class="space-y-2">
						<Label for="subject"
							>Subject <span class="text-muted-foreground"
								>(optional)</span
							></Label
						>
						<Input
							id="subject"
							type="text"
							bind:value={subject}
							placeholder="Enter subject..."
							disabled={sending}
						/>
						{#if subject.length > 100}
							<p class="text-xs text-muted-foreground">
								{subject.length} characters
							</p>
						{/if}
					</div>
				{/if}

				{#if messageType === "mms"}
					<AttachmentUploader bind:attachments disabled={sending} />
				{/if}

				<div class="space-y-2">
					<Label for="text">
						Message Text
						{#if messageType === "mms"}
							<span class="text-muted-foreground"
								>(optional if attachments added)</span
							>
						{/if}
					</Label>
					<Textarea
						id="text"
						bind:value={text}
						rows={5}
						placeholder="Enter your message..."
						disabled={sending}
					/>
					{#if messageType === "mms" && attempted && !text.trim() && readyAttachments.length === 0}
						<p class="text-sm text-destructive">
							Add message text or at least one attachment.
						</p>
					{/if}
				</div>

				<div class="space-y-2">
					<Label for="device">Device</Label>
					<select
						id="device"
						class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						value={deviceId}
						onchange={(e) =>
							(deviceId = (e.target as HTMLSelectElement).value)}
						disabled={sending}
					>
						<option value="">Any device</option>
						{#if dc.loading}
							<option value="" disabled>Loading devices...</option
							>
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
					{sending ? "Sending..." : "Send Message"}
				</Button>
			</form>
		</CardContent>
	</Card>
</div>
