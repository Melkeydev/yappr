import axios from "axios";

// Axios instance that always includes the cookie the backend sets
export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080",
  withCredentials: true,
});

/**
 * Log in with email + password.
 * Returns { id, username } on success (matches your backend shape).
 */
export async function login(email: string, password: string) {
  const { data } = await api.post("/api/users/login", { email, password });
  return data as { id: string; username: string };
}
