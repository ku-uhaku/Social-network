"use client";

import { displayName, resolveMediaSrc } from "@/lib/utils";
import { useRouter } from "next/navigation";

export default function Avatar({ user, size = 64 }) {
  const router = useRouter();
  const src = resolveMediaSrc(user?.avatar);
  const name = displayName(user);

  return (
    <div className="avatar-image" onClick={() => router.push(`/profile/${user?.username}`)} style={{ width: size, height: size }}>
      {src ? (
        <img className="avatar-image__img" src={src} alt={name} />
      ) : (
        <span className="avatar-image__fallback" style={{ fontSize: Math.round(size * 0.45) }}>
          {(name || "?").slice(0, 1).toUpperCase()}
        </span>
      )}
    </div>
  );
}
