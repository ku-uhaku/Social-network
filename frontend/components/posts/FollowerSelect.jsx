"use client";

import Avatar from "@/components/shared/Avatar";
import "@/css/createPost.css";

export default function FollowerSelect({
  followers = [],
  loading = false,
  selectedViewers = [],
  onToggleViewer = () => {},
  selectable = true,
}) {
  if (loading) {
    return <div>Loading followers...</div>;
  }

  if (followers.length === 0) {
    return <div>No followers found</div>;
  }

  return (
    <div className="followersChecklist">
      {followers.map((follower) => {
        const checked = selectedViewers.includes(follower.id);
        const handleChange = (event) =>
          onToggleViewer(follower.id, event.target.checked);
        return (
          <label key={follower.id} className="checkboxLabel">
            {selectable && (
              <>
                <input
                  type="checkbox"
                  checked={checked}
                  onChange={handleChange}
                />
                <span className="checkboxCustom"></span>
              </>
            )}
            <div className="followerInfo">
              {follower.avatar && (
                <Avatar
                  avatar={follower.avatar}
                  username={follower.username}
                  size={32}
                />
              )}
              <span className="followerName">
                {follower.first_name} {follower.last_name} (@{follower.username})
              </span>
            </div>
          </label>
        );
      })}
    </div>
  );
}
