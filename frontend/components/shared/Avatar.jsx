import { resolveMediaSrc } from "@/lib/utils";

export default function Avatar({ avatar, name, size = 64 }) {
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