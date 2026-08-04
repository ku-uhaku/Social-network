"use client";

const TOGGLE_OFF = "/images/toggle_off.png";
const TOGGLE_ON = "/images/toggle_on.png";

export default function CharmToggle({
  checked = false,
  onChange,
  label,
  title,
}) {
  return (
    <button
      type="button"
      className={`charmToggle ${checked ? "on" : "off"}`}
      title={title || label}
      onClick={() => onChange?.(!checked)}
    >
      <img
        className="charmToggleCharm"
        src={TOGGLE_ON}
        alt=""
      />
      <img
        className="charmToggleSwitch"
        src={checked ? TOGGLE_ON : TOGGLE_OFF}
        alt={checked ? "On" : "Off"}
      />
    </button>
  );
}