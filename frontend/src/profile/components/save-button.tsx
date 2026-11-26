import { useEffect, useLayoutEffect, useRef } from "react";
import { cn } from "@/lib/utils";

export type SaveState = "idle" | "saving" | "saved";

const CLOSE_MS = 420;
// Check draw (280ms) plus its pop (200ms delay + 360ms), from index.css
const CHECK_MS = 560;
const SHRUNK_PX = 32;
const ARC_START = 26; // percent of the circle the spinning arc covers
const SPIN_DEG_PER_MS = 0.42;

const prefersReducedMotion = () =>
  window.matchMedia("(prefers-reduced-motion: reduce)").matches;

// Submit button that shrinks into a spinning ring while saving, closes the
// ring on success, draws a check, then grows back into the button.
export function SaveButton({
  state,
  disabled,
  onSavedShown,
}: {
  state: SaveState;
  disabled: boolean;
  onSavedShown: () => void;
}) {
  const buttonRef = useRef<HTMLButtonElement>(null);
  const arcRef = useRef<SVGCircleElement>(null);
  const checkRef = useRef<SVGPathElement>(null);
  const fullWidth = useRef(0);
  const anim = useRef({ angle: -90, len: ARC_START, raf: 0 });
  const onSavedShownRef = useRef(onSavedShown);
  onSavedShownRef.current = onSavedShown;

  const paint = () => {
    const { angle, len } = anim.current;
    arcRef.current?.setAttribute("stroke-dasharray", `${len} 100`);
    arcRef.current?.setAttribute("transform", `rotate(${angle} 12 12)`);
  };

  // Shrink to the ring when saving starts; grow back when returning to idle.
  useLayoutEffect(() => {
    const button = buttonRef.current;
    if (!button) return;
    if (state === "saving") {
      fullWidth.current = button.offsetWidth;
      button.style.width = `${fullWidth.current}px`;
      void button.offsetWidth;
      button.style.width = `${SHRUNK_PX}px`;
    } else if (state === "idle" && button.style.width) {
      button.style.width = `${fullWidth.current}px`;
      const clear = () => (button.style.width = "");
      button.addEventListener("transitionend", clear, { once: true });
      return () => button.removeEventListener("transitionend", clear);
    }
  }, [state]);

  // Drive the arc: spin while saving, then close it and draw the check.
  useEffect(() => {
    const a = anim.current;
    const check = checkRef.current;
    const reduce = prefersReducedMotion();
    cancelAnimationFrame(a.raf);

    if (state === "idle") {
      a.angle = -90;
      a.len = ARC_START;
      check?.classList.remove("save-check-drawn");
      paint();
      return;
    }

    if (state === "saving") {
      a.angle = -90;
      a.len = ARC_START;
      check?.classList.remove("save-check-drawn");
      paint();
      if (reduce) return;
      let last = performance.now();
      const spin = (t: number) => {
        a.angle += SPIN_DEG_PER_MS * (t - last);
        last = t;
        paint();
        a.raf = requestAnimationFrame(spin);
      };
      a.raf = requestAnimationFrame(spin);
      return () => cancelAnimationFrame(a.raf);
    }

    // state === "saved"
    let holdTimer = 0;
    const showCheck = () => {
      check?.classList.add("save-check-drawn");
      holdTimer = window.setTimeout(
        () => onSavedShownRef.current(),
        reduce ? 0 : CHECK_MS,
      );
    };
    if (reduce) {
      a.len = 100;
      paint();
      showCheck();
    } else {
      const start = performance.now();
      const startLen = a.len;
      let last = start;
      const close = (t: number) => {
        const p = Math.min(1, (t - start) / CLOSE_MS);
        const eased = 1 - (1 - p) ** 3;
        a.len = startLen + (100 - startLen) * eased;
        a.angle += SPIN_DEG_PER_MS * (1 - eased) * (t - last);
        last = t;
        paint();
        if (p < 1) a.raf = requestAnimationFrame(close);
        else showCheck();
      };
      a.raf = requestAnimationFrame(close);
    }
    return () => {
      cancelAnimationFrame(a.raf);
      clearTimeout(holdTimer);
    };
  }, [state]);

  const busy = state !== "idle";

  return (
    <button
      ref={buttonRef}
      type="submit"
      data-slot="save"
      disabled={disabled || busy}
      aria-label={busy ? "Saving changes" : undefined}
      className={cn(
        "relative flex h-8 items-center justify-center overflow-hidden rounded-lg bg-dark-green px-4 font-host text-xs leading-[normal] whitespace-nowrap text-off-white outline-none",
        "transition-[width,padding,border-radius,opacity,background-color] duration-300 ease-[cubic-bezier(0.4,0,0.2,1)] motion-reduce:transition-none",
        "hover:bg-dark-green/90 focus-visible:ring-[3px] focus-visible:ring-dark-green/30",
        busy
          ? "cursor-progress rounded-2xl px-0"
          : "disabled:cursor-not-allowed disabled:opacity-50",
      )}
    >
      <span
        className={cn("transition-opacity duration-150", busy && "opacity-0")}
      >
        Save Changes
      </span>
      <svg
        aria-hidden="true"
        viewBox="0 0 24 24"
        className={cn(
          "absolute inset-0 m-auto size-[18px] overflow-visible opacity-0 transition-opacity duration-150",
          busy && "opacity-100",
        )}
      >
        <circle
          ref={arcRef}
          cx="12"
          cy="12"
          r="10"
          fill="none"
          stroke="currentColor"
          strokeWidth="2.25"
          strokeLinecap="round"
          pathLength={100}
          strokeDasharray={`${ARC_START} 100`}
          transform="rotate(-90 12 12)"
        />
        <path
          ref={checkRef}
          className="save-check"
          d="M8 12.4l2.7 2.7 5.3-5.7"
          fill="none"
          stroke="currentColor"
          strokeWidth="2.25"
          strokeLinecap="round"
          strokeLinejoin="round"
          pathLength={100}
        />
      </svg>
    </button>
  );
}
