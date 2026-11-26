// e.g. "07/20/2027", as in the hi-fi
export function formatDueDate(value: string) {
  return new Date(value).toLocaleDateString("en-US", {
    month: "2-digit",
    day: "2-digit",
    year: "numeric",
  });
}
