"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { createPost } from "@/lib/api/posts";
import ImageUploader from "@/components/shared/ImageUploader";

export default function CreatePostPage() {
  const router = useRouter();
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [privacy, setPrivacy] = useState("public");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    setSubmitting(true);

    try {
      const payload = {
        title: title.trim(),
        content: content.trim(),
        privacy,
      };

      const response = await createPost(payload);
      const post = response?.data;
      if (post?.id) {
        router.push(`/posts/${post.id}`);
      } else {
        setError("Unexpected response from the server.");
      }
    } catch (err) {
      setError(err?.message || "Could not create post.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className="postsContainer">
      <div className="postFormWrapper">
        <h1 className="feedTitle">Create a new post</h1>
        <form className="postForm" onSubmit={handleSubmit}>
          {error && <div className="postsError">{error}</div>}

          <div className="field">
            <label htmlFor="title">Title</label>
            <input
              id="title"
              value={title}
              onChange={(event) => setTitle(event.target.value)}
              placeholder="Post title"
              required
            />
          </div>

          <div className="field">
            <label htmlFor="content">Content</label>
            <textarea
              id="content"
              value={content}
              onChange={(event) => setContent(event.target.value)}
              placeholder="Write your post here..."
              required
            />
          </div>

          <div className="row">
            <div className="field">
              <label htmlFor="privacy">Privacy</label>
              <select
                id="privacy"
                value={privacy}
                onChange={(event) => setPrivacy(event.target.value)}
              >
                <option value="public">Public</option>
                <option value="almost_private">Almost Private</option>
                <option value="private">Private</option>
              </select>
            </div>
            <div className="field">
              <label htmlFor="image_url">Image (optional)</label>
              <ImageUploader avatar={null} name={title || "Post"} size={64} />
              
            </div>
          </div>

          <p className="guideText">You can add an image URL to include a visual preview on the post page.</p>

          <div className="postFormActions">
            <button className="button postFormButton" type="submit" disabled={submitting}>
              {submitting ? "Creating..." : "Publish post"}
            </button>
          </div>
        </form>
      </div>
    </section>
  );
}
