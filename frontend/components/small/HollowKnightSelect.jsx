export default function HollowKnightSelect({ label, options, className = "", ...props }) {
  return (
    <div className={`hk-field ${className}`}>
      {label && <label>{label}</label>}
      <select {...props}>
        {options.map((opt) => (
          <option key={opt.value} value={opt.value}>
            {opt.label}
          </option>
        ))}
      </select>
    </div>
  );
}
