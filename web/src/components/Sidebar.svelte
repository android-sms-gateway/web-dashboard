<script lang="ts">
	let { current, onNavigate, onLogout, isOpen, onToggle, onClose }: {
		current: string;
		onNavigate: (page: string) => void;
		onLogout: () => void;
		isOpen: boolean;
		onToggle: () => void;
		onClose?: () => void;
	} = $props();

	const items = [
		{ id: 'dashboard', label: 'Dashboard', icon: '📊' },
		{ id: 'messages', label: 'Messages', icon: '✉️' },
		{ id: 'devices', label: 'Devices', icon: '📱' },
		{ id: 'webhooks', label: 'Webhooks', icon: '🔗' },
		{ id: 'settings', label: 'Settings', icon: '⚙️' },
		{ id: 'tokens', label: 'API Tokens', icon: '🔑' },
	];

	function handleNav(page: string) {
		onNavigate(page);
		onClose?.();
	}

	$effect(() => {
		if (!isOpen) return;
		function listener(e: KeyboardEvent) { if (e.key === "Escape") onClose?.(); }
		window.addEventListener("keydown", listener);
		return () => window.removeEventListener("keydown", listener);
	});
</script>

{#if isOpen}
	<div class="backdrop" onclick={onClose} role="presentation"></div>
{/if}

<aside class="sidebar" class:open={isOpen}>
	<div class="sidebar-header">
		<div class="brand">SMS Dashboard</div>
		<button class="close-btn" onclick={onToggle} aria-label="Close sidebar">×</button>
	</div>

	<nav>
		{#each items as item}
			<button
				class="nav-item"
				class:active={current === item.id}
				onclick={() => handleNav(item.id)}
			>
				<span class="icon">{item.icon}</span>
				{item.label}
			</button>
		{/each}
	</nav>

	<div class="spacer"></div>

	<button class="nav-item logout" onclick={() => { onLogout(); onClose?.(); }}>
		<span class="icon">🚪</span>
		Logout
	</button>
</aside>

<style>
	.backdrop {
		display: none;
	}

	.sidebar {
		width: 220px;
		background: #1e293b;
		color: #e2e8f0;
		display: flex;
		flex-direction: column;
		padding: 0;
	}

	.sidebar-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.brand {
		padding: 1.25rem 1rem;
		font-size: 1rem;
		font-weight: 700;
		color: #f8fafc;
		border-bottom: 1px solid #334155;
		flex: 1;
	}

	.close-btn {
		display: none;
	}

	nav {
		padding: 0.5rem 0;
	}

	.nav-item {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		width: 100%;
		padding: 0.625rem 1rem;
		background: none;
		border: none;
		color: #94a3b8;
		font-size: 0.875rem;
		cursor: pointer;
		text-align: left;
	}

	.nav-item:hover {
		background: #334155;
		color: #e2e8f0;
	}

	.nav-item.active {
		background: #2563eb;
		color: white;
	}

	.icon {
		font-size: 1rem;
		width: 1.25rem;
		text-align: center;
	}

	.spacer {
		flex: 1;
	}

	.logout {
		border-top: 1px solid #334155;
		margin-top: auto;
	}

	@media (max-width: 768px) {
		.backdrop {
			display: block;
			position: fixed;
			inset: 0;
			background: rgba(0, 0, 0, 0.4);
			z-index: 99;
		}

		.sidebar {
			position: fixed;
			top: 0;
			left: 0;
			bottom: 0;
			z-index: 100;
			transform: translateX(-100%);
			transition: transform 0.2s ease;
			box-shadow: 4px 0 20px rgba(0, 0, 0, 0.3);
		}

		.sidebar.open {
			transform: translateX(0);
		}

		.close-btn {
			display: block;
			background: none;
			border: none;
			color: #94a3b8;
			font-size: 1.5rem;
			cursor: pointer;
			padding: 0.75rem 1rem;
			line-height: 1;
		}

		.close-btn:hover {
			color: #e2e8f0;
		}
	}
</style>
