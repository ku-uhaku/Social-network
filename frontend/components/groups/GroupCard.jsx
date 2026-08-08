import Link from "next/link";

export default function GroupCard({ group }) {
  const createdAt = group.created_at
    ? new Date(group.created_at).toLocaleDateString() // TODO: investigate where toLocaleDateString is used and see if it can be made into a util func if all backend dates are stored the same way
    : "";

  return (
    <Link href={`/group/${group.id}`} className="groupCard">
      <div className="groupCardHeader">
        <h2 className="groupCardTitle">{group.title}</h2>
      </div>
      <p className="groupCardDescription">{group.description}</p>
      {createdAt && <div className="groupCardMeta">Created {createdAt}</div>}
    </Link>
  );
}
