import Link from "next/link";
import { formatDate } from "@/lib/utils";

export default function GroupCard({ group, unreadCount = 0 }) {
  const createdAt = formatDate(group.created_at);

  return (
    <Link href={`/group/${group.id}`} className="groupCard">
      <div className="groupCardHeader">
        <h2 className="groupCardTitle">{group.title}</h2>
        {unreadCount > 0 && <span className="groupChatBadge">{unreadCount}</span>}
      </div>
      <p className="groupCardDescription">{group.description}</p>
      {createdAt && <div className="groupCardMeta">Created {createdAt}</div>}
    </Link>
  );
}
