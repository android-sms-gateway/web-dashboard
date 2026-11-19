import { describe, expect, it } from 'vitest';
import { sortedStates } from './message-timeline';

describe('sortedStates', () => {
	it('sorts states chronologically, oldest first', () => {
		const states = {
			Sent: '2026-11-19T12:00:00Z',
			Pending: '2026-11-19T11:00:00Z',
			Delivered: '2026-11-19T13:00:00Z',
		};

		expect(sortedStates(states).map((entry) => entry.state)).toEqual([
			'Pending',
			'Sent',
			'Delivered',
		]);
	});

	it('formats valid timestamps with toLocaleString', () => {
		const timestamp = '2026-11-19T11:00:00Z';

		expect(sortedStates({ Pending: timestamp })[0].time).toBe(
			new Date(timestamp).toLocaleString(),
		);
	});

	it('renders invalid timestamps as dash and places them last', () => {
		const timestamp = '2026-11-19T12:00:00Z';
		const result = sortedStates({
			Pending: '',
			Sent: timestamp,
			Failed: 'not-a-date',
		});

		expect(result.map((entry) => entry.state)).toEqual(['Sent', 'Pending', 'Failed']);
		expect(result.map((entry) => entry.time)).toEqual([
			new Date(timestamp).toLocaleString(),
			'—',
			'—',
		]);
	});

	it('keeps insertion order for equal timestamps', () => {
		const timestamp = '2026-11-19T12:00:00Z';

		expect(sortedStates({ A: timestamp, B: timestamp, C: timestamp }).map((e) => e.state)).toEqual([
			'A',
			'B',
			'C',
		]);
	});

	it('handles an empty states map', () => {
		expect(sortedStates({})).toEqual([]);
	});
});
