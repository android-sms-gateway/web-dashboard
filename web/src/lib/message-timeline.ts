interface TimelineEntry {
	state: string;
	time: string;
}

export function sortedStates(states: Record<string, string>): TimelineEntry[] {
	return Object.entries(states)
		.map(([state, value], index) => {
			const timestamp = value ? Date.parse(value) : NaN;
			return { state, index, timestamp, valid: !Number.isNaN(timestamp) };
		})
		.sort((a, b) => {
			if (a.valid && b.valid) return a.timestamp - b.timestamp || a.index - b.index;
			if (a.valid !== b.valid) return a.valid ? -1 : 1;
			return a.index - b.index;
		})
		.map(({ state, timestamp, valid }) => ({
			state,
			time: valid ? new Date(timestamp).toLocaleString() : '—',
		}));
}
