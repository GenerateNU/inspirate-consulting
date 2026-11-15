import type { ChatMessage } from "../api/models";

export const POLL_INTERVAL_MS = 5000;

// Consecutive messages from one sender within this window share a bubble group.
const GROUP_GAP_MS = 5 * 60 * 1000;
// A gap this long (or a new day) starts a new section with a timestamp header.
const SECTION_GAP_MS = 60 * 60 * 1000;

export interface MessageGroup {
  senderId: string;
  messages: ChatMessage[];
}

export type TimelineItem =
  | { kind: "separator"; key: string; at: Date }
  | { kind: "group"; key: string; group: MessageGroup };

const isSameDay = (a: Date, b: Date) =>
  a.getFullYear() === b.getFullYear() &&
  a.getMonth() === b.getMonth() &&
  a.getDate() === b.getDate();

// Expects messages oldest first.
export function buildTimeline(messages: ChatMessage[]): TimelineItem[] {
  const items: TimelineItem[] = [];
  let group: MessageGroup | null = null;
  let prev: Date | null = null;

  for (const message of messages) {
    const at = new Date(message.created_at);
    const gap = prev ? at.getTime() - prev.getTime() : Infinity;

    if (!prev || gap > SECTION_GAP_MS || !isSameDay(prev, at)) {
      items.push({ kind: "separator", key: `sep-${message.id}`, at });
      group = null;
    }
    if (!group || group.senderId !== message.sender_id || gap > GROUP_GAP_MS) {
      group = { senderId: message.sender_id, messages: [] };
      items.push({ kind: "group", key: message.id, group });
    }
    group.messages.push(message);
    prev = at;
  }

  return items;
}

export const formatTime = (date: Date) =>
  date.toLocaleTimeString([], { hour: "numeric", minute: "2-digit" });

export function formatSectionHeader(date: Date) {
  const now = new Date();
  const yesterday = new Date(now);
  yesterday.setDate(now.getDate() - 1);

  if (isSameDay(date, now)) return `Today ${formatTime(date)}`;
  if (isSameDay(date, yesterday)) return `Yesterday ${formatTime(date)}`;

  const day = date.toLocaleDateString([], {
    month: "short",
    day: "numeric",
    year: date.getFullYear() === now.getFullYear() ? undefined : "numeric",
  });
  return `${day}, ${formatTime(date)}`;
}

export function formatListTime(date: Date) {
  const now = new Date();
  if (isSameDay(date, now)) return formatTime(date);

  const weekAgo = new Date(now);
  weekAgo.setDate(now.getDate() - 6);
  if (date > weekAgo) return date.toLocaleDateString([], { weekday: "short" });

  return date.toLocaleDateString([], { month: "short", day: "numeric" });
}

const UUID_PATTERN =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

export const isUUID = (value: string) => UUID_PATTERN.test(value);
