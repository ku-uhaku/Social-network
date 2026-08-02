export default function HollowKnightCheckbox({ label, className = "", ...props }) {
  return (
    <label className={`hk-choice ${className}`}>
      <input type="checkbox" {...props} />
      <span className="hk-choice-box" />
      <span>{label}</span>
    </label>
  );
}
