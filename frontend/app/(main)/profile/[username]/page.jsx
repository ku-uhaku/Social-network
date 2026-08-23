"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams, notFound } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { getUserProfile, getUserPosts, updateProfilePrivacy, followUser, unfollowUser } from "@/lib/api/user";
import PostCard from "@/components/posts/PostCard";
import Avatar from "@/components/shared/Avatar";
import CharmToggle from "@/components/shared/CharmToggle";
import FollowListModal from "@/components/profile/FollowListModal";
import { formatDate } from "@/lib/utils";

export default function ProfilePage() {
  const { username } = useParams();
  const { user: currentUser, refresh: refreshAuth } = useAuth();

  const [profile, setProfile] = useState(null);
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [postsLoading, setPostsLoading] = useState(false);
  const [toggleLoading, setToggleLoading] = useState(false);
  const [followLoading, setFollowLoading] = useState(false);
  const [modalType, setModalType] = useState(null); // 'followers' | 'following' | null

  const isOwner = currentUser && profile && currentUser.id === profile.id;
  const isPrivate = profile &&  profile.is_public === 0;

  const fetchProfile = useCallback(async () => {
    try {
      const response = await getUserProfile(username);
      setProfile(response?.data || response);
    } catch(err) {
        notFound();
    }
  }, [username]);

  const fetchPosts = useCallback(async () => {
    setPostsLoading(true);
    try {
      const response = await getUserPosts(username);
      const data = response?.data || response;
      setPosts(Array.isArray(data) ? data : []);
    } catch {
      // Private profiles return 403
      setPosts([]);
    } finally {
      setPostsLoading(false);
    }
  }, [username]);

  useEffect(() => {
    let cancelled = false;
    async function init() {
      try {
        await Promise.all([fetchProfile(), fetchPosts()]);
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    init();
    return () => { cancelled = true; };
  }, [fetchProfile, fetchPosts]);

  const handleTogglePrivacy = async () => {
    if (!profile || toggleLoading) return;
    setToggleLoading(true);
    try {
      const newValue = profile.is_public === 0 ? 1 : 0;
      await updateProfilePrivacy(newValue);
      setProfile((prev) => ({ ...prev, is_public: newValue }));
      await refreshAuth();
    } catch {
      // ignore
    } finally {
      setToggleLoading(false);
    }
  };

  const handleFollow = async () => {
    if (!profile || followLoading) return;
    setFollowLoading(true);
    try {
      const wasAccepted = profile.follow_status === "accepted";
      if (wasAccepted || profile.follow_status === "pending") {
        await unfollowUser(profile.id);
        setProfile((prev) => ({
          ...prev,
          follow_status: "none",
          followers_count: wasAccepted ? Math.max(0, prev.followers_count - 1) : prev.followers_count,
        }));
      } else {
        const response = await followUser(profile.id);
        const newStatus = response?.data?.status || "accepted";
        setProfile((prev) => ({
          ...prev,
          follow_status: newStatus,
          followers_count: newStatus === "accepted" ? prev.followers_count + 1 : prev.followers_count,
        }));
      }
    } catch {
      // ignore
    } finally {
      setFollowLoading(false);
    }
  };

  if (loading) {
    return (
      <section className="profilePage">
        <div className="profilePlaceholder">Loading profile...</div>
      </section>
    );
  }

  if (!profile) {
    notFound();
  }

  const followButtonLabel = (() => {
    if (profile.follow_status === "accepted") return "Unfollow";
    if (profile.follow_status === "pending") return "Requested";
    return "Follow";
  })();

  const showPosts = isOwner || (profile && profile.is_public==1) || profile?.follow_status === "accepted";
  // console.log("profile:::::",profile);
  
  return (
    <section className="profilePage">
      <div className="profileCard">
        <div className="profileCardTop">
          <Avatar avatar={profile.avatar} username={profile.username} size={96} />
          {isOwner && (
            <div className="profileToggleRow">
              <span className="profileToggleLabel">{profile.is_public ? "Public" : "Private"}</span>
              <CharmToggle
                checked={profile.is_public === 1}
                onChange={handleTogglePrivacy}
                title={profile.is_public === 1 ? "Switch to private" : "Switch to public"}
              />
            </div>
          )}
        </div>

        <h1 className="profileUsername">{profile.username}</h1>

        <img className="profileSeparator" src="/images/profile_separator.png" alt="" />
        {
          (!loading && showPosts ) ? (
            <div className="profileDetails">
              <div className="profileDetail"><strong>FullName:</strong> {profile.first_name} {profile.last_name}</div>
              <div className="profileDetail"><strong>Email:</strong> {profile.email}</div>
              <div className="profileDetail"><strong>Gender:</strong> {profile.gender}</div>
              <div className="profileDetail"><strong>Date of Birth:</strong> {profile.date_of_birth}</div>
              {profile.about_me && (
                <div className="profileDetail"><strong>About me:</strong> {profile.about_me}</div>
              )}
              <div className="profileDetail"><strong>Joined:</strong> {formatDate(profile.created_at)}</div>
            </div>
          ):(
             <div className="profilePrivateNotice">
            This account is private.
          </div>
          )
        }
        <div className="profileStats">
          <button
            type="button"
            className="profileStatButton"
            onClick={() => setModalType("followers")}
          >
            <strong>{profile.followers_count}</strong> Followers
          </button>
          <button
            type="button"
            className="profileStatButton"
            onClick={() => setModalType("following")}
          >
            <strong>{profile.following_count}</strong> Following
          </button>
        </div>

        {!isOwner && (
          <button
            type="button"
            className={`profileFollowButton ${profile.follow_status}`}
            onClick={handleFollow}
            disabled={followLoading}
          >
            {followLoading ? "..." : followButtonLabel}
          </button>
        )}
      </div>

      <img className="postCommentSeparator" src="/images/post_comment_separator.png" alt="" />

      <div className="profilePosts">
        <h2 className="profilePostsTitle">Posts by {username}</h2>
        {postsLoading && <div className="profilePlaceholder">Loading posts...</div>}
        {!postsLoading && showPosts && posts.length === 0 && (
          <div className="profilePlaceholder">No posts yet.</div>
        )}
        {!postsLoading && showPosts && posts.map((post) => (
          <PostCard key={post.id} post={post} isFeed={true} />
        ))}
        {!postsLoading && !showPosts && (
          <div className="profilePrivateNotice">
            This account is private.
          </div>
        )}
      </div>

      {modalType && (
        <FollowListModal
          userId={profile.id}
          type={modalType}
          onClose={() => setModalType(null)}
        />
      )}
    </section>
  );
}
