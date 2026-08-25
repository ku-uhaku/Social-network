"use client";

import Avatar from "@/components/shared/Avatar";
import "@/css/createPost.css";

export default function UsersSelect({
  users = [],
  loading = false,
  selected = [],
  onToggle = () => {},
  selectable = true,
}) {
  if (loading) {
    return <div>Loading users...</div>;
  }

  if (users.length === 0) {
    return <div>No users found</div>;
  }

  return (
    <div className="usersChecklist">
      {users.map((user) => {
        const checked = selected.includes(user.id);
        const handleChange = (event) =>
          onToggle(user.id, event.target.checked);
        return (
          <label key={user.id} className="checkboxLabel">
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
            <div className="userSelectInfo">
              {user.avatar && (
                <Avatar
                  avatar={user.avatar}
                  username={user.username}
                  size={32}
                />
              )}
              <span className="userSelectName">
                {user.first_name} {user.last_name}
              </span>
            </div>
          </label>
        );
      })}
    </div>
  );
}
