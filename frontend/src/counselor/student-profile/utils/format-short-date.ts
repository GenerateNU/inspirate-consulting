// e.g. "Sep 23", as in the hi-fi's "Due Sep 23" and "Completed Sep 20"
export function formatShortDate(value: string) {
  return new Date(value).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
  });
}
