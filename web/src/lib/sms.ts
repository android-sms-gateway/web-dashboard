export type SmsEncoding = 'GSM-7' | 'UCS-2';

export interface SmsStats {
	chars: number;
	encoding: SmsEncoding | null;
	parts: number;
	partUsage: number[];
	limits: { single: number; multi: number };
}

const GSM7_EXT = 0x1b;

const GSM7_DEFAULT = new Map<string, number>(
	[
		'@£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞ\x1bÆæßÉ',
		' !"#¤%&\'()*+,-./0123456789:;<=>?',
		'¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿',
		'abcdefghijklmnopqrstuvwxyzäöñüà',
	]
		.join('')
		.split('')
		.map((c, i) => [c, i] as const),
);

const GSM7_EXTENSION = new Map<string, number>([
	['\f', 0x0a],
	['\n', 0x0d],
	['^', 0x14],
	['{', 0x28],
	['}', 0x29],
	['\\', 0x2f],
	['[', 0x3c],
	['~', 0x3d],
	[']', 0x3e],
	['|', 0x40],
	['€', 0x65],
]);

// A part without a concatenation UDH holds 140 octets of user data:
// 160 GSM-7 septets or 70 UCS-2 (UTF-16) characters. A concatenated part
// loses 6 octets to the 8-bit concatenation UDH: 153 septets or 67 units.
const GSM7_SINGLE_LIMIT = 160;
const GSM7_MULTI_LIMIT = 153;
const UCS2_SINGLE_LIMIT = 70;
const UCS2_MULTI_LIMIT = 67;

const HIGH_SURROGATE_START = 0xd800;
const HIGH_SURROGATE_END = 0xdc00;

// Encodes the text to GSM-7 septets (extension runes cost two: escape + code)
// or null when any rune is outside the GSM-7 default alphabet.
function encodeGsm7(text: string): number[] | null {
	const septets: number[] = [];
	for (const c of text) {
		const plain = GSM7_DEFAULT.get(c);
		if (plain !== undefined) {
			septets.push(plain);
			continue;
		}
		const ext = GSM7_EXTENSION.get(c);
		if (ext === undefined) {
			return null;
		}
		septets.push(GSM7_EXT, ext);
	}
	return septets;
}

// Converts the text to UTF-16 code units (surrogate pairs count as two).
function utf16Units(text: string): number[] {
	const units: number[] = [];
	for (let i = 0; i < text.length; i++) {
		units.push(text.charCodeAt(i));
	}
	return units;
}

// Splits a GSM-7 septet array into concatenated parts, never splitting an
// escape pair across a boundary (mirrors warthog618 chunk7Bit).
function chunkGsm7(septets: number[], bs: number): number[][] {
	const chunks: number[][] = [];
	let start = 0;
	let end = bs;
	while (end < septets.length) {
		if (septets[end - 1] === GSM7_EXT && septets[end - 2] !== GSM7_EXT) {
			end--;
		}
		chunks.push(septets.slice(start, end));
		start = end;
		end = start + bs;
	}
	chunks.push(septets.slice(start));
	return chunks;
}

// Splits a UTF-16 unit array into concatenated parts, never splitting a
// surrogate pair across a boundary (mirrors warthog618 chunkUCS2).
function chunkUcs2(units: number[], bs: number): number[][] {
	const chunks: number[][] = [];
	let start = 0;
	let end = bs;
	while (end < units.length) {
		const last = units[end - 1];
		if (HIGH_SURROGATE_START <= last && last < HIGH_SURROGATE_END) {
			end--;
		}
		chunks.push(units.slice(start, end));
		start = end;
		end = start + bs;
	}
	chunks.push(units.slice(start));
	return chunks;
}

export function smsStats(text: string): SmsStats {
	const chars = [...text].length;
	if (text === '') {
		return {
			chars,
			encoding: null,
			parts: 0,
			partUsage: [],
			limits: { single: GSM7_SINGLE_LIMIT, multi: GSM7_MULTI_LIMIT },
		};
	}

	const septets = encodeGsm7(text);
	if (septets !== null) {
		if (septets.length <= GSM7_SINGLE_LIMIT) {
			return {
				chars,
				encoding: 'GSM-7',
				parts: 1,
				partUsage: [septets.length],
				limits: { single: GSM7_SINGLE_LIMIT, multi: GSM7_MULTI_LIMIT },
			};
		}
		const chunks = chunkGsm7(septets, GSM7_MULTI_LIMIT);
		return {
			chars,
			encoding: 'GSM-7',
			parts: chunks.length,
			partUsage: chunks.map((c) => c.length),
			limits: { single: GSM7_SINGLE_LIMIT, multi: GSM7_MULTI_LIMIT },
		};
	}

	const units = utf16Units(text);
	if (units.length <= UCS2_SINGLE_LIMIT) {
		return {
			chars,
			encoding: 'UCS-2',
			parts: 1,
			partUsage: [units.length],
			limits: { single: UCS2_SINGLE_LIMIT, multi: UCS2_MULTI_LIMIT },
		};
	}
	const chunks = chunkUcs2(units, UCS2_MULTI_LIMIT);
	return {
		chars,
		encoding: 'UCS-2',
		parts: chunks.length,
		partUsage: chunks.map((c) => c.length),
		limits: { single: UCS2_SINGLE_LIMIT, multi: UCS2_MULTI_LIMIT },
	};
}