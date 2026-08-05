const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// Derive the WebSocket URL from the HTTP API base
export function getWebSocketUrl() {
  const url = new URL(API_BASE);
  url.protocol = "ws:";
  url.pathname = "/ws";
  return url.toString();
}