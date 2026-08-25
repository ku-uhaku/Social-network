"use client";

import { useEffect, useState } from "react";
import { getAllUsers, getGroupMembers, inviteUsers } from "@/lib/api/groups";
import NailButton from "@/components/shared/NailButton";
import UsersSelect from "@/components/shared/UsersSelect";

export default function UsersModal({ groupId, inviteMode = false, onClose }) {
  const [users, setUsers] = useState([]);
  const [selected, setSelected] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        let list;
        if (inviteMode) {
          const [usersRes, membersRes] = await Promise.all([
            getAllUsers(),
            getGroupMembers(groupId),
          ]);
          // Skip members
          const memberIds = new Set((membersRes?.data || []).map((m) => m.id));
          list = (usersRes?.data || []).filter((u) => !memberIds.has(u.id));
        } else {
          const membersRes = await getGroupMembers(groupId);
          list = membersRes?.data || [];
        }
        if (!cancelled) setUsers(list);
      } catch (err) {
        if (!cancelled) setError(err?.message || "Could not load users.");
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [groupId, inviteMode]);

  function toggleUser(userId, checked) {
    setSelected((prev) =>
      checked ? [...prev, userId] : prev.filter((id) => id !== userId)
    );
  }

  async function handleInvite() {
    if (selected.length === 0) {
      setError("Select at least one user to invite.");
      return;
    }
    setError("");
    setSubmitting(true);
    try {
      await inviteUsers(groupId, selected);
      onClose();
    } catch (err) {
      setError(err?.message || "Could not send invitations.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="groupInviteOverlay">
      <div className="groupInvitePanel">
        <h2 className="groupInviteTitle">{inviteMode ? "Invite members" : "Members"}</h2>
        {error && <div className="postsError">{error}</div>}

        <div className="groupInviteList">
          <UsersSelect
            users={users}
            loading={loading}
            selected={selected}
            onToggle={toggleUser}
            selectable={inviteMode}
          />
        </div>

        <div className="groupInviteActions">
          {inviteMode && (
            <NailButton onClick={handleInvite} disabled={submitting || loading}>
              {submitting ? "Sending..." : "Invite"}
            </NailButton>
          )}
          <NailButton onClick={onClose}>{inviteMode ? "Cancel" : "Close"}</NailButton>
        </div>
      </div>
    </div>
  );
}