<script lang="ts">
	import type { AttachmentDraft } from "$lib/types";
	import { formatSize } from "$lib/utils";
	import { Button } from "$lib/components/ui/button/index.js";
	import { File as FileIcon, X } from "lucide-svelte";

	interface Props {
		draft: AttachmentDraft;
		disabled?: boolean;
		onremove: () => void;
	}

	let { draft, disabled = false, onremove }: Props = $props();

	let isImage = $derived(draft.contentType.startsWith("image/"));
	let ext = $derived(
		draft.name.split(".").pop()?.toUpperCase().slice(0, 3) ?? "",
	);
</script>

<li
	class="flex items-center gap-3 rounded-md border bg-card p-2 {draft.status ===
	'error'
		? 'border-destructive/50'
		: ''}"
>
	{#if isImage && draft.objectUrl}
		<img
			src={draft.objectUrl}
			alt=""
			class="h-10 w-10 shrink-0 rounded object-cover"
		/>
	{:else}
		<span
			class="flex h-10 w-10 shrink-0 items-center justify-center rounded bg-muted"
		>
			{#if ext && ext !== "FILE"}
				<span class="font-mono text-xs text-muted-foreground"
					>{ext}</span
				>
			{:else}
				<FileIcon class="h-4 w-4 text-muted-foreground" />
			{/if}
		</span>
	{/if}

	<div class="min-w-0 flex-1">
		<p class="truncate text-sm font-medium">{draft.name}</p>
		{#if draft.status === "error"}
			<p class="text-xs text-destructive">Could not read file</p>
		{:else}
			<p class="truncate text-xs text-muted-foreground">
				{draft.contentType} - {formatSize(draft.size)}
				{#if draft.size === 0}
					(empty file)
				{/if}
				{#if draft.status === "pending"}
					(reading...)
				{/if}
			</p>
		{/if}
	</div>

	<Button
		variant="ghost"
		size="icon"
		class="h-11 w-11 shrink-0"
		aria-label={`Remove ${draft.name}`}
		{disabled}
		onclick={onremove}
	>
		<X class="h-4 w-4" />
	</Button>
</li>
