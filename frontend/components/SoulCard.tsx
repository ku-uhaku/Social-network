import HollowKnightButton from "./HollowKnightButton";

export default function SoulCard() {
  return (
    <div className="max-w-md w-full bg-hk-surface border border-hk-border rounded-xl p-6 shadow-2xl transition-all duration-300 hover:border-hk-soul hover:shadow-[0_0_20px_rgba(0,229,255,0.2)]">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold text-hk-bone tracking-wide">
          Soul Vessel
        </h2>
        <span className="px-2.5 py-1 text-xs font-medium bg-hk-soul/10 text-hk-soul border border-hk-soul/30 rounded-full">
          Full
        </span>
      </div>

      <p className="text-hk-muted text-sm mb-6 leading-relaxed">
        Gather soul by striking enemies with the Nail to focus energy or cast powerful spells.
      </p>


    </div>
  );
}