import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";

/**
 * Asks before leaving while `dirty`. Covers:
 * - `guard(action)`: in-page exits such as a back button
 * - clicks on in-app links (e.g. the sidebar)
 * - closing or reloading the tab (the browser's own prompt)
 *
 * The browser back button is not covered.
 **/
export function useLeaveGuard(dirty: boolean) {
  const navigate = useNavigate();
  const [pending, setPending] = useState<(() => void) | null>(null);
  const dirtyRef = useRef(dirty);
  dirtyRef.current = dirty;

  useEffect(() => {
    if (!dirty) return;
    const onBeforeUnload = (e: BeforeUnloadEvent) => e.preventDefault();
    window.addEventListener("beforeunload", onBeforeUnload);
    return () => window.removeEventListener("beforeunload", onBeforeUnload);
  }, [dirty]);

  // Capture phase runs before React Router's <Link> handler, which skips
  // navigating when the event is already defaultPrevented.
  useEffect(() => {
    const onClick = (e: MouseEvent) => {
      if (!dirtyRef.current || e.defaultPrevented || e.button !== 0) return;
      if (e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) return;
      const link = (e.target as Element | null)?.closest?.("a[href]");
      if (!(link instanceof HTMLAnchorElement)) return;
      if (link.target && link.target !== "_self") return;
      if (link.hasAttribute("download")) return;
      const url = new URL(link.href, window.location.href);
      if (url.origin !== window.location.origin) return;
      const to = url.pathname + url.search + url.hash;
      const here =
        window.location.pathname +
        window.location.search +
        window.location.hash;
      if (to === here) return;
      e.preventDefault();
      setPending(() => () => void navigate(to));
    };
    document.addEventListener("click", onClick, true);
    return () => document.removeEventListener("click", onClick, true);
  }, [navigate]);

  const guard = (action: () => void) => {
    if (dirtyRef.current) setPending(() => action);
    else action();
  };

  return {
    guard,
    dialog: {
      open: pending !== null,
      onOpenChange: (open: boolean) => {
        if (!open) setPending(null);
      },
      onConfirm: () => {
        const action = pending;
        setPending(null);
        action?.();
      },
    },
  };
}
