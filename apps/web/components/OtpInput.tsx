import { ChangeEvent, useRef, KeyboardEvent, ClipboardEvent } from "react";

interface Props {
    value: string;
    onChange: (value: string) => void;
}

const LENGTH = 6;

export default function OtpInput({ value, onChange }: Props) {
    const refs = useRef<(HTMLInputElement | null)[]>([]);

    const digits = value.padEnd(LENGTH, "").split("").slice(0, LENGTH);
 
    const handleChange = (index: number, e: ChangeEvent<HTMLInputElement>) => {
        const char = e.target.value.replace(/\D/g, "").slice(-1);
        const next = digits.map((d, i) => (i === index ? char : d)).join("").trimEnd();
        onChange(next);
        if (char && index < LENGTH - 1) {
        refs.current[index + 1]?.focus();
        }
    };

    const handleKeyDown = (index: number, e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Backspace") {
          if (digits[index]) {
            const next = digits.map((d, i) => (i === index ? "" : d)).join("").trimEnd();
            onChange(next);
          } else if (index > 0) {
            refs.current[index - 1]?.focus();
            const next = digits.map((d, i) => (i === index - 1 ? "" : d)).join("").trimEnd();
            onChange(next);
          }
        }
        if (e.key === "ArrowLeft" && index > 0) refs.current[index - 1]?.focus();
        if (e.key === "ArrowRight" && index < LENGTH - 1) refs.current[index + 1]?.focus();
      };
     
      const handlePaste = (e: ClipboardEvent<HTMLInputElement>) => {
        e.preventDefault();
        const pasted = e.clipboardData.getData("text").replace(/\D/g, "").slice(0, LENGTH);
        onChange(pasted);
        refs.current[Math.min(pasted.length, LENGTH - 1)]?.focus();
      };
    return (
        <div className="flex gap-2 justify-center">
        {Array.from({ length: LENGTH }).map((_, i) => (
          <input
            key={i}
            ref={(el) => { refs.current[i] = el; }}
            type="text"
            inputMode="numeric"
            pattern="[0-9]*"
            maxLength={1}
            value={digits[i] ?? ""}
            onChange={(e) => handleChange(i, e)}
            onKeyDown={(e) => handleKeyDown(i, e)}
            onPaste={handlePaste}
            onFocus={(e) => e.target.select()}
            className="w-12 h-14 text-center text-xl font-bold border-2 border-gray-300 rounded-lg
                       focus:outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-100
                       transition-colors"
            aria-label={`Digit ${i + 1}`}
          />
        ))}
      </div> 
    )
}
