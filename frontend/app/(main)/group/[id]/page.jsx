"use client";

import { use, useEffect, useState } from "react";
import { notFound } from "next/navigation";
import Link from "next/link";
import { getGroup, getGroupFeed, joinGroup, leaveGroup, inviteUsers, getAllUsers, getGroupMembers } from "@/lib/api/groups";
import PostCard from "@/components/posts/PostCard";
import NailButton from "@/components/shared/NailButton";
import UsersSelect from "@/components/shared/UsersSelect";
import Avatar from "@/components/shared/Avatar";
import "@/css/groups.css";

const pageLimit = 10;

export default function GroupDetailPage({ params }) {
  const { id } = use(params);
  const [group, setGroup] = useState(null);
  const [membership, setMembership] = useState("none");
  const [loadingGroup, setLoadingGroup] = useState(true);
  const [groupError, setGroupError] = useState("");
  const [posts, setPosts] = useState([]);
  const [feedLoaded, setFeedLoaded] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [nextCursor, setNextCursor] = useState(null);
  const [hasMore, setHasMore] = useState(false);
  const [actionError, setActionError] = useState("");
  const [inviteOpen, setInviteOpen] = useState(false);
  const [membersOpen, setMembersOpen] = useState(false);

  const groupId = Number(id);
  if (isNaN(groupId)) {
    notFound();
  }

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const response = await getGroup(groupId);
        if (cancelled) return;
        setGroup(response?.data?.group || null);
        setMembership(response?.data?.membership || "none");
      } catch (err) {
        if (err?.message === "Group not found") {
          notFound();
        }
        if (!cancelled) setGroupError(err?.message || "Could not load group.");
      } finally {
        if (!cancelled) setLoadingGroup(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [groupId]);

  function applyFeed(data, append) {
    setPosts((prev) => (append ? [...prev, ...(data.posts || [])] : data.posts || []));
    setNextCursor(data.next_cursor || null);
    setHasMore(Boolean(data.has_more));
  }


  useEffect(() => {
    if (membership !== "accepted") return;
    let cancelled = false;
    (async () => {
      try {
        const response = await getGroupFeed(groupId, { limit: pageLimit });
        if (cancelled) return;
        applyFeed(response?.data || {}, false);
        setFeedLoaded(true);
      } catch (err) {
        if (!cancelled) setActionError(err?.message || "Could not load group feed.");
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [groupId, membership]);

  async function loadMore() {
    if (loadingMore || !nextCursor) return;
    setLoadingMore(true);
    try {
      const response = await getGroupFeed(groupId, { limit: pageLimit, cursor: nextCursor });
      applyFeed(response?.data || {}, true);
    } catch (err) {
      setActionError(err?.message || "Could not load more posts.");
    } finally {
      setLoadingMore(false);
    }
  }

  async function handleJoin() {
    setActionError("");
    try {
      await joinGroup(groupId);
      const response = await getGroup(groupId);
      setMembership(response?.data?.membership || "none");
    } catch (err) {
      setActionError(err?.message || "Could not join group.");
    }
  }

  async function handleLeave() {
    setActionError("");
    try {
      await leaveGroup(groupId);
      setMembership("none");
      setPosts([]);
      setFeedLoaded(false);
      setNextCursor(null);
      setHasMore(false);
    } catch (err) {
      setActionError(err?.message || "Could not leave group.");
    }
  }

  if (loadingGroup) {
    return <div className="postsPlaceholder">Loading group…</div>;
  }

  if (groupError) {
    return <div className="postsError">{groupError}</div>;
  }

  if (!group) {
    notFound();
  }

  return (
    <section className="postsContainer">
      <div className="groupHeader">
        <div className="groupHeaderInfo">
          <h1 className="feedTitle">{group.title}</h1>
          <p className="groupHeaderDescription">{group.description}</p>
        </div>

        <div className="groupHeaderActions">
          {membership === "accepted" && (
            <>
              <Link href={`/posts/create?group_id=${groupId}`}>
                <NailButton>Create post</NailButton>
              </Link>
              <Link href={`/group/${groupId}/events`}>
                <NailButton>Events</NailButton>
              </Link>
              <NailButton onClick={() => setMembersOpen(true)}>Members</NailButton>
              <NailButton onClick={() => setInviteOpen(!inviteOpen)}>
                {inviteOpen ? "Close invite" : "Invite"}
              </NailButton>
              <NailButton onClick={handleLeave}>Leave group</NailButton>
            </>
          )}
          {membership === "pending" && (
            <span className="groupPendingLabel">pending</span>
          )}
          {membership === "none" && (
            <NailButton onClick={handleJoin}>Join group</NailButton>
          )}
        </div>
      </div>

      {actionError && <div className="postsError">{actionError}</div>}

      {inviteOpen && membership === "accepted" && (
        <InviteModal
          groupId={groupId}
          onClose={() => setInviteOpen(false)}
          onInvited={() => setActionError("")}
        />
      )}

      {membersOpen && membership === "accepted" && (
        <MembersModal
          groupId={groupId}
          onClose={() => setMembersOpen(false)}
        />
      )}

      {membership === "accepted" ? (
        <>
          {!feedLoaded && !actionError && (
            <div className="postsPlaceholder">Loading group feed…</div>
          )}
          {feedLoaded && posts.length === 0 && !actionError && (
            <div className="postsPlaceholder">No posts in this group yet.</div>
          )}

          <div className="feedList">
            {posts.map((post) => (
              <PostCard key={post.id} post={post} isFeed={true} />
            ))}
          </div>

          {feedLoaded && hasMore && (
            <div className="feedLoadMore">
              <NailButton onClick={loadMore} disabled={loadingMore}>
                {loadingMore ? "Loading…" : "Load more"}
              </NailButton>
            </div>
          )}
          {feedLoaded && !hasMore && posts.length > 0 && (
            <div className="feedEnd">{"You've reached the end of the feed."}</div>
          )}
        </>
      ) : (
        <div className="postsPlaceholder">Join this group to see its posts.</div>
      )}
    </section>
  );
}

function InviteModal({ groupId, onClose, onInvited }) {
  const [users, setUsers] = useState([]);
  const [selected, setSelected] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const [usersRes, membersRes] = await Promise.all([
          getAllUsers(),
          getGroupMembers(groupId),
        ]);
        if (cancelled) return;
        const allUsers = usersRes?.data || [];
        const memberIds = new Set((membersRes?.data || []).map((m) => m.id));
        // Skip users that are already members of the group
        setUsers(allUsers.filter((u) => !memberIds.has(u.id)));
      } catch (err) {
        if (!cancelled) setError(err?.message || "Could not load users.");
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [groupId]);

  function toggleUser(userId, checked) {
    setSelected((prev) =>
      checked ? [...prev, userId] : prev.filter((id) => id !== userId)
    );
  }

  async function handleInvite(event) {
    event.preventDefault();
    if (selected.length === 0) {
      setError("Select at least one user to invite.");
      return;
    }
    setError("");
    setSubmitting(true);
    try {
      await inviteUsers(groupId, selected);
      onInvited();
      onClose();
    } catch (err) {
      setError(err?.message || "Could not send invitations.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="groupInviteOverlay">
      <form className="groupInvitePanel" onSubmit={handleInvite}>
        <h2 className="groupInviteTitle">Invite members</h2>
        {error && <div className="postsError">{error}</div>}

        <div className="groupInviteList">
          <UsersSelect
            users={users}
            loading={loading}
            selected={selected}
            onToggle={toggleUser}
          />
        </div>

        <div className="groupInviteActions">
          <NailButton type="submit" disabled={submitting || loading}>
            {submitting ? "Sending..." : "Invite"}
          </NailButton>
          <NailButton onClick={onClose}>Cancel</NailButton>
        </div>
      </form>
    </div>
  );
}

function MembersModal({ groupId, onClose }) {
  const [members, setMembers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;
    getGroupMembers(groupId)
      .then((res) => {
        if (!cancelled) setMembers(res?.data || []);
      })
      .catch((err) => {
        if (!cancelled) setError(err?.message || "Could not load members.");
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [groupId]);

  return (
    <div className="groupInviteOverlay">
      <div className="groupInvitePanel">
        <h2 className="groupInviteTitle">Members</h2>
        {error && <div className="postsError">{error}</div>}

        {loading ? (
          <div className="postsPlaceholder">Loading members…</div>
        ) : members.length === 0 ? (
          <div className="postsPlaceholder">No members yet.</div>
        ) : (
          <div className="groupMemberList">
            {members.map((m) => (
              <Link
                key={m.id}
                href={`/profile/${m.username}`}
                className="groupMemberItem"
                onClick={onClose}
              >
                <Avatar avatar={m.avatar} username={m.username} size={32} />
                <span className="followerName">
                  {m.first_name} {m.last_name} (@{m.username})
                </span>
              </Link>
            ))}
          </div>
        )}

        <div className="groupInviteActions">
          <NailButton onClick={onClose}>Close</NailButton>
        </div>
      </div>
    </div>
  );
}

