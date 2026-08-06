import { API_BASE } from "@/lib/utils";

export function getWebSocketUrl() {
  const url = new URL(API_BASE);
  url.protocol = "ws:";
  url.pathname = "/ws";
  return url.toString();
}