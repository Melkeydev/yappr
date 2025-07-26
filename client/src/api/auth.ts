import axios from "axios";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "https://server.yappr.chat";

export const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
  timeout: 10000,
});

api.interceptors.request.use(
  (config) => {
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    return Promise.reject(error);
  },
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
export async function signup(
  username: string,
  email: string,
  password: string,
) {
  const { data } = await api.post("/api/users/signup", {
    username,
    email,
    password,
  });
  return data as { id: string; username: string; email?: string };
}
