import { describe, expect, it } from 'vitest';
import { parseTTL } from './ttl';

describe('parseTTL', () => {
	it('accepts empty and numeric TTL values', () => {
		expect(parseTTL('')).toBeUndefined();
		expect(parseTTL(1800)).toBe(1800);
	});

	it('rejects invalid TTL values', () => {
		expect(parseTTL('abc')).toBeNull();
		expect(parseTTL('1.5')).toBeNull();
		expect(parseTTL(0)).toBeNull();
		expect(parseTTL(-1)).toBeNull();
	});
});
