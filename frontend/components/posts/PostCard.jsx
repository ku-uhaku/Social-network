import Link from "next/link";

export default function PostCard({ post, isFeed = true }) {
  const snippet = isFeed && post.content.length > 180 ? `${post.content.slice(0, 180).trim()}...` : post.content;
  const createdAt = new Date(post.created_at).toLocaleDateString();

  const cardContent = (
    <>
      <div className="postCardHeader">
        <div>
          <h2 className="postTitle">{post.title}</h2>
          <div className="postMeta">
            {post.privacy} · {createdAt} · {post.comments_count} comments
          </div>
        </div>
        {!isFeed && (
          <Link href={`/posts/${post.id}`} className="button hk-button">
            View
          </Link>
        )}
      </div>

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