import {
  useEffect,
  useRef,
  useState,
  type FormEvent,
  type KeyboardEvent,
} from "react";
import { Link } from "react-router-dom";
import { ArrowLeft, Bot, Navigation } from "lucide-react";
import { useAskLlm } from "../api/endpoints";
import { errorMessage } from "../api/client";
import { cn } from "../lib/utils";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { formatTime } from "./chat-utils";

type AssistantMessage = {
  id: number;
  fromUser: boolean;
  text: string;
  at: Date;
  error?: boolean;
};

const GREETING =
  "Hi! Ask me about your students' tasks or chats, like “find chats where a student talks about hating literature class”.";

export function AssistantAvatar({ className }: { className?: string }) {
  return (
    <span
      aria-hidden
      className={cn(
        "flex shrink-0 items-center justify-center rounded-full bg-charcoal text-mint",
        className,
      )}
    >
      <Bot className="size-1/2" />
    </span>
  );
}

/** The AI assistant, shown as its own chat in the conversation pane. */
export default function AssistantConversation({ canUse }: { canUse: boolean }) {
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
        <AssistantAvatar className="size-12" />
        <div className="min-w-0">
          <h2 className="truncate text-2xl font-semibold">Assistant</h2>
          <p className="truncate text-sm text-neutral-500">
            Searches your students' tasks and chats
          </p>
        </div>
      </header>

      {canUse ? (
        <AssistantChat />
      ) : (
        <div className="flex flex-1 items-center justify-center p-6 text-neutral-500">
          The assistant is only available to counselors.
        </div>
      )}
    </div>
  );
}

function AssistantChat() {
  const [messages, setMessages] = useState<AssistantMessage[]>(() => [
    { id: 0, fromUser: false, text: GREETING, at: new Date() },
  ]);
  const [text, setText] = useState("");
  const nextId = useRef(1);
  const bottomRef = useRef<HTMLDivElement>(null);
  const askLlm = useAskLlm();
  const canSend = text.trim() !== "" && !askLlm.isPending;

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ block: "end" });
  }, [messages, askLlm.isPending]);

  const append = (message: Omit<AssistantMessage, "id" | "at">) =>
    setMessages((prev) => [
      ...prev,
      { ...message, id: nextId.current++, at: new Date() },
    ]);

  const send = () => {
    if (!canSend) return;
    const question = text.trim();
    setText("");
    append({ fromUser: true, text: question });
    askLlm.mutate(
      { data: { input: question } },
      {
        onSuccess: (result) =>
          append({
            fromUser: false,
            text: result.answer || "I couldn't find anything for that.",
          }),
        onError: (error) =>
          append({ fromUser: false, text: errorMessage(error), error: true }),
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
    <>
      <div className="min-h-0 flex-1 overflow-y-auto px-6 md:px-10">
        <div className="flex flex-col gap-6 py-6">
          {messages.map((m) =>
            m.fromUser ? (
              <div key={m.id} className="flex flex-col items-end gap-2">
                <Bubble sent>{m.text}</Bubble>
                <span className="text-xs text-neutral-400">
                  {formatTime(m.at)}
                </span>
              </div>
            ) : (
              <ReceivedMessage key={m.id} at={m.at} error={m.error}>
                {m.text}
              </ReceivedMessage>
            ),
          )}
          {askLlm.isPending && (
            <ReceivedMessage>
              <span className="text-neutral-500">
                Searching<span className="animate-pulse">...</span>
              </span>
            </ReceivedMessage>
          )}
          <div ref={bottomRef} />
        </div>
      </div>

      <form onSubmit={handleSubmit} className="px-6 pt-2 pb-6 md:px-10">
        <div className="flex items-end gap-4">
          <Textarea
            rows={1}
            value={text}
            onChange={(e) => setText(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Ask the assistant..."
            aria-label="Message the assistant"
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
      </form>
    </>
  );
}

function Bubble({
  sent,
  error,
  children,
}: {
  sent: boolean;
  error?: boolean;
  children: React.ReactNode;
}) {
  return (
    <div
      className={cn(
        "max-w-[min(32rem,85%)] rounded-lg px-5 py-3 break-words whitespace-pre-wrap",
        sent
          ? "bg-forest text-white"
          : error
            ? "border border-red-200 bg-red-50 text-red-700"
            : "border border-neutral-200 bg-white text-neutral-900",
      )}
    >
      {children}
    </div>
  );
}

function ReceivedMessage({
  at,
  error,
  children,
}: {
  at?: Date;
  error?: boolean;
  children: React.ReactNode;
}) {
  return (
    <div className="flex flex-col gap-1">
      <div className="flex items-end gap-3">
        <AssistantAvatar className="size-10" />
        <div className="flex min-w-0 flex-1 flex-col items-start gap-2">
          <span className="text-sm text-neutral-600">Assistant</span>
          <Bubble sent={false} error={error}>
            {children}
          </Bubble>
        </div>
      </div>
      {at && (
        <span className="pl-13 text-xs text-neutral-400">{formatTime(at)}</span>
      )}
    </div>
  );
}
