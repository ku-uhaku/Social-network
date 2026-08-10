import Link from "next/link";
import { resolveMediaSrc, formatDate } from "@/lib/utils";
import Avatar from "@/components/shared/Avatar";

const trimLength = 180;

export default function PostCard({ post, isFeed = true }) {
  const content = post.content || "";
  const snippet = isFeed && content.length > trimLength ? `${content.slice(0, trimLength).trim()}...` : content;
  const createdAt = formatDate(post.created_at);
  const imageSrc = resolveMediaSrc(post.image_url);
  const author = post.user || {};
  const authorName = author.username || "Unknown";

  const cardContent = (
    <>
      <div className="postCardHeader">
        <div className="postAuthorInfo">
          <Avatar avatar={author.avatar} username={authorName} size={40} />
          <div>
            <div className="postAuthorName">{authorName}</div>
            <div className="postMeta">
              {post.privacy} · {createdAt} · {post.comments_count} comments
            </div>
          </div>
        </div>

        <img className="postCardHeaderDecoration" src="/images/post_corner_decoration.png" alt="" />
      </div>

      <h2 className="postTitle">{post.title}</h2>
      <img className="postTitleSeparator" src="/images/post_title_separator.png" alt="" />

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
