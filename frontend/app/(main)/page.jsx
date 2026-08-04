"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { getFeed } from "@/lib/api/posts";
import PostCard from "@/components/posts/PostCard";
import NailButton from "@/components/shared/NailButton";

export default function HomePage() {
  const { user } = useAuth();
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function loadFeed() {
      try {
        const response = await getFeed();
        setPosts(response?.data || []);
      } catch (err) {
        setError(err?.message || "Could not load feed.");
      } finally {
        setLoading(false);
      }
    }

    loadFeed();
  }, []);

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
    </section>
  );
}
