function resolveMediaSrc(value) {
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

  return `/media/${normalized}`;
}

export default function ImageUploader({ avatar, name, size = 64 }) {
  const label = name
    ? name
        .split(" ")
        .map((part) => part.charAt(0).toUpperCase())
        .slice(0, 2)
        .join("")
    : "?";

  const src = resolveMediaSrc(avatar);

  return (
    <div className="avatar-image" style={{ width: size, height: size }}>
      {src ? (
        <img className="avatar-image__img" src={src} alt={name || "avatar"} />
      ) : (
        <span className="avatar-image__fallback">{label}</span>
      )}
    </div>
  );
}
