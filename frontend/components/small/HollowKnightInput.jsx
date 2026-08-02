export default function HollowKnightInput({ label, className = "", ...props }) {
  return (
    <div className={`hk-field ${className}`}>
      {label && <label>{label}</label>}
      <input {...props} />
    </div>
  );
}
