<script lang="ts">
	import type { OutgoingAttachment } from "$lib/types";
	import { formatSize } from "$lib/utils";
	import { File as FileIcon } from "lucide-svelte";

	interface Props {
		attachment: OutgoingAttachment;
	}

	let { attachment }: Props = $props();

	let broken = $state(false);
	let isImage = $derived(
		attachment.data != null &&
			attachment.contentType.toLowerCase().startsWith("image/"),
	);
	let ext = $derived(
		attachment.name?.split(".").pop()?.toUpperCase().slice(0, 3) ?? "",
	);
	let sizeBytes = $derived(
		attachment.data
			? Math.floor(
					(attachment.data.length * 3) / 4 -
						(attachment.data.endsWith("==")
							? 2
							: attachment.data.endsWith("=")
								? 1
								: 0),
				)
			: undefined,
	);
</script>

{#if isImage && !broken}
	<figure class="space-y-2 rounded-md border bg-card p-3">
		<img
			src="data:{attachment.contentType};base64,{attachment.data}"
			alt={attachment.name ?? ""}
			class="h-40 w-full rounded bg-muted object-contain"
			loading="lazy"
			onerror={() => (broken = true)}
		/>
		<figcaption class="break-all text-sm font-medium">
			{attachment.name ?? "attachment"}
		</figcaption>
	</figure>
{:else}
	<div class="space-y-2 rounded-md border bg-card p-3">
		<div
			class="flex h-10 w-10 items-center justify-center rounded bg-muted"
		>
			{#if ext && ext !== "FILE"}
				<span class="font-mono text-xs text-muted-foreground"
					>{ext}</span
				>
			{:else}
				<FileIcon class="h-4 w-4 text-muted-foreground" />
			{/if}
		</div>
		<p class="break-all text-sm font-medium">
			{attachment.name ?? "attachment"}
		</p>
		<p class="text-xs text-muted-foreground">
			{attachment.contentType}{sizeBytes != null
				? ` - ${formatSize(sizeBytes)}`
				: ""}
		</p>
	</div>
{/if}
