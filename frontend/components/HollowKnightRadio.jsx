export default function HollowKnightRadio({ label, className = "", ...props }) {
  return (
    <label className={`hk-choice ${className}`}>
      <input type="radio" {...props} />
      <span className="hk-choice-box round" />
      <span>{label}</span>
    </label>
  );
}
