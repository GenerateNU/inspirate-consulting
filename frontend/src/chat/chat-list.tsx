import { useState, type FormEvent } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import type { ChatSummary } from "../api/models";
import { cn } from "../lib/utils";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import UserAvatar from "@/components/user-avatar";
import { formatListTime, isUUID } from "./chat-utils";

interface ChatListProps {
  chats: ChatSummary[];
  currentUserId: string | null;
  isLoading: boolean;
  error: string | null;
}

export default function ChatList({
  chats,
  currentUserId,
  isLoading,
  error,
}: ChatListProps) {
  return (
    <div className="flex min-h-0 flex-1 flex-col">
      <div className="border-b border-neutral-200 px-5 py-5">
        <h1 className="text-2xl font-semibold">Messages</h1>
        <NewChatForm currentUserId={currentUserId} />
      </div>

      <ul className="min-h-0 flex-1 overflow-y-auto">
        {isLoading &&
          Array.from({ length: 4 }, (_, i) => (
            <li
              key={i}
              aria-label="Loading conversations"
              className="flex items-center gap-3 px-5 py-3"
            >
              <Skeleton className="size-10 rounded-full bg-neutral-200" />
              <div className="flex flex-1 flex-col gap-2">
                <Skeleton className="h-4 w-1/2 bg-neutral-200" />
                <Skeleton className="h-3 w-4/5 bg-neutral-200" />
              </div>
            </li>
          ))}
        {error && <li className="p-5 text-sm text-red-600">{error}</li>}
        {!isLoading && !error && chats.length === 0 && (
          <li className="p-5 text-sm text-neutral-500">
            No conversations yet.
          </li>
        )}
        {chats.map((chat) => (
          <li key={chat.other_user_id}>
            <ChatListItem chat={chat} currentUserId={currentUserId} />
          </li>
        ))}
      </ul>
    </div>
  );
}

function ChatListItem({
  chat,
  currentUserId,
}: {
  chat: ChatSummary;
  currentUserId: string | null;
}) {
  const unread = chat.unread_count > 0;
  const preview =
    chat.last_sender_id === currentUserId
      ? `You: ${chat.last_message}`
      : chat.last_message;

  return (
    <NavLink
      to={`/chats/${chat.other_user_id}`}
      className={({ isActive }) =>
        cn(
          "flex items-center gap-3 px-5 py-3 transition-colors hover:bg-neutral-100",
          isActive && "bg-neutral-100",
        )
      }
    >
      <UserAvatar
        name={chat.other_user_name}
        className="size-10"
        fallbackClassName="bg-mist text-xs text-neutral-600"
      />
      <div className="min-w-0 flex-1">
        <div className="flex items-baseline justify-between gap-2">
          <span className={cn("truncate", unread && "font-semibold")}>
            {chat.other_user_name}
          </span>
          <span className="shrink-0 text-xs text-neutral-400">
            {formatListTime(new Date(chat.last_message_at))}
          </span>
        </div>
        <div className="flex items-center justify-between gap-2">
          <span
            className={cn(
              "truncate text-sm",
              unread ? "text-neutral-900" : "text-neutral-500",
            )}
          >
            {preview}
          </span>
          {unread && (
            <Badge className="bg-forest text-white">{chat.unread_count}</Badge>
          )}
        </div>
      </div>
    </NavLink>
  );
}

function NewChatForm({ currentUserId }: { currentUserId: string | null }) {
  const navigate = useNavigate();
  const [userId, setUserId] = useState("");
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const id = userId.trim();
    if (!isUUID(id)) {
      setError("Enter a valid user ID.");
      return;
    }
    if (id === currentUserId) {
      setError("You can't message yourself.");
      return;
    }
    setError(null);
    setUserId("");
    void navigate(`/chats/${id}`);
  };

  return (
    <form onSubmit={handleSubmit} className="mt-3 flex flex-col gap-1">
      <div className="flex gap-2">
        <Input
          value={userId}
          onChange={(e) => {
            setUserId(e.target.value);
            setError(null);
          }}
          placeholder="Paste a user ID to start a chat"
          aria-label="User ID"
          aria-invalid={error !== null}
          className="bg-white"
        />
        <Button type="submit" variant="outline">
          Start
        </Button>
      </div>
      {error && <p className="text-xs text-red-600">{error}</p>}
    </form>
  );
}
