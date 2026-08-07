"use client";


const emojies = (function() {
  const start = 0x1f600; // see https://unicode.org/emoji/charts/full-emoji-list.html
  const offset = 20;
  const list = [];

  for (let i = 0; i < offset; i++) {
    const codepoint = start + i;
    const emoji = String.fromCodePoint(codepoint);
    list.push(emoji);
  }
  return list;
})()

export default function EmojiPicker({ onPick }) {
  return (
    <div className="chatEmojiTray">
      {emojies.map((e, i) => (
        <button key={i} type="button" className="chatEmoji" onClick={() => onPick(e)}>
          {e}
        </button>
      ))}
    </div>
  );
}
