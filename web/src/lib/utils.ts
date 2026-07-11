import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function stateBadgeVariant(state: string): 'default' | 'secondary' | 'destructive' | 'outline' {
	switch (state) {
		case 'Delivered': return 'default';
		case 'Failed': return 'destructive';
		case 'Sent': return 'secondary';
		case 'Processed': return 'secondary';
		default: return 'outline';
	}
}
