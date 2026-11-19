export function parseTTL(value: unknown): number | undefined | null {
	if (value === null || value === undefined) return undefined;

	const text = String(value).trim();
	if (!text) return undefined;
	if (!/^\d+$/.test(text)) return null;

	const parsed = Number(text);
	if (!Number.isSafeInteger(parsed) || parsed <= 0) return null;
	return parsed;
}
