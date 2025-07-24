import axios from "axios";

// Axios instance that always includes the cookie the backend sets
export const api = axios.create({
  baseURL: "", // Use empty string to use relative URLs for proxy
  withCredentials: true,
  timeout: 10000, // 10 second timeout
});

// Add request/response interceptors for better debugging
api.interceptors.request.use(
  (config) => {
    console.log(`Making ${config.method?.toUpperCase()} request to ${config.url}`);
    return config;
  },
  (error) => {
    console.error('Request interceptor error:', error);
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => {
    console.log(`Response from ${response.config.url}:`, response.status);
    return response;
  },
  (error) => {
    console.error('Response interceptor error:', error);
    return Promise.reject(error);
  }
);

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
