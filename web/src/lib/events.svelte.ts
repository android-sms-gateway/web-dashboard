import type { SseEvent } from './types';

const VALID_EVENT_TYPES = [
	'message.received',
	'message.state_changed',
	'device.status_changed',
	'stats.updated',
] as const;

function isValidSseEvent(data: unknown): data is SseEvent {
	if (typeof data !== 'object' || data === null) return false;
	const d = data as Record<string, unknown>;
	if (typeof d.type !== 'string') return false;
	if (!(VALID_EVENT_TYPES as readonly string[]).includes(d.type)) return false;
	if (typeof d.payload !== 'object' || d.payload === null) return false;
	return true;
}

type EventType = SseEvent['type'];
type Handler = (event: SseEvent) => void;

export interface Toast {
	id: number;
	type: string;
	message: string;
	createdAt: number;
}

export const events = $state({
	status: 'disconnected' as 'disconnected' | 'connecting' | 'connected',
	error: null as string | null,
	toasts: [] as Toast[],
});

let source: EventSource | null = null;
const handlers = new Map<EventType, Set<Handler>>();
let nextId = 1;
const TOAST_TTL_MS = 5000;

export function connect(): void {
	if (source) return;

	events.status = 'connecting';
	events.error = null;

	source = new EventSource('/api/v1/events', { withCredentials: true });

	source.onopen = () => {
		events.status = 'connected';
		events.error = null;
	};

	source.onerror = () => {
		events.status = 'disconnected';
		events.error = 'connection lost';
	};

	source.addEventListener('message', (e) => {
		try {
			const parsed = JSON.parse((e as MessageEvent<string>).data);
			if (!isValidSseEvent(parsed)) {
				console.warn('Invalid SSE event structure', parsed);
				return;
			}
			handleIncoming(parsed);
		} catch (err) {
			console.warn('Failed to parse SSE event', err);
		}
	});
}

export function disconnect(): void {
	if (!source) return;

	source.close();
	source = null;
	events.status = 'disconnected';
	events.error = null;
	events.toasts = [];
	handlers.clear();
}

export function on(type: EventType, handler: Handler): () => void {
	let set = handlers.get(type);
	if (!set) {
		set = new Set();
		handlers.set(type, set);
	}
	set.add(handler);

	return () => {
		set!.delete(handler);
		if (set!.size === 0) handlers.delete(type);
	};
}

export function dismissToast(id: number): void {
	events.toasts = events.toasts.filter((t) => t.id !== id);
}

export function showToast(type: string, message: string): void {
	const toast: Toast = {
		id: nextId++,
		type,
		message,
		createdAt: Date.now(),
	};
	events.toasts = [toast, ...events.toasts];
	setTimeout(() => dismissToast(toast.id), TOAST_TTL_MS);
}

function handleIncoming(event: SseEvent): void {
	dispatch(event);
	addToast(event);
}

function dispatch(event: SseEvent): void {
	const set = handlers.get(event.type);
	if (!set) return;

	for (const handler of set) {
		try {
			handler(event);
		} catch (err) {
			console.error('SSE handler error', err);
		}
	}
}

function addToast(event: SseEvent): void {
	if (event.type === 'stats.updated') return;

	const message = formatToastMessage(event);
	const toast: Toast = {
		id: nextId++,
		type: event.type,
		message,
		createdAt: Date.now(),
	};

	events.toasts = [toast, ...events.toasts];

	setTimeout(() => dismissToast(toast.id), TOAST_TTL_MS);
}

function formatToastMessage(event: SseEvent): string {
	switch (event.type) {
		case 'message.received':
			return `SMS from ${event.payload.sender}: ${event.payload.message.slice(0, 50)}`;
		case 'message.state_changed':
			return `Message ${event.payload.messageId} → ${event.payload.state}`;
		case 'device.status_changed':
			return `Device ${event.payload.deviceId} ${event.payload.isOnline ? 'online' : 'offline'}`;
		case 'stats.updated':
			return 'Stats updated';
	}
}
