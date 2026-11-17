import { afterEach, describe, expect, it, vi } from 'vitest';
import { sendMessage } from './api';

const request = { phoneNumbers: ['+15555550100'], text: 'hello' };

function stubFetch(body: string, status: number) {
	const fetchMock = vi.fn().mockResolvedValue(
		new Response(body, {
			status,
			headers: { 'Content-Type': 'application/json' },
		})
	);
	vi.stubGlobal('fetch', fetchMock);
	return fetchMock;
}

afterEach(() => {
	vi.unstubAllGlobals();
});

describe('request error messages', () => {
	it('surfaces the server message instead of the raw JSON body', async () => {
		const fetchMock = stubFetch('{"message":"invalid phone number","code":400}', 400);

		await expect(sendMessage(request)).rejects.toThrow(/^invalid phone number$/);
		expect(fetchMock).toHaveBeenCalledOnce();
	});

	it('surfaces a queue-limit reason for 503 responses', async () => {
		stubFetch('{"message":"queue limits exceeded","code":503}', 503);

		await expect(sendMessage(request)).rejects.toThrow(/^queue limits exceeded$/);
	});

	it('falls back to the raw body when it is not a JSON error payload', async () => {
		stubFetch('upstream exploded', 502);

		await expect(sendMessage(request)).rejects.toThrow('upstream exploded');
	});

	it('ignores a non-string message field', async () => {
		stubFetch('{"message":42,"code":500}', 500);

		await expect(sendMessage(request)).rejects.toThrow('{"message":42,"code":500}');
	});

	it('falls back to a status line when the body is empty', async () => {
		stubFetch('', 502);

		await expect(sendMessage(request)).rejects.toThrow('API error: 502');
	});
});
