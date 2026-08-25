"use client";

import { resolveMediaSrc } from "@/lib/utils";
import { useRouter } from "next/navigation";

export default function Avatar({ avatar, username, size = 64 }) {
  const router = useRouter();
  const src = resolveMediaSrc(avatar);

  return (
    <div className="avatar-image" onClick={() => router.push(`/profile/${username}`)} style={{ width: size, height: size }}>
      {src ? (
        <img className="avatar-image__img" src={src} alt={username} />
      ) : (
        <span className="avatar-image__fallback" style={{ fontSize: Math.round(size * 0.45) }}>
          {(username || "?").slice(0, 1).toUpperCase()}
        </span>
      )}
    </div>
  );
}
