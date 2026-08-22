import assert from 'node:assert/strict';
import { afterEach, describe, it } from 'node:test';

import { errorMessage, formatApiError, fromServer, getSettings, toServer, updateSettings } from './api.ts';

type FetchStub = (input: string, init?: RequestInit) => Promise<Response>;

function stubFetch(stub: FetchStub) {
	globalThis.fetch = ((input: RequestInfo | URL, init?: RequestInit) =>
		stub(String(input), init)) as typeof fetch;
}

const FULL_SERVER_SETTINGS = {
	encryption: { passphrase: 's3cret' },
	messages: {
		send_interval_min: 5,
		send_interval_max: 30,
		limit_period: 'PerMinute',
		limit_value: 10,
		sim_selection_mode: 'RoundRobin',
		log_lifetime_days: 7,
		processing_order: 'FIFO'
	},
	ping: { interval_seconds: 60 },
	logs: { lifetime_days: 14 },
	webhooks: { internet_required: false, retry_count: 3, signing_key: 'wh-key' },
	gateway: { cloud_url: 'https://gw.example.com', private_token: 'tok', notification_channel: 'SSE_ONLY' }
};

describe('fromServer', () => {
	it('maps snake_case sections to camelCase', () => {
		const mapped = fromServer(FULL_SERVER_SETTINGS);

		assert.deepEqual(mapped, {
			encryption: { passphrase: 's3cret' },
			messages: {
				sendIntervalMin: 5,
				sendIntervalMax: 30,
				limitPeriod: 'PerMinute',
				limitValue: 10,
				simSelectionMode: 'RoundRobin',
				logLifetimeDays: 7,
				processingOrder: 'FIFO'
			},
			ping: { intervalSeconds: 60 },
			logs: { lifetimeDays: 14 },
			webhooks: { internetRequired: false, retryCount: 3, signingKey: 'wh-key' },
			gateway: { cloudUrl: 'https://gw.example.com', privateToken: 'tok', notificationChannel: 'SSE_ONLY' }
		});
	});

	it('preserves falsy values like false and 0', () => {
		const mapped = fromServer({ webhooks: { internet_required: false, retry_count: 0 } });

		assert.deepEqual(mapped, { webhooks: { internetRequired: false, retryCount: 0 } });
	});

	it('omits null values', () => {
		const mapped = fromServer({ messages: { limit_period: null }, ping: {} });

		assert.deepEqual(mapped, {});
	});

	it('drops unknown keys so ghost snake_case keys never reach form state', () => {
		const mapped = fromServer({
			receiver: { content_provider_enabled: true },
			messages: { work_hours_start: '09:00', limit_period: 'PerDay' },
			unknown_section: { foo: 1 }
		});

		assert.deepEqual(mapped, { messages: { limitPeriod: 'PerDay' } });
	});

	it('returns empty object for empty input', () => {
		assert.deepEqual(fromServer({}), {});
	});

	it('is stable across a full map/unmap/map cycle', () => {
		const once = fromServer(FULL_SERVER_SETTINGS);
		const twice = fromServer(toServer(once));

		assert.deepEqual(twice, once);
	});
});

describe('toServer', () => {
	it('maps camelCase sections back to snake_case (round-trip)', () => {
		const roundTripped = toServer(fromServer(FULL_SERVER_SETTINGS));

		assert.deepEqual(roundTripped, FULL_SERVER_SETTINGS);
	});

	it('omits undefined, null and empty-string values', () => {
		const payload = toServer({
			messages: { sendIntervalMin: undefined, limitPeriod: null, processingOrder: '' },
			ping: {},
			encryption: { passphrase: 'keep-me' }
		});

		assert.deepEqual(payload, {
			encryption: { passphrase: 'keep-me' }
		} as Record<string, unknown>);
	});

	it('keeps falsy values like false and 0', () => {
		const payload = toServer({ webhooks: { internetRequired: false, retryCount: 0 } });

		assert.deepEqual(payload, { webhooks: { internet_required: false, retry_count: 0 } });
	});
});

describe('settings API mapping integration', () => {
	afterEach(() => {
		delete (globalThis as { fetch?: unknown }).fetch;
	});

	it('getSettings maps the server response to camelCase', async () => {
		let requestedPath = '';
		stubFetch(async (input) => {
			requestedPath = input;
			return new Response(JSON.stringify(FULL_SERVER_SETTINGS), { status: 200 });
		});

		const result = await getSettings();

		assert.equal(requestedPath, '/api/v1/settings');
		assert.deepEqual(result.messages.limitPeriod, 'PerMinute');
		assert.deepEqual(result.webhooks.signingKey, 'wh-key');
		assert.deepEqual(result.gateway.privateToken, 'tok');
	});

	it('updateSettings sends a snake_case PATCH body and maps the response back', async () => {
		let capturedPath = '';
		let capturedBody = '';
		stubFetch(async (input, init) => {
			capturedPath = input;
			capturedBody = String(init?.body ?? '');
			return new Response(JSON.stringify(FULL_SERVER_SETTINGS), { status: 200 });
		});

		const result = await updateSettings({
			messages: { limitPeriod: 'PerMinute', sendIntervalMin: 5 },
			encryption: { passphrase: 'abc' }
		});

		assert.equal(capturedPath, '/api/v1/settings');
		const sent = JSON.parse(capturedBody) as Record<string, unknown>;
		assert.deepEqual(sent, {
			messages: { limit_period: 'PerMinute', send_interval_min: 5 },
			encryption: { passphrase: 'abc' }
		});
		assert.deepEqual(result.encryption?.passphrase, 's3cret');
	});
});

describe('error detail surfacing', () => {
	afterEach(() => {
		delete (globalThis as { fetch?: unknown }).fetch;
	});

	it('formatApiError combines status and detail on one line', () => {
		assert.equal(formatApiError(400, '{"message":"failed on the isdefault tag"}\n'), 'API error 400: {"message":"failed on the isdefault tag"}');
		assert.equal(formatApiError(500, ''), 'API error 500');
	});

	it('formatApiError truncates long bodies', () => {
		const message = formatApiError(502, 'x'.repeat(1000));

		assert.ok(message.length < 400);
		assert.ok(message.startsWith('API error 502:'));
	});

	it('getSettings rejects with status and body detail included', async () => {
		stubFetch(async () => new Response('{"message":"boom detail"}', { status: 401 }));

		await assert.rejects(getSettings(), (err: Error) => {
			assert.match(err.message, /401/);
			assert.match(err.message, /boom detail/);
			return true;
		});
	});

	it('errorMessage falls back for unknown throwables', () => {
		assert.equal(errorMessage(new Error('detail')), 'detail');
		assert.equal(errorMessage('string error'), 'string error');
		assert.equal(errorMessage(undefined), 'Unknown error');
	});
});
