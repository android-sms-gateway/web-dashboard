<script lang="ts">
	import { page } from '$app/stores';
	import { Button } from '$lib/components/ui/button/index.js';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle,
	} from '$lib/components/ui/card/index.js';
</script>

<div class="mx-auto max-w-4xl space-y-6">
	<div>
		<h1 class="text-3xl font-bold tracking-tight">Dashboard</h1>
		<p class="text-muted-foreground">
			Welcome back, {$page.data.user?.login ?? 'User'}
		</p>
	</div>

	<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
		<Card>
			<CardHeader class="pb-2">
				<CardTitle class="text-sm font-medium text-muted-foreground">
					Devices Online
				</CardTitle>
			</CardHeader>
			<CardContent>
				<p class="text-2xl font-bold">—</p>
			</CardContent>
		</Card>
		<Card>
			<CardHeader class="pb-2">
				<CardTitle class="text-sm font-medium text-muted-foreground">
					Messages Sent
				</CardTitle>
			</CardHeader>
			<CardContent>
				<p class="text-2xl font-bold">—</p>
			</CardContent>
		</Card>
		<Card>
			<CardHeader class="pb-2">
				<CardTitle class="text-sm font-medium text-muted-foreground">
					Messages Pending
				</CardTitle>
			</CardHeader>
			<CardContent>
				<p class="text-2xl font-bold">—</p>
			</CardContent>
		</Card>
	</div>

	<Card>
		<CardHeader>
			<CardTitle>Getting Started</CardTitle>
			<CardDescription>
				This dashboard is under construction. More features coming soon.
			</CardDescription>
		</CardHeader>
		<CardContent>
			<p class="text-sm text-muted-foreground">
				Pages available in future phases: Messages, Devices, Webhooks, Settings,
				API Tokens.
			</p>
		</CardContent>
	</Card>
</div>
