import { StoredToken } from '../types'

// Simple in-memory token store (use database in production)
class TokenStore {
  private tokens: Map<string, StoredToken> = new Map()
  private codeVerifiers: Map<string, string> = new Map()

  saveToken(userId: string, token: StoredToken): void {
    this.tokens.set(userId, token)
  }

  getToken(userId: string): StoredToken | undefined {
    return this.tokens.get(userId)
  }

  deleteToken(userId: string): void {
    this.tokens.delete(userId)
  }

  // Store code verifier temporarily during OAuth flow
  saveCodeVerifier(state: string, codeVerifier: string): void {
    this.codeVerifiers.set(state, codeVerifier)
  }

  getCodeVerifier(state: string): string | undefined {
    return this.codeVerifiers.get(state)
  }

  deleteCodeVerifier(state: string): void {
    this.codeVerifiers.delete(state)
  }

  getAllUserIds(): string[] {
    return Array.from(this.tokens.keys())
  }
}

export const tokenStore = new TokenStore()
