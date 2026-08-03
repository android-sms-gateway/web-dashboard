export function niceMax(values: number[]): number {
	const v = Math.max(0, ...values);
	if (v <= 0) return 1;
	const pow = Math.pow(10, Math.floor(Math.log10(v)));
	const frac = v / pow;
	if (frac <= 1) return pow;
	if (frac <= 2) return 2 * pow;
	if (frac <= 5) return 5 * pow;
	return 10 * pow;
}

export function tickValues(max: number): number[] {
	const mid = Math.round(max / 2);
	return Array.from(new Set([0, mid, max]));
}

export function formatTick(v: number): string {
	if (v >= 1000) {
		const k = v / 1000;
		return `${k % 1 === 0 ? k : k.toFixed(1)}k`;
	}
	return String(v);
}

export function xLabelIndexes(n: number): number[] {
	if (n <= 1) return n === 1 ? [0] : [];
	const step = Math.ceil(n / 6);
	const idx: number[] = [];
	for (let i = 0; i < n; i += step) idx.push(i);
	if (n - 1 - idx[idx.length - 1] >= 2) idx.push(n - 1);
	return idx;
}
