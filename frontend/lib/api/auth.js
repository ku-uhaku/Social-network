import { apiFetch } from "./fetcher";

export function register(formData) {
  return apiFetch("/api/v1/auth/register", {
    method: "POST",
    body: formData,
  });
}

export function login({ identifier, password }) {
  return apiFetch("/api/v1/auth/login", {
    method: "POST",
    body: { login: identifier, password },
  });
}

export function me() {
  return apiFetch("/api/v1/auth/me", { method: "GET" });
}

export function logout() {
  return apiFetch("/api/v1/auth/logout", { method: "POST" });
}