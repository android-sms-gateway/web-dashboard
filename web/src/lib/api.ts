import type {
	CreateTokenRequest,
	CreateTokenResponse,
	CreateWebhookRequest,
	Device,
	DeviceSettings,
	ListMessagesParams,
	ListMessagesResponse,
	LoginRequest,
	LoginResponse,
	Me,
	MessageDetail,
	SendMessageRequest,
	Stats,
	TrendsResponse,
	Webhook,
} from './types';

const BASE = '/api/v1';

// The backend speaks snake_case (Go json tags) while the UI binds camelCase
// fields from DeviceSettings (types.ts). These tables map between the two
// shapes; unknown keys are dropped so raw snake_case never leaks into form
// state.
const CAMEL_TO_SNAKE: Record<string, Record<string, string>> = {
	encryption: { passphrase: 'passphrase' },
	messages: {
		sendIntervalMin: 'send_interval_min',
		sendIntervalMax: 'send_interval_max',
		limitPeriod: 'limit_period',
		limitValue: 'limit_value',
		simSelectionMode: 'sim_selection_mode',
		logLifetimeDays: 'log_lifetime_days',
		processingOrder: 'processing_order'
	},
	ping: { intervalSeconds: 'interval_seconds' },
	logs: { lifetimeDays: 'lifetime_days' },
	webhooks: {
		internetRequired: 'internet_required',
		retryCount: 'retry_count',
		signingKey: 'signing_key'
	},
	gateway: {
		cloudUrl: 'cloud_url',
		privateToken: 'private_token',
		notificationChannel: 'notification_channel'
	}
};

function invertSections(map: Record<string, Record<string, string>>): Record<string, Record<string, string>> {
	const result: Record<string, Record<string, string>> = {};
	for (const [section, keys] of Object.entries(map)) {
		result[section] = {};
		for (const [from, to] of Object.entries(keys)) {
			result[section][to] = from;
		}
	}
	return result;
}

const SNAKE_TO_CAMEL = invertSections(CAMEL_TO_SNAKE);

function isMissing(value: unknown): boolean {
	return value === undefined || value === null || value === '';
}

function mapSettings(
	source: Record<string, unknown>,
	keyMap: Record<string, Record<string, string>>
): Record<string, Record<string, unknown>> {
	const result: Record<string, Record<string, unknown>> = {};
	for (const [section, keys] of Object.entries(keyMap)) {
		const raw = source[section];
		if (raw === null || typeof raw !== 'object' || Array.isArray(raw)) continue;
		const mapped: Record<string, unknown> = {};
		for (const [fromKey, toKey] of Object.entries(keys)) {
			const value = (raw as Record<string, unknown>)[fromKey];
			if (!isMissing(value)) mapped[toKey] = value;
		}
		if (Object.keys(mapped).length > 0) result[section] = mapped;
	}
	return result;
}

/** Maps a snake_case server payload to the camelCase DeviceSettings shape. */
export function fromServer(raw: Record<string, unknown>): DeviceSettings {
	return mapSettings(raw, SNAKE_TO_CAMEL) as DeviceSettings;
}

/** Maps camelCase form state back to the snake_case server payload shape. */
export function toServer(settings: DeviceSettings): Record<string, Record<string, unknown>> {
	return mapSettings(settings as Record<string, unknown>, CAMEL_TO_SNAKE);
}

const MAX_ERROR_DETAIL_LENGTH = 300;

export function formatApiError(status: number, detail: string): string {
	const normalized = detail.replace(/\s+/g, ' ').trim();
	if (!normalized) return `API error ${status}`;
	const short =
		normalized.length > MAX_ERROR_DETAIL_LENGTH
			? `${normalized.slice(0, MAX_ERROR_DETAIL_LENGTH)}...`
			: normalized;
	return `API error ${status}: ${short}`;
}

export function errorMessage(error: unknown): string {
	if (error instanceof Error && error.message) return error.message;
	if (typeof error === 'string' && error) return error;
	return 'Unknown error';
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
	const res = await fetch(`${BASE}${path}`, {
		...options,
		credentials: options?.credentials ?? 'include',
		headers: {
			'Content-Type': 'application/json',
			...options?.headers,
		},
	});

	if (!res.ok) {
		let detail = '';
		try {
			detail = await res.text();
		} catch {}
		throw new Error(formatApiError(res.status, detail));
	}

	if (res.status === 204) {
		return undefined as T;
	}

	return res.json();
}

export function login(data: LoginRequest) {
	return request<LoginResponse>('/auth/login', {
		method: 'POST',
		body: JSON.stringify(data),
	});
}

export function logout() {
	return request<void>('/auth/logout', { method: 'POST' });
}

export function me() {
	return request<Me>('/auth/me');
}

export function stats() {
	return request<Stats>('/stats');
}

export function trends(days: number) {
	return request<TrendsResponse>(`/stats/trends?days=${days}`);
}

export function listMessages(params?: ListMessagesParams) {
	const query = new URLSearchParams();
	if (params) {
		for (const [key, value] of Object.entries(params)) {
			if (value !== undefined) query.set(key, String(value));
		}
	}
	const qs = query.toString();
	return request<ListMessagesResponse>(`/messages${qs ? `?${qs}` : ''}`);
}

export function getMessage(id: string) {
	return request<MessageDetail>(`/messages/${encodeURIComponent(id)}`);
}

export function sendMessage(data: SendMessageRequest) {
	return request<MessageDetail>('/messages', {
		method: 'POST',
		body: JSON.stringify(data),
	});
}

export function listDevices() {
	return request<Device[]>('/devices');
}

export function deleteDevice(id: string) {
	return request<void>(`/devices/${encodeURIComponent(id)}`, {
		method: 'DELETE',
	});
}

export function listWebhooks() {
	return request<Webhook[]>('/webhooks');
}

export function createWebhook(data: CreateWebhookRequest) {
	return request<Webhook>('/webhooks', {
		method: 'POST',
		body: JSON.stringify(data),
	});
}

export function deleteWebhook(id: string) {
	return request<void>(`/webhooks/${encodeURIComponent(id)}`, {
		method: 'DELETE',
	});
}

export async function getSettings(): Promise<DeviceSettings> {
	const raw = await request<Record<string, unknown>>('/settings');
	return fromServer(raw);
}

export async function updateSettings(data: DeviceSettings): Promise<DeviceSettings> {
	const raw = await request<Record<string, unknown>>('/settings', {
		method: 'PATCH',
		body: JSON.stringify(toServer(data)),
	});
	return fromServer(raw);
}

export function createToken(data: CreateTokenRequest) {
	return request<CreateTokenResponse>('/tokens', {
		method: 'POST',
		body: JSON.stringify(data),
	});
}

export function revokeToken(jti: string) {
	return request<void>(`/tokens/${encodeURIComponent(jti)}`, {
		method: 'DELETE',
	});
}
