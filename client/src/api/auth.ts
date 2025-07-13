import axios from "axios";

// Axios instance that always includes the cookie the backend sets
export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080",
  withCredentials: true,
});

/**
 * Log in with email + password.
 * Returns { id, username, email } on success (matches your backend shape).
 */
export async function login(email: string, password: string) {
  const { data } = await api.post("/api/users/login", { email, password });
  return data as { id: string; username: string; email?: string };
}

/**
 * Sign up with username, email + password.
 * Returns { id, username, email } on success.
 */
export async function signup(username: string, email: string, password: string) {
  const { data } = await api.post("/api/users/signup", { username, email, password });
  return data as { id: string; username: string; email?: string };
}
