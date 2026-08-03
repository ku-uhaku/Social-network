"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { getPost } from "@/lib/api/posts";
import PostCard from "@/components/posts/PostCard";

export default function PostDetailPage({ params }) {
  const { postId } = params;
  const router = useRouter();
  const { user, loading } = useAuth();
  const [post, setPost] = useState(null);
  const [loadingPost, setLoadingPost] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!loading && !user) {
      router.replace("/login");
    }
  }, [loading, user, router]);

  useEffect(() => {
    async function loadPost() {
      try {
        const response = await getPost(postId);
        setPost(response?.data || null);
      } catch (err) {
        setError(err?.message || "Could not load post.");
      } finally {
        setLoadingPost(false);
      }
    }

    loadPost();
  }, [postId]);

  if (loading || (!user && !loading)) {
    return <div className="postsPlaceholder">Checking authentication…</div>;
  }

  if (loadingPost) {
    return <div className="postsPlaceholder">Loading post…</div>;
  }

  if (error) {
    return <div className="postsError">{error}</div>;
  }

  if (!post) {
    return <div className="postsPlaceholder">Post not found.</div>;
  }

  return (
    <section className="postsContainer">
      <div className="postDetails">
        <PostCard post={post} isFeed={false} />
      </div>
    </section>
  );
}
