"use client";

import { resolveMediaSrc, formatDate } from "@/lib/utils";
import Avatar from "@/components/shared/Avatar";

export default function CommentCard({ comment }) {
  const author = comment.user || {};
  const name = author.username || "Unknown";
  const createdAt = formatDate(comment.created_at);
  const imageSrc = resolveMediaSrc(comment.image_url);

  return (
    <article className="commentCard">
      <div className="commentCardHeader">
        <Avatar avatar={author.avatar} username={name} size={40} />
        <div className="commentCardMeta">
          <span className="commentAuthor">{name}</span>
          <span className="commentDate">{createdAt}</span>
        </div>
      </div>

      <h3 className="commentTitle">{comment.title}</h3>
      <p className="commentContent">{comment.content}</p>

      {imageSrc && (
        <img className="commentImage" src={imageSrc} alt="comment image" />
      )}
    </article>
  );
}