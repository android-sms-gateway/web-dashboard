<script lang="ts">
	import { onDestroy } from "svelte";
	import type { AttachmentDraft } from "$lib/types";
	import { cn, formatSize } from "$lib/utils";
	import { Label } from "$lib/components/ui/label/index.js";
	import { Alert, AlertDescription } from "$lib/components/ui/alert/index.js";
	import { TriangleAlert } from "lucide-svelte";
	import FileCard from "./FileCard.svelte";

	interface Props {
		attachments: AttachmentDraft[];
		disabled?: boolean;
	}

	let { attachments = $bindable(), disabled = false }: Props = $props();

	let dragging = $state(false);
	let inputEl: HTMLInputElement;
	let readers = new Map<AttachmentDraft, FileReader>();

	let totalBytes = $derived(attachments.reduce((sum, a) => sum + a.size, 0));
	let fileCount = $derived(attachments.length);

	onDestroy(() => {
		for (const [draft, reader] of readers) {
			reader.abort();
		}
		readers.clear();
		for (const a of attachments) {
			if (a.objectUrl) URL.revokeObjectURL(a.objectUrl);
		}
	});

	function handleFiles(files: Iterable<File>) {
		const seen = new Set(
			attachments.map(
				(a) => `${a.file.name}|${a.file.size}|${a.file.lastModified}`,
			),
		);
		const fresh: AttachmentDraft[] = [];
		for (const file of files) {
			const key = `${file.name}|${file.size}|${file.lastModified}`;
			if (seen.has(key)) continue;
			seen.add(key);
			fresh.push({
				file,
				name: file.name,
				contentType: file.type || "application/octet-stream",
				size: file.size,
				data: "",
				status: "pending",
			});
		}
		if (fresh.length === 0) return;
		attachments.push(...fresh);
		const start = attachments.length - fresh.length;
		for (let i = 0; i < fresh.length; i += 1) readFile(attachments[start + i]);
	}

	function readFile(draft: AttachmentDraft) {
		const reader = new FileReader();
		readers.set(draft, reader);
		reader.onload = () => {
			const result = reader.result;
			if (typeof result !== "string") {
				draft.status = "error";
				return;
			}
			const comma = result.indexOf(",");
			draft.data = comma >= 0 ? result.slice(comma + 1) : result;
			if (draft.contentType.startsWith("image/")) {
				try {
					draft.objectUrl = URL.createObjectURL(draft.file);
				} catch {
					// object URL unsupported - icon tile is used instead
				}
			}
			draft.status = "ready";
			readers.delete(draft);
		};
		reader.onerror = () => {
			draft.status = "error";
			readers.delete(draft);
		};
		reader.readAsDataURL(draft.file);
	}

	function removeDraft(draft: AttachmentDraft) {
		const i = attachments.indexOf(draft);
		if (i === -1) return;
		const reader = readers.get(draft);
		if (reader) {
			reader.abort();
			readers.delete(draft);
		}
		if (draft.objectUrl) URL.revokeObjectURL(draft.objectUrl);
		attachments.splice(i, 1);
	}
</script>

<div class="space-y-2">
	<Label
		>Attachments <span class="text-muted-foreground">(optional)</span
		></Label
	>

	<button
		type="button"
		class={cn(
			"w-full rounded-md border border-dashed p-6 text-center text-sm text-muted-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
			dragging
				? "border-primary bg-accent"
				: "border-input bg-muted/50 hover:bg-muted",
		)}
		aria-label="Add attachments"
		aria-describedby="attachments-hint"
		{disabled}
		onclick={() => inputEl?.click()}
		ondragover={(e) => {
			e.preventDefault();
			dragging = true;
		}}
		ondragenter={(e) => {
			e.preventDefault();
			dragging = true;
		}}
		ondragleave={() => {
			dragging = false;
		}}
		ondrop={(e) => {
			e.preventDefault();
			dragging = false;
			if (disabled) return;
			if (e.dataTransfer) handleFiles(e.dataTransfer.files);
		}}
	>
		Drag and drop files here, or <span class="underline">browse</span>
	</button>

	<input
		class="sr-only"
		type="file"
		multiple
		aria-hidden="true"
		tabindex="-1"
		bind:this={inputEl}
		onchange={(e) => {
			const input = e.target as HTMLInputElement;
			if (input.files) handleFiles(input.files);
			input.value = "";
		}}
	/>

	<p id="attachments-hint" class="text-xs text-muted-foreground">
		Adding at least one attachment or message text is required
	</p>

	{#if fileCount > 0}
		<ul
			class="space-y-2 {fileCount >= 5
				? 'sm:grid sm:grid-cols-2 sm:gap-2 sm:space-y-0'
				: ''}"
		>
			{#each attachments as draft (draft.file.name + draft.file.size + draft.file.lastModified)}
				<FileCard
					{draft}
					{disabled}
					onremove={() => removeDraft(draft)}
				/>
			{/each}
		</ul>

		<p
			class="flex items-center gap-1 text-xs {totalBytes >
			10 * 1024 * 1024
				? 'text-amber-700 dark:text-amber-400'
				: 'text-muted-foreground'}"
		>
			{#if totalBytes > 10 * 1024 * 1024}
				<TriangleAlert class="h-3.5 w-3.5 shrink-0" />
			{/if}
			{fileCount}
			{fileCount === 1 ? "file" : "files"}, {formatSize(totalBytes)} total
		</p>

		{#if totalBytes > 25 * 1024 * 1024}
			<Alert>
				<TriangleAlert class="h-4 w-4" />
				<AlertDescription>
					Large attachments may fail to send or be rejected by
					carriers.
				</AlertDescription>
			</Alert>
		{/if}
	{/if}
</div>
