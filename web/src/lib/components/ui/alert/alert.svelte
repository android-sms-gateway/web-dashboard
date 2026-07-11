<script lang="ts">
	import { cn } from '$lib/utils';
	import { cva, type VariantProps } from 'class-variance-authority';

	const alertVariants = cva(
		'relative w-full rounded-lg border p-4 [&>svg+div]:translate-y-[-3px] [&>svg]:absolute [&>svg]:left-4 [&>svg]:top-4 [&>svg]:text-foreground',
		{
			variants: {
				variant: {
					default: 'bg-background text-foreground',
					destructive:
						'border-destructive/50 text-destructive dark:border-destructive [&>svg]:text-destructive',
				},
			},
			defaultVariants: {
				variant: 'default',
			},
		},
	);

	interface Props extends VariantProps<typeof alertVariants> {
		variant?: 'default' | 'destructive';
		class?: string;
		children?: import('svelte').Snippet;
	}

	let {
		variant = 'default',
		class: className,
		children,
	}: Props = $props();
</script>

<div class={cn(alertVariants({ variant }), className)} role="alert">
	{@render children?.()}
</div>
