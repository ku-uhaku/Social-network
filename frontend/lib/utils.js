export const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

// formatDate: convert sql DATETIME
export function formatDate(value, options) {
  if (!value) return "";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "";
  return date.toLocaleDateString(undefined, options);
}

export function isOldEnough(birthday) {
    const birthDate = new Date(birthday);
    console.log(birthDate);
    
    const today = new Date();
    console.log("today",today);
    
    
    let age = today.getFullYear() - birthDate.getFullYear();
     
    const month = today.getMonth() - birthDate.getMonth();

    if (
        month < 0 ||
        (month === 0 && today.getDate() < birthDate.getDate())
    ) {
        age--;
    }
    if (age>100){
      return false
    }
    
    
    return age >= 16;
}

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

// formatMessageTime: "14:05" for today, "Mar 4 14:05" otherwise
export function formatMessageTime(value) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "";

  const time = date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  if (date.toDateString() === new Date().toDateString()) return time;

  return `${date.toLocaleDateString([], { month: "short", day: "numeric" })} ${time}`;
}
