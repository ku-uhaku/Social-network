// app/page.tsx
import HollowKnightInput from "@/components/HollowKnightInput";
import HollowKnightCheckbox from "@/components/HollowKnightCheckbox";
import HollowKnightSelect from "@/components/HollowKnightSelect";
import HollowKnightButton from "@/components/HollowKnightButton";

export default function SettingsPage() {
  return (
    <main className="min-h-screen bg-hk-abyss flex items-center justify-center p-6 text-hk-bone">
      <div className="w-full max-w-md bg-hk-surface/40 border border-hk-border rounded-xl p-8 backdrop-blur-sm shadow-2xl flex flex-col gap-6">
        
        {/* Title */}
        <div className="text-center">
          <h2 className="text-2xl font-serif tracking-[0.2em] text-hk-bone uppercase">
            Knight Profile
          </h2>
          <div className="w-24 h-[1px] bg-gradient-to-r from-transparent via-hk-border to-transparent my-2 mx-auto" />
        </div>

        {/* Inputs */}
        <HollowKnightInput label="Vessel Name" placeholder="e.g. Ghost" />

        <HollowKnightSelect
          label="Nail Upgrade Level"
          options={[
            { value: "old", label: "Old Nail" },
            { value: "sharpened", label: "Sharpened Nail" },
            { value: "chiseled", label: "Chiseled Nail" },
            { value: "pure", label: "Pure Nail" },
          ]}
        />

        {/* Checkboxes */}
        <div className="flex flex-col gap-2 pt-2">
          <HollowKnightCheckbox label="Enable Wayward Compass" />
          <HollowKnightCheckbox label="Show Infection Glow" />
        </div>

        {/* Submit Action */}
        <div className="pt-4 flex justify-center">
          <HollowKnightButton>Save Settings</HollowKnightButton>
        </div>

      </div>
    </main>
  );
}