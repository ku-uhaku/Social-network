"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { getFeed } from "@/lib/api/posts";
import PostCard from "@/components/posts/PostCard";
import NailButton from "@/components/shared/NailButton";

const PAGE_LIMIT = 10;

export default function HomePage() {
  const { user } = useAuth();
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [nextCursor, setNextCursor] = useState(null);
  const [hasMore, setHasMore] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    async function loadFeed() {
      try {
        const response = await getFeed({ limit: PAGE_LIMIT });
        const data = response?.data || {};
        setPosts(data.posts || []);
        setNextCursor(data.next_cursor || null);
        setHasMore(Boolean(data.has_more));
      } catch (err) {
        setError(err?.message || "Could not load feed.");
      } finally {
        setLoading(false);
      }
    }

    loadFeed();
  }, []);

  async function loadMore() {
    if (loadingMore || !nextCursor) return;
    setLoadingMore(true);
    try {
      const response = await getFeed({ limit: PAGE_LIMIT, cursor: nextCursor });
      const data = response?.data || {};
      setPosts((prev) => [...prev, ...(data.posts || [])]);
      setNextCursor(data.next_cursor || null);
      setHasMore(Boolean(data.has_more));
    } catch (err) {
      setError(err?.message || "Could not load more posts.");
    } finally {
      setLoadingMore(false);
    }
  }

  return (
    <section className="postsContainer">
      <div className="feedHeader">
        <div>
          <h1 className="feedTitle">Home feed</h1>
          <p className="feedSubtitle">
            {user ? `Welcome back, ${user.username}.` : "Your latest posts appear below."}
          </p>
        </div>
        <Link href="/posts/create" className="createPostLink">
          <NailButton>Create post</NailButton>
        </Link>
      </div>

      {loading && <div className="postsPlaceholder">Loading your feed…</div>}
      {error && <div className="postsError">{error}</div>}
      {!loading && !error && posts.length === 0 && (
        <div className="postsPlaceholder">No posts available yet. Create the first one.</div>
      )}

      <div className="feedList">
        {posts.map((post) => (
          <PostCard key={post.id} post={post} isFeed={true} />
        ))}
      </div>

      {!loading && !error && hasMore && (
        <div className="feedLoadMore">
          <NailButton onClick={loadMore} disabled={loadingMore}>
            {loadingMore ? "Loading…" : "Load more"}
          </NailButton>
        </div>
      )}
      {!loading && !error && !hasMore && posts.length > 0 && (
        <div className="feedEnd">{"You've reached the end of the feed."}</div>
      )}
    </section>
  );
}