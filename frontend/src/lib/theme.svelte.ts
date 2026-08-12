type Theme = 'light' | 'dark';

let current = $state<Theme>('light');

function apply(theme: Theme) {
	document.documentElement.setAttribute('data-theme', theme);
	localStorage.setItem('theme', theme);
	current = theme;
}

export const theme = {
	get current() {
		return current;
	},
	init() {
		current = (document.documentElement.getAttribute('data-theme') as Theme | null) ?? 'light';
	},
	toggle() {
		apply(current === 'dark' ? 'light' : 'dark');
	}
};
