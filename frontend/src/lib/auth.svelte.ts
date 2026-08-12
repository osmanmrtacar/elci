import { browser } from '$app/environment';

const STORAGE_KEY = 'elci_token';

let token = $state<string | null>(browser ? localStorage.getItem(STORAGE_KEY) : null);

export const auth = {
	get token() {
		return token;
	},
	get isAuthenticated() {
		return token !== null;
	},
	login(newToken: string) {
		token = newToken;
		localStorage.setItem(STORAGE_KEY, newToken);
	},
	logout() {
		token = null;
		localStorage.removeItem(STORAGE_KEY);
	}
};
