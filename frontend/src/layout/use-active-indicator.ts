import { useLayoutEffect, useState } from "react";

type Indicator = { top: number; height: number; visible: boolean };

/**
 * Sliding highlight that sits behind the active nav item. It is measured
 * against the container so a single pill can travel between the main menu
 * and the footer.
 *
 * Nav items mark themselves with `data-nav-path={to}` so the active one can
 * be found by its path.
 *
 * The container is tracked with a callback ref because the mobile sheet
 * mounts a fresh one each time it opens, and the effect must re-run then.
 */
export function useActiveIndicator(activePath: string | undefined) {
  const [container, setContainer] = useState<HTMLDivElement | null>(null);
  const [indicator, setIndicator] = useState<Indicator>({
    top: 0,
    height: 0,
    visible: false,
  });
  // Skip the transition on the very first measurement so the pill doesn't
  // slide in from the top of the sidebar on page load.
  const [ready, setReady] = useState(false);

  useLayoutEffect(() => {
    if (!container) return;

    const measure = () => {
      const active = activePath
        ? container.querySelector<HTMLElement>(
            `[data-nav-path="${activePath}"]`,
          )
        : null;
      if (!active) {
        setIndicator((prev) => ({ ...prev, visible: false }));
        return;
      }
      const containerRect = container.getBoundingClientRect();
      const activeRect = active.getBoundingClientRect();
      setIndicator({
        top: activeRect.top - containerRect.top,
        height: activeRect.height,
        visible: true,
      });
    };

    measure();
    const frame = requestAnimationFrame(() => setReady(true));
    const observer = new ResizeObserver(measure);
    observer.observe(container);
    return () => {
      cancelAnimationFrame(frame);
      observer.disconnect();
    };
  }, [container, activePath]);

  return { containerRef: setContainer, indicator, ready };
}
