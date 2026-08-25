"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import Avatar from "@/components/shared/Avatar";
import { getFollowers, getFollowing } from "@/lib/api/user";
import { useRouter } from "next/navigation";

export default function FollowListModal({ userId, type, onClose }) {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const router=useRouter()

  useEffect(() => {
    async function load() {
      setLoading(true);
      setError("");
      try {
        const response =
          type === "followers"
            ? await getFollowers(userId)
            : await getFollowing(userId);
        const data = response?.data || response || [];
        setUsers(Array.isArray(data) ? data : []);
      } catch (err) {
        if (
          err.status === 401 ||
          err.status === 403 ||
          err.status === 404 ||
          err.status >= 500
        ) {
          router.push(
            `/error?message=${encodeURIComponent(err.statusText)}`
          );
        }
        setError(err?.message || "Could not load list.");
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [userId, type]);

  const handleBackdrop = (e) => {
    if (e.target === e.currentTarget) onClose();
  };

  return (
    <div className="followListOverlay" onClick={handleBackdrop}>
      <div className="followListModal">
        <div className="followListModalHeader">
          <h2 className="followListModalTitle">
            {type === "followers" ? "Followers" : "Following"}
          </h2>
          <button
            type="button"
            className="followListCloseButton"
            onClick={onClose}
          >
            &times;
          </button>
        </div>

        <div className="followListModalBody">
          {loading && <div className="followListPlaceholder">Loading...</div>}
          {error && <div className="followListError">{error}</div>}
          {!loading && !error && users.length === 0 && (
            <div className="followListPlaceholder">No users found.</div>
          )}
          {!loading &&
            !error &&
            users.map((u) => (
              <Link
                key={u.id}
                href={`/profile/${u.username}`}
                className="followListItem"
                onClick={onClose}
              >
                <Avatar avatar={u.avatar} username={u.username} size={40} />
                <div className="followListItemInfo">
                  <span className="followListItemName">{u.first_name} {u.last_name}</span>
                  <span className="followListItemUsername">@{u.username}</span>
                </div>
              </Link>
            ))}
        </div>
      </div>
    </div>
  );
}