import { describe, expect, it } from 'vitest';
import { smsStats } from './sms';

// Mirrors at-gateway/internal/modem/gsm7_test.go: the part-count table used by
// SegmentCount. A part holds 160 GSM-7 septets (153 once the concatenation
// UDH is present; extension runes like the Euro sign cost two septets) or
// 70 UCS-2 characters (67 with the UDH). Non-GSM-7 runes force UCS-2.
describe('smsStats part counts (parity with gsm7_test.go)', () => {
	const cases = [
		{ name: 'ascii single part', text: 'a'.repeat(160), want: 1, encoding: 'GSM-7' },
		{ name: 'ascii just over one part', text: 'a'.repeat(161), want: 2, encoding: 'GSM-7' },
		{ name: 'ascii part boundary', text: 'a'.repeat(153), want: 1, encoding: 'GSM-7' },
		{ name: 'ascii two full parts', text: 'a'.repeat(306), want: 2, encoding: 'GSM-7' },
		{ name: 'ascii just over two parts', text: 'a'.repeat(307), want: 3, encoding: 'GSM-7' },
		{ name: 'extended rune two septets', text: '€'.repeat(80), want: 1, encoding: 'GSM-7' },
		{ name: 'extended rune over one part', text: '€'.repeat(81), want: 2, encoding: 'GSM-7' },
		{ name: 'ucs2 single part', text: 'п'.repeat(70), want: 1, encoding: 'UCS-2' },
		{ name: 'ucs2 just over one part', text: 'п'.repeat(71), want: 2, encoding: 'UCS-2' },
		{ name: 'ucs2 two full parts', text: 'п'.repeat(134), want: 2, encoding: 'UCS-2' },
		{ name: 'ucs2 just over two parts', text: 'п'.repeat(135), want: 3, encoding: 'UCS-2' },
		{ name: 'emoji surrogate pair', text: '👍', want: 1, encoding: 'UCS-2' },
		{ name: 'protocol ceiling', text: 'a'.repeat(255 * 153), want: 255, encoding: 'GSM-7' },
		{ name: 'over protocol ceiling', text: 'a'.repeat(255 * 153 + 1), want: 256, encoding: 'GSM-7' },
		{ name: 'empty', text: '', want: 0, encoding: null },
	] as const;

	for (const { name, text, want, encoding } of cases) {
		it(name, () => {
			const stats = smsStats(text);
			expect(stats.parts).toBe(want);
			expect(stats.encoding).toBe(encoding);
			if (want === 0) {
				expect(stats.partUsage).toEqual([]);
			} else {
				expect(stats.partUsage.length).toBe(want);
				const limit = want === 1 ? stats.limits.single : stats.limits.multi;
				expect(Math.max(...stats.partUsage)).toBeLessThanOrEqual(limit);
			}
		});
	}
});

describe('smsStats char counting', () => {
	it('counts code points, not UTF-16 units', () => {
		expect(smsStats('👍').chars).toBe(1);
	});

	it('counts BMP characters one-to-one', () => {
		expect(smsStats('abcde').chars).toBe(5);
	});
});

describe('smsStats per-part limits', () => {
	it('uses 160 septets for a single GSM-7 part', () => {
		expect(smsStats('a').limits.single).toBe(160);
	});

	it('uses 70 units for a single UCS-2 part', () => {
		expect(smsStats('п').limits.single).toBe(70);
	});
});

// Pins the GSM-7 table edge cases against the Go encoder (warthog618):
// LF and CR live in the default alphabet (one septet each), while extension
// runes like form feed cost two septets (escape + code).
describe('smsStats escape handling', () => {
	it('line feed and carriage return are single default septets', () => {
		expect(smsStats('\n').partUsage).toEqual([1]);
		expect(smsStats('\r').partUsage).toEqual([1]);
	});

	it('form feed is an extension rune costing two septets', () => {
		expect(smsStats('\f').partUsage).toEqual([2]);
	});

	it('never splits an escape pair across a concatenation boundary', () => {
		// 'a'*152 + ESC + 0x0a + 'b'*310 = 464 septets; the 153-septet slice
		// would cut the escape pair, so the first part backs off to 152
		// (matches warthog618, Go encoder: [152,153,153,6]).
		expect(smsStats('a'.repeat(152) + '\f' + 'b'.repeat(310)).partUsage).toEqual([
			152, 153, 153, 6,
		]);
	});
});

// Pins literal U+001B parity with warthog618: warthog618's default alphabet
// maps a literal ESC rune to a single septet (index 27) and its chunk7Bit
// backoff fires only when the boundary septet is an ESC not preceded by an
// ESC. With 151 basic chars, ESC, form feed (ESC + 0x0a), and 7 more septets,
// the 153-septet first part ends "…ESC ESC" and the form-feed code opens part
// two: [153, 8], 2 parts (at-gateway SegmentCount: 2 parts).
describe('smsStats literal ESC parity (warthog618 chunk7Bit)', () => {
	it('mirrors the Go encoder part count for ESC-adjacent chunking', () => {
		const stats = smsStats('a'.repeat(151) + '\x1b' + '\f' + 'b'.repeat(7));
		expect(stats.encoding).toBe('GSM-7');
		expect(stats.parts).toBe(2);
		expect(stats.partUsage).toEqual([153, 8]);
	});

	it('a lone literal ESC is a single GSM-7 part', () => {
		const stats = smsStats('\x1b');
		expect(stats.encoding).toBe('GSM-7');
		expect(stats.parts).toBe(1);
		expect(stats.partUsage).toEqual([1]);
	});
});