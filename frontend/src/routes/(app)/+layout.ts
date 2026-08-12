// This whole section is an authenticated, localStorage-driven app shell —
// there's nothing to server-render, and SSR would just mismatch against the
// client-only auth state anyway.
export const ssr = false;
