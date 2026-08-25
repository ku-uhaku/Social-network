"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams, notFound } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { getPost, getComments } from "@/lib/api/posts";
import PostCard from "@/components/posts/PostCard";
import CommentCard from "@/components/posts/CommentCard";
import CommentCreate from "@/components/posts/CommentCreate";
import UsersSelect from "@/components/shared/UsersSelect";

export default function PostDetailPage() {
  const { postId } = useParams();
  const router = useRouter();
  const { user, loading } = useAuth();

  if (isNaN(Number(postId))) {
    notFound(); // invalid id
  }
  const [post, setPost] = useState(null);
  const [comments, setComments] = useState([]);
  const [loadingPost, setLoadingPost] = useState(true);
  const [loadingComments, setLoadingComments] = useState(true);
  const [error, setError] = useState("");
  const [commentsError, setCommentsError] = useState("");

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
        if (err?.message === "Post not found") {
          notFound();
        }
        setError(err?.message || "Could not load post.");
      } finally {
        setLoadingPost(false);
      }
    }

    loadPost();
  }, [postId]);

  useEffect(() => {
    async function loadComments() {
      try {
        const response = await getComments(postId);
        setComments(response?.data || []);
      } catch (err) {
        setCommentsError(err?.message || "Could not load comments.");
      } finally {
        setLoadingComments(false);
      }
    }

    loadComments();
  }, [postId]);

  function handleCommentCreated(comment) {
    setComments((prev) => [...prev, comment]);
    setPost((prev) => (prev ? { ...prev, comments_count: prev.comments_count + 1 } : prev));
  }

  if (loading || !user) {
    return <div className="postsPlaceholder">Checking authentication...</div>;
  }

  if (loadingPost) {
    return <div className="postsPlaceholder">Loading post...</div>;
  }

  if (error) {
    return <div className="postsError">{error}</div>;
  }

  if (!post) {
    notFound();
  }

  return (
    <section className="postsContainer">
      <div className="postDetails">
        <PostCard post={post} isFeed={false} />

        {post.privacy === "private" && post.viewers?.length > 0 && (
          <div className="postViewersSection">
            <h3 className="postViewersTitle">Visible only to:</h3>
            <UsersSelect users={post.viewers} selectable={false} />
          </div>
        )}

        <img className="postCommentSeparator" src="/images/post_comment_separator.png" alt="" />

        <div className="commentsSection">
          <h2 className="commentsTitle">Comments ({post.comments_count})</h2>

          <CommentCreate postId={postId} onCreated={handleCommentCreated} />

          {loadingComments ? (
            <div className="postsPlaceholder">Loading comments...</div>
          ) : commentsError ? (
            <div className="postsError">{commentsError}</div>
          ) : comments.length === 0 ? (
            <div className="postsPlaceholder">No comments yet. Be the first to comment!</div>
          ) : (
            <div className="commentsList">
              {comments.map((comment) => (
                <CommentCard key={comment.id} comment={comment} />
              ))}
            </div>
          )}
        </div>
      </div>
    </section>
  );
}
