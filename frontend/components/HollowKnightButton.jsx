export default function HollowKnightButton({ children, className = "", ...props }) {
  return (
    <button className={`hk-button ${className}`} {...props}>
      {children}
    </button>
  );
}
