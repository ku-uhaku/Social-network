"use client";

import { useEffect, useState } from "react";
import { useParams, notFound } from "next/navigation";
import Link from "next/link";
import { getGroup, getGroupFeed, joinGroup, leaveGroup } from "@/lib/api/groups";
import { useAuth } from "@/contexts/AuthContext";
import { useGroupChat } from "@/contexts/GroupChatContext";
import PostCard from "@/components/posts/PostCard";
import NailButton from "@/components/shared/NailButton";
import UsersModal from "@/components/groups/UsersModal";
import GroupChat from "@/components/chat/GroupChat";
import "@/css/groups.css";

const pageLimit = 10;

export default function GroupDetailPage() {
  const { id } = useParams();
  const { user } = useAuth();
  const { unread } = useGroupChat();
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
  const [chatOpen, setChatOpen] = useState(false);

  const groupId = Number(id);
  if (isNaN(groupId)) notFound();
  const unreadCount = unread[groupId] || 0;

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
      applyFeed({}, false);
      setFeedLoaded(false);
    } catch (err) {
      setActionError(err?.message || "Could not leave group.");
    }
  }

  if (loadingGroup) return <div className="postsPlaceholder">Loading group...</div>;
  if (groupError) return <div className="postsError">{groupError}</div>;
  if (!group) notFound();

  return (
    <section className="postsContainer">
      <div className="groupHeader">
        <div className="groupHeaderInfo">
          <h1 className="feedTitle">{group.title}</h1>
          <p className="groupHeaderDescription">{group.description}</p>
        </div>

        <img className="group_title_separator" src="/images/group_title_separator.png" alt="" />

        <div className="groupHeaderActions">
          {membership === "accepted" && (
            <>
              <Link href={`/posts/create?group_id=${groupId}`}>
                <NailButton>Create post</NailButton>
              </Link>
              <Link href={`/group/${groupId}/events`}>
                <NailButton>Events</NailButton>
              </Link>
              <span className="groupChatButton">
                <NailButton onClick={() => setChatOpen(true)}>Chat</NailButton>
                {unreadCount > 0 && <span className="groupChatBadge">{unreadCount}</span>}
              </span>
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
        <UsersModal groupId={groupId} inviteMode onClose={() => setInviteOpen(false)} />
      )}

      {membersOpen && membership === "accepted" && (
        <UsersModal groupId={groupId} onClose={() => setMembersOpen(false)} />
      )}

      {membership === "accepted" ? (
        <>
          {!feedLoaded && !actionError && (
            <div className="postsPlaceholder">Loading group feed...</div>
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
                {loadingMore ? "Loading..." : "Load more"}
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

      {chatOpen && membership === "accepted" && (
        <GroupChat
          groupId={groupId}
          title={group.title}
          meId={user?.id}
          onClose={() => setChatOpen(false)}
        />
      )}
    </section>
  );
}

