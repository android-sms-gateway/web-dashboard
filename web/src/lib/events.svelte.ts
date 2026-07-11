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

export const events = $state({
	status: 'disconnected' as 'disconnected' | 'connecting' | 'connected',
	error: null as string | null,
});

let source: EventSource | null = null;
const handlers = new Map<EventType, Set<Handler>>();

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
		if (source?.readyState === EventSource.CLOSED) {
			events.status = 'disconnected';
			events.error = 'connection lost';
		}
	};

	source.addEventListener('message', (e) => {
		try {
			const parsed = JSON.parse((e as MessageEvent<string>).data);
			if (!isValidSseEvent(parsed)) {
				console.warn('Invalid SSE event structure', parsed);
				return;
			}
			dispatch(parsed);
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

function dispatch(event: SseEvent): void {
	const set = handlers.get(event.type);
	if (!set) return;

	for (const handler of set) {
		try {
			handler(event);
		} catch (err) {
			console.warn('SSE handler error', err);
		}
	}
}
