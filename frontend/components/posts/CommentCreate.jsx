"use client";

import { useState } from "react";
import { createComment } from "@/lib/api/posts";
import ImageUploadButton from "@/components/shared/ImageUploadButton";

export default function CommentCreate({ postId, onCreated }) {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [image, setImage] = useState(null);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [resetKey, setResetKey] = useState(0); // used to force refresh ImageUploadButton to remove preview

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    setSubmitting(true);

    try {
      const payload = image
        ? (() => {
            const formData = new FormData();
            formData.append("post_id", postId);
            formData.append("title", title.trim());
            formData.append("content", content.trim());
            formData.append("image", image);
            return formData;
          })()
        : {
            post_id: Number(postId),
            title: title.trim(),
            content: content.trim(),
          };

      const response = await createComment(payload);
      const comment = response?.data;
      if (comment?.id) {
        setTitle("");
        setContent("");
        setImage(null);
        setResetKey((prev) => prev + 1);
        onCreated?.(comment);
      } else {
        setError("Unexpected response from the server.");
      }
    } catch (err) {
      setError(err?.message || "Could not add comment.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form className="commentForm" onSubmit={handleSubmit}>
      {error && <div className="postsError">{error}</div>}

      <div className="field">
        <label htmlFor="comment-title">Title</label>
        <input
          id="comment-title"
          value={title}
          onChange={(event) => setTitle(event.target.value)}
          placeholder="Comment title"
          required
        />
      </div>

      <div className="field">
        <label htmlFor="comment-content">Content</label>
        <textarea
          id="comment-content"
          value={content}
          onChange={(event) => setContent(event.target.value)}
          placeholder="Write your comment here..."
          required
        />
      </div>

      <div className="commentFormRow">
        <ImageUploadButton
          key={resetKey}
          label="Image (optional)"
          value={image}
          onChange={setImage}
        />
        <button className="button commentFormButton" type="submit" disabled={submitting}>
          {submitting ? "Adding..." : "Add comment"}
        </button>
      </div>
    </form>
  );
}