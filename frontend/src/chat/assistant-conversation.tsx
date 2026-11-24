import {
  useState,
  type FormEvent,
  type KeyboardEvent,
  type ReactNode,
} from "react";
import { Link } from "react-router-dom";
import { ArrowLeft, Navigation, Sparkles } from "lucide-react";
import { useAskLlm } from "../api/endpoints";
import { errorMessage } from "../api/client";
import { cn } from "../lib/utils";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Textarea } from "@/components/ui/textarea";

// Lives only in this component; the assistant conversation isn't saved.
interface AssistantMessage {
  id: number;
  sent: boolean;
  text: string;
}

export default function AssistantConversation() {
  const [messages, setMessages] = useState<AssistantMessage[]>([]);
  const askLlm = useAskLlm();

  const addMessage = (sent: boolean, text: string) =>
    setMessages((prev) => [...prev, { id: prev.length, sent, text }]);

  const ask = (question: string) => {
    addMessage(true, question);
    askLlm.mutate(
      { data: { input: question } },
      { onSuccess: (res) => addMessage(false, res.answer) },
    );
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
        <AssistantAvatar className="size-12" />
        <h2 className="truncate text-2xl font-semibold">Assistant</h2>
      </header>

      <div className="flex min-h-0 flex-1 flex-col-reverse overflow-y-auto px-6 md:px-10">
        <div className="flex flex-col gap-6 py-6">
          {messages.length === 0 && (
            <p className="text-center text-sm text-neutral-500">
              Ask about your students' tasks or messages.
            </p>
          )}

          {messages.map((m) =>
            m.sent ? (
              <div key={m.id} className="flex flex-col items-end">
                <Bubble text={m.text} sent />
              </div>
            ) : (
              <AssistantReply key={m.id}>
                <Bubble text={m.text} sent={false} />
              </AssistantReply>
            ),
          )}

          {askLlm.isPending && (
            <AssistantReply>
              <Skeleton
                aria-label="Assistant is thinking"
                className="h-12 w-56 rounded-lg bg-neutral-200"
              />
            </AssistantReply>
          )}
        </div>
      </div>

      <Composer
        onSend={ask}
        isPending={askLlm.isPending}
        error={askLlm.error ? errorMessage(askLlm.error) : null}
      />
    </div>
  );
}

export function AssistantAvatar({ className }: { className?: string }) {
  return (
    <div
      className={cn(
        "flex shrink-0 items-center justify-center rounded-full bg-mist text-neutral-600",
        className,
      )}
    >
      <Sparkles className="size-1/2" />
    </div>
  );
}

function Bubble({ text, sent }: { text: string; sent: boolean }) {
  return (
    <div
      className={cn(
        "max-w-[min(32rem,85%)] rounded-lg px-5 py-3 break-words whitespace-pre-wrap",
        sent
          ? "bg-forest text-white"
          : "border border-neutral-200 bg-white text-neutral-900",
      )}
    >
      {text}
    </div>
  );
}

function AssistantReply({ children }: { children: ReactNode }) {
  return (
    <div className="flex items-end gap-3">
      <AssistantAvatar className="size-10" />
      <div className="flex min-w-0 flex-1 flex-col items-start gap-2">
        <span className="text-sm text-neutral-600">Assistant</span>
        {children}
      </div>
    </div>
  );
}

function Composer({
  onSend,
  isPending,
  error,
}: {
  onSend: (question: string) => void;
  isPending: boolean;
  error: string | null;
}) {
  const [text, setText] = useState("");
  const canSend = text.trim() !== "" && !isPending;

  const send = () => {
    if (!canSend) return;
    onSend(text.trim());
    setText("");
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
          placeholder="Ask the assistant..."
          aria-label="Question"
          className="max-h-40 min-h-0 flex-1 resize-none rounded-xl border-neutral-200 bg-white px-6 py-5 text-base shadow-none md:text-base"
        />
        <Button
          type="submit"
          variant="ghost"
          size="icon-lg"
          disabled={!canSend}
          aria-label="Ask"
          className="mb-3 text-forest hover:bg-transparent hover:text-forest/80 disabled:opacity-40"
        >
          <Navigation className="size-8 rotate-45 fill-current" />
        </Button>
      </div>
      {error && <p className="mt-2 text-sm text-red-600">{error}</p>}
    </form>
  );
}
