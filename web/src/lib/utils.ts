import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

/** Merges class names using clsx and tailwind-merge for conflict-free Tailwind CSS classes. */
export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

/** Formats a byte count into a human-readable string (B, KB, or MB). */
export function formatSize(bytes: number): string {
	if (bytes < 1024) return `${bytes} B`;
	if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
	return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

/** Returns the badge variant style for a given processing state string. */
export function stateBadgeVariant(state: string): 'default' | 'secondary' | 'destructive' | 'outline' {
	switch (state) {
		case 'Delivered': return 'default';
		case 'Failed': return 'destructive';
		case 'Sent': return 'secondary';
		case 'Processed': return 'secondary';
		default: return 'outline';
	}
}
