// most words a module description can have
export const descriptionWordLimit = 250;

export function countWords(text: string) {
  return text.trim() === "" ? 0 : text.trim().split(/\s+/).length;
}
