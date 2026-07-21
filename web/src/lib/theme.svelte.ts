type Theme = 'light' | 'dark' | 'system';

let current = $state<Theme>('system');
let mediaQuery: MediaQueryList | null = null;
let mediaHandler: ((e: MediaQueryListEvent) => void) | null = null;

function apply(resolved: 'light' | 'dark') {
	if (typeof document === 'undefined') return;
	const root = document.documentElement;
	if (resolved === 'dark') {
		root.classList.add('dark');
	} else {
		root.classList.remove('dark');
	}
}

function getResolved(): 'light' | 'dark' {
	return current === 'system' ? (mediaQuery?.matches ? 'dark' : 'light') : current;
}

function watchSystem() {
	if (typeof window === 'undefined') return;
	mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
	mediaHandler = () => {
		if (current === 'system') {
			apply(getResolved());
		}
	};
	mediaQuery.addEventListener('change', mediaHandler);
}

function unwatchSystem() {
	if (mediaQuery && mediaHandler) {
		mediaQuery.removeEventListener('change', mediaHandler);
		mediaQuery = null;
		mediaHandler = null;
	}
}

function load() {
	if (typeof localStorage === 'undefined') return;
	const stored = localStorage.getItem('theme') as Theme | null;
	current = stored === 'light' || stored === 'dark' ? stored : 'system';
	watchSystem();
	apply(getResolved());
}

function save(t: Theme) {
	current = t;
	if (typeof localStorage !== 'undefined') {
		localStorage.setItem('theme', t);
	}
	if (t === 'system') {
		watchSystem();
	} else {
		unwatchSystem();
	}
	apply(getResolved());
}

function toggle() {
	if (current === 'light') save('dark');
	else if (current === 'dark') save('system');
	else save('light');
}

function theme(): Theme {
	return current;
}

function resolvedTheme(): 'light' | 'dark' {
	return getResolved();
}

export { theme, resolvedTheme, load, toggle };
