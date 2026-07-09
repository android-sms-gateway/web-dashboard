import { me } from '$lib/api';
import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const ssr = false;

export async function load({ url }: Parameters<LayoutLoad>[0]) {
	if (url.pathname.startsWith('/login') || url.pathname.startsWith('/logout')) {
		return;
	}

	try {
		const user = await me();
		return { user };
	} catch {
		throw redirect(302, '/login');
	}
}
