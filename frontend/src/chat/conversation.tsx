import {
  useEffect,
  useMemo,
  useRef,
  useState,
  type FormEvent,
  type KeyboardEvent,
} from "react";
import { Link } from "react-router-dom";
import { useInfiniteQuery, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft, Navigation } from "lucide-react";
import {
  getListChatMessagesQueryKey,
  getListChatsQueryKey,
  listChatMessages,
  useCreateChatMessage,
  useMarkChatRead,
} from "../api/endpoints";
import type { ChatMessage } from "../api/models";
import { errorMessage } from "../api/client";
import { cn } from "../lib/utils";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Textarea } from "@/components/ui/textarea";
import UserAvatar from "@/components/user-avatar";
import {
  POLL_INTERVAL_MS,
  buildTimeline,
  formatSectionHeader,
  formatTime,
  type MessageGroup,
} from "./chat-utils";

const PAGE_SIZE = 30;

interface ConversationProps {
  currentUserId: string | null;
  otherUserId: string;
  otherUserName: string;
}

export default function Conversation({
  currentUserId,
  otherUserId,
  otherUserName,
}: ConversationProps) {
  const queryClient = useQueryClient();
  const scrollerRef = useRef<HTMLDivElement>(null);
  const topSentinelRef = useRef<HTMLDivElement>(null);

  const messagesQuery = useInfiniteQuery({
    queryKey: [...getListChatMessagesQueryKey(otherUserId), "infinite"],
    queryFn: ({ pageParam, signal }) =>
      listChatMessages(
        otherUserId,
        { limit: PAGE_SIZE, cursor: pageParam },
        { signal },
      ),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: (lastPage) => lastPage.next_cursor ?? undefined,
    refetchInterval: POLL_INTERVAL_MS,
  });
  const { hasNextPage, isFetchingNextPage, fetchNextPage } = messagesQuery;

  // Pages arrive newest first; flip to oldest first for display.
  const messages = useMemo(() => {
    const seen = new Set<string>();
    const all: ChatMessage[] = [];
    for (const page of messagesQuery.data?.pages ?? []) {
      for (const m of page.messages ?? []) {
        if (!seen.has(m.id)) {
          seen.add(m.id);
          all.push(m);
        }
      }
    }
    return all.reverse();
  }, [messagesQuery.data]);

  const timeline = useMemo(() => buildTimeline(messages), [messages]);

  const invalidateChat = () => {
    void queryClient.invalidateQueries({
      queryKey: getListChatMessagesQueryKey(otherUserId),
    });
    void queryClient.invalidateQueries({ queryKey: getListChatsQueryKey() });
  };

  // Load older messages when the top of the history scrolls into view.
  useEffect(() => {
    const sentinel = topSentinelRef.current;
    if (!sentinel || !hasNextPage) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting && !isFetchingNextPage) void fetchNextPage();
      },
      { root: scrollerRef.current, rootMargin: "200px 0px 0px 0px" },
    );
    observer.observe(sentinel);
    return () => observer.disconnect();
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  const { mutate: markRead, isPending: isMarkingRead } = useMarkChatRead();
  const hasUnread = messages.some(
    (m) => m.sender_id === otherUserId && m.read_at === null,
  );

  useEffect(() => {
    if (!hasUnread || isMarkingRead) return;
    markRead(
      { userId: otherUserId },
      {
        onSuccess: () => {
          void queryClient.invalidateQueries({
            queryKey: getListChatMessagesQueryKey(otherUserId),
          });
          void queryClient.invalidateQueries({
            queryKey: getListChatsQueryKey(),
          });
        },
      },
    );
  }, [hasUnread, isMarkingRead, markRead, otherUserId, queryClient]);

  const scrollToBottom = () => {
    // The scroller is column-reverse, so 0 is the bottom.
    scrollerRef.current?.scrollTo({ top: 0 });
  };

  return (
    <div className="flex min-h-0 flex-1 flex-col">
      <header className="flex items-center gap-4 border-b border-neutral-300 px-6 py-5 md:px-10">
        <Button
          asChild
          variant="ghost"
          size="icon"
          className="text-neutral-500 md:hidden"
        >
          <Link to="/chats" aria-label="Back to conversations">
            <ArrowLeft className="size-5" />
          </Link>
        </Button>
        <UserAvatar
          name={otherUserName}
          className="size-12"
          fallbackClassName="bg-mist text-sm text-neutral-600"
        />
        <h2 className="truncate text-2xl font-semibold">
          {otherUserName || " "}
        </h2>
      </header>

      <div
        ref={scrollerRef}
        className="flex min-h-0 flex-1 flex-col-reverse overflow-y-auto px-6 md:px-10"
      >
        <div className="flex flex-col gap-6 py-6">
          <div ref={topSentinelRef} />
          {isFetchingNextPage && <MessagesSkeleton rows={2} />}
          {messagesQuery.isLoading && <MessagesSkeleton rows={5} />}
          {messagesQuery.error && (
            <p className="text-center text-sm text-red-600">
              {errorMessage(messagesQuery.error)}
            </p>
          )}
          {messagesQuery.isSuccess && messages.length === 0 && (
            <p className="text-center text-sm text-neutral-500">
              No messages yet. Say hello!
            </p>
          )}

          {timeline.map((item) =>
            item.kind === "separator" ? (
              <p
                key={item.key}
                className="text-center text-sm text-neutral-400"
              >
                {formatSectionHeader(item.at)}
              </p>
            ) : item.group.senderId === currentUserId ? (
              <SentGroup key={item.key} group={item.group} />
            ) : (
              <ReceivedGroup
                key={item.key}
                group={item.group}
                senderName={otherUserName}
              />
            ),
          )}
        </div>
      </div>

      <Composer
        recipientId={otherUserId}
        onSent={() => {
          invalidateChat();
          scrollToBottom();
        }}
      />
    </div>
  );
}

