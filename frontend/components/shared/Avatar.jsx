"use client";

import { resolveMediaSrc } from "@/lib/utils";
import { useRouter } from "next/navigation";

export default function Avatar({ avatar, username, size = 64 }) {
  const router = useRouter();
  const src = resolveMediaSrc(avatar) ||"/public/images/defaulte_avatar.jpeg";

  return (
    <div className="avatar-image" onClick={() => router.push(`/profile/${username}`)} style={{ width: size, height: size }}>
      {src ? (
        <img className="avatar-image__img" src={src} alt={username} />
      ) : (
        <span className="avatar-image__fallback">{username}</span>
      )}
    </div>
  );
}
