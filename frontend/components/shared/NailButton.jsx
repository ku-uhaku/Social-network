"use client";

const BUTTON_HOVER_DECORATION = "/images/button_hover_decoration.png";

export default function NailButton({
  children,
  className = "",
  type = "button",
  onClick,
  disabled = false,
  title,
}) {
  return (
    <button
      type={type}
      className={`nailButton ${className}`}
      onClick={onClick}
      disabled={disabled}
      title={title}
    >
      <img
        className="nailButtonDecoration nailButtonDecorationLeft"
        src={BUTTON_HOVER_DECORATION}
        alt=""
      />
      <span className="nailButtonLabel">{children}</span>
      <img
        className="nailButtonDecoration nailButtonDecorationRight"
        src={BUTTON_HOVER_DECORATION}
        alt=""
      />
    </button>
  );
}