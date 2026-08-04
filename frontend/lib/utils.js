const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export function resolveMediaSrc(value) {
  // image path is either a path in "backend/media" or "http.." path for testing
  if (!value) return null;

  const trimmed = value.trim();
  if (trimmed.startsWith("http://") || trimmed.startsWith("https://")) {
    return trimmed;
  }

  let normalized = trimmed.replace(/\\/g, "/");
  const lower = normalized.toLowerCase();
  const mediaIndex = lower.lastIndexOf("/media/");
  if (mediaIndex >= 0) {
    normalized = normalized.slice(mediaIndex + 7);
  }

  while (normalized.startsWith("/")) {
    normalized = normalized.slice(1);
  }

  return `${API_BASE}/media/${normalized}`;
}