// Placeholder bubbles that alternate sides, like a real conversation.
function MessagesSkeleton({ rows }: { rows: number }) {
  const widths = ["w-48", "w-72", "w-40", "w-64", "w-56"];
  return (
    <div aria-label="Loading messages" className="flex flex-col gap-4">
      {Array.from({ length: rows }, (_, i) =>
        i % 2 === 0 ? (
          <div key={i} className="flex items-end gap-3">
            <Skeleton className="size-10 rounded-full bg-neutral-200" />
            <Skeleton
              className={cn("h-12 rounded-lg bg-neutral-200", widths[i % 5])}
            />
          </div>
        ) : (
          <Skeleton
            key={i}
            className={cn(
              "h-12 self-end rounded-lg bg-neutral-200",
              widths[i % 5],
            )}
          />
        ),
      )}
    </div>
  );
}

function Bubble({ message, sent }: { message: ChatMessage; sent: boolean }) {
  return (
    <div
      className={cn(
        "max-w-[min(32rem,85%)] rounded-lg px-5 py-3 break-words whitespace-pre-wrap",
        sent
          ? "bg-forest text-white"
          : "border border-neutral-200 bg-white text-neutral-900",
      )}
    >
      {message.message}
    </div>
  );
}

function ReceivedGroup({
  group,
  senderName,
}: {
  group: MessageGroup;
  senderName: string;
}) {
  const last = group.messages[group.messages.length - 1];

  return (
    <div className="flex flex-col gap-1">
      <div className="flex items-end gap-3">
        <UserAvatar
          name={senderName}
          className="size-10"
          fallbackClassName="bg-mist text-xs text-neutral-600"
        />
        <div className="flex min-w-0 flex-1 flex-col items-start gap-2">
          <span className="text-sm text-neutral-600">{senderName}</span>
          {group.messages.map((m) => (
            <Bubble key={m.id} message={m} sent={false} />
          ))}
        </div>
      </div>
      <span className="pl-13 text-xs text-neutral-400">
        {formatTime(new Date(last.created_at))}
      </span>
    </div>
  );
}

function SentGroup({ group }: { group: MessageGroup }) {
  const last = group.messages[group.messages.length - 1];

  return (
    <div className="flex flex-col items-end gap-2">
      {group.messages.map((m) => (
        <Bubble key={m.id} message={m} sent />
      ))}
      <span className="text-xs text-neutral-400">
        {last.read_at && <span className="text-neutral-500">Read · </span>}
        {formatTime(new Date(last.created_at))}
      </span>
    </div>
  );
}

function Composer({
  recipientId,
  onSent,
}: {
  recipientId: string;
  onSent: () => void;
}) {
  const [text, setText] = useState("");
  const createMessage = useCreateChatMessage();
  const canSend = text.trim() !== "" && !createMessage.isPending;

  const send = () => {
    if (!canSend) return;
    createMessage.mutate(
      { data: { recipient_id: recipientId, message: text.trim() } },
      {
        onSuccess: () => {
          setText("");
          onSent();
        },
      },
    );
  };

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    send();
  };

  // Enter sends; Shift+Enter inserts a newline. Skip while an IME is composing.
  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey && !e.nativeEvent.isComposing) {
      e.preventDefault();
      send();
    }
  };

  return (
    <form onSubmit={handleSubmit} className="px-6 pt-2 pb-6 md:px-10">
      <div className="flex items-end gap-4">
        <Textarea
          rows={1}
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Your message..."
          aria-label="Message"
          className="max-h-40 min-h-0 flex-1 resize-none rounded-xl border-neutral-200 bg-white px-6 py-5 text-base shadow-none md:text-base"
        />
        <Button
          type="submit"
          variant="ghost"
          size="icon-lg"
          disabled={!canSend}
          aria-label="Send message"
          className="mb-3 text-forest hover:bg-transparent hover:text-forest/80 disabled:opacity-40"
        >
          <Navigation className="size-8 rotate-45 fill-current" />
        </Button>
      </div>
      {createMessage.error && (
        <p className="mt-2 text-sm text-red-600">
          {errorMessage(createMessage.error)}
        </p>
      )}
    </form>
  );
}
