import Link from "next/link";
import { resolveMediaSrc } from "@/lib/utils";
import Avatar from "@/components/shared/Avatar";

const trimLength = 180;

export default function PostCard({ post, isFeed = true }) {
  const snippet = isFeed && post.content.length > trimLength ? `${post.content.slice(0, trimLength).trim()}...` : post.content;
  const createdAt = new Date(post.created_at).toLocaleDateString();
  const imageSrc = resolveMediaSrc(post.image_url);
  const author = post.user || {};
  const authorName = [author.first_name, author.last_name].filter(Boolean).join(" ") || author.username || "Unknown";

  const cardContent = (
    <>
      <div className="postCardHeader">
        <div className="postAuthorInfo">
          <Avatar avatar={author.avatar} name={authorName} size={40} />
          <div>
            <div className="postAuthorName">{authorName}</div>
            <div className="postMeta">
              {post.privacy} · {createdAt} · {post.comments_count} comments
            </div>
          </div>
        </div>
      </div>

      <h2 className="postTitle">{post.title}</h2>

      {imageSrc && (
        <img className="postCardImage" src={imageSrc} alt="image" />
      )}

      <p className="postSnippet">{snippet}</p>
    </>
  );

  if (isFeed) {
    return (
      <Link href={`/posts/${post.id}`} className="postCard postCardLink">
        {cardContent}
      </Link>
    );
  }

  return <article className="postCard">{cardContent}</article>;
}