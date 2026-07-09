export interface Stats {
	messagesSent: number;
	messagesPending: number;
	messagesFailed: number;
	devicesOnline: number;
	devicesActive: number;
	devicesTotal: number;
}

export interface Me {
	login: string;
}

export interface LoginRequest {
	login: string;
	password: string;
}

export interface LoginResponse {
	login: string;
}

export interface Device {
	id: string;
	name: string;
	lastSeen: string;
	isOnline: boolean;
	createdAt: string;
}

export interface MessageListItem {
	id: string;
	deviceId: string;
	state: string;
	recipients: RecipientState[];
	textPreview: string;
	createdAt: string;
}

export interface MessageDetail {
	id: string;
	deviceId: string;
	state: string;
	isHashed: boolean;
	isEncrypted: boolean;
	recipients: RecipientState[];
	states: Record<string, string>;
	textMessage?: { text: string };
	dataMessage?: { data: string; port: number };
	hashedMessage?: { hash: string };
}

export interface RecipientState {
	phoneNumber: string;
	state: string;
	error?: string;
}

export interface SendMessageRequest {
	phoneNumbers: string[];
	text: string;
	deviceId?: string;
	simNumber?: number;
	priority?: number;
}

export interface ListMessagesParams {
	limit?: number;
	offset?: number;
	state?: string;
	deviceId?: string;
	from?: string;
	to?: string;
}

export interface ListMessagesResponse {
	items: MessageListItem[];
	total: number;
}

export interface Webhook {
	id: string;
	url: string;
	event: string;
	deviceId?: string;
}

export interface CreateWebhookRequest {
	url: string;
	event: string;
	deviceId?: string;
}

export interface DeviceSettings {
	encryption?: { passphrase?: string };
	messages?: {
		sendIntervalMin?: number;
		sendIntervalMax?: number;
		limitPeriod?: 'Disabled' | 'PerMinute' | 'PerHour' | 'PerDay';
		limitValue?: number;
		simSelectionMode?: 'OSDefault' | 'RoundRobin' | 'Random';
		logLifetimeDays?: number;
		processingOrder?: 'LIFO' | 'FIFO';
	};
	ping?: { intervalSeconds?: number };
	logs?: { lifetimeDays?: number };
	webhooks?: {
		internetRequired?: boolean;
		retryCount?: number;
		signingKey?: string;
	};
	gateway?: {
		cloudUrl?: string;
		privateToken?: string;
		notificationChannel?: 'AUTO' | 'SSE_ONLY';
	};
}

export interface CreateTokenRequest {
	ttl?: number;
	scopes: string[];
}

export interface CreateTokenResponse {
	id: string;
	token_type: string;
	access_token: string;
	refresh_token?: string;
	expires_at: string;
}

export interface MessageReceivedPayload {
	sender: string;
	message: string;
}

export interface MessageStateChangedPayload {
	messageId: string;
	state: string;
}

export interface DeviceStatusChangedPayload {
	deviceId: string;
	isOnline: boolean;
}

export type SseEvent =
	| { type: 'message.received'; payload: MessageReceivedPayload }
	| { type: 'message.state_changed'; payload: MessageStateChangedPayload }
	| { type: 'device.status_changed'; payload: DeviceStatusChangedPayload }
	| { type: 'stats.updated'; payload: Record<string, never> };
