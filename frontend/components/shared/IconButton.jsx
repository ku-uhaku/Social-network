import Link from "next/link";
import "@/css/iconButton.css";

export default function IconButton({ icon, label, href, onClick, active, children }) {
  const className = `iconButton ${active ? "active" : ""}`;
  const content = (
    <>
      <img className="iconButtonIcon" src={icon} />
      <span>{label}</span>
      {children}
    </>
  );

  if (href) {
    return (
      <Link href={href} className={className} title={label}>
        {content}
      </Link>
    );
  }

  return (
    <button type="button" className={className} title={label} onClick={onClick}>
      {content}
    </button>
  );
}
