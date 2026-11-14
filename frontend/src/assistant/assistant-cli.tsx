import { useEffect, useRef, useState } from "react";
import { errorMessage } from "../api/client";
import { useAskLlm } from "../api/endpoints";

type Line = {
  kind: "user" | "bot" | "error" | "system";
  text: string;
};

const PROMPT = "counselor@inspirate:~$";

const WELCOME: Line[] = [
  {
    kind: "system",
    text: "Inspirate assistant — ask about your students' tasks or chats.",
  },
  { kind: "system", text: 'Type "help" for commands.' },
];

const HELP: Line[] = [
  { kind: "system", text: "help    show this message" },
  { kind: "system", text: "clear   clear the screen" },
  { kind: "system", text: "↑ / ↓   browse previous questions" },
  { kind: "system", text: "Anything else is sent to the assistant." },
];

const lineColor: Record<Line["kind"], string> = {
  user: "text-cream",
  bot: "text-sage",
  error: "text-red-400",
  system: "text-mist/70",
};

/** Terminal-style chat with the LLM assistant. */
export default function AssistantCli() {
  const [lines, setLines] = useState<Line[]>(WELCOME);
  const [input, setInput] = useState("");
  const [history, setHistory] = useState<string[]>([]);
  const [historyIndex, setHistoryIndex] = useState<number | null>(null);
  const { mutate, isPending } = useAskLlm();

  const inputRef = useRef<HTMLInputElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ block: "end" });
  }, [lines, isPending]);

  useEffect(() => {
    if (!isPending) inputRef.current?.focus();
  }, [isPending]);

  const append = (...newLines: Line[]) =>
    setLines((prev) => [...prev, ...newLines]);

  const submit = () => {
    const text = input.trim();
    setInput("");
    setHistoryIndex(null);
    if (!text) return;

    setHistory((prev) => [...prev, text]);

    if (text === "clear") {
      setLines([]);
      return;
    }
    append({ kind: "user", text });
    if (text === "help") {
      append(...HELP);
      return;
    }

    mutate(
      { data: { input: text } },
      {
        onSuccess: (result) =>
          append({ kind: "bot", text: result.answer || "(no response)" }),
        onError: (error) =>
          append({ kind: "error", text: `error: ${errorMessage(error)}` }),
      },
    );
  };

  const handleKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") {
      event.preventDefault();
      submit();
    } else if (event.key === "ArrowUp" && history.length > 0) {
      event.preventDefault();
      const next =
        historyIndex === null
          ? history.length - 1
          : Math.max(0, historyIndex - 1);
      setHistoryIndex(next);
      setInput(history[next]);
    } else if (event.key === "ArrowDown" && historyIndex !== null) {
      event.preventDefault();
      const next = historyIndex + 1;
      if (next >= history.length) {
        setHistoryIndex(null);
        setInput("");
      } else {
        setHistoryIndex(next);
        setInput(history[next]);
      }
    }
  };

  return (
    <div
      className="flex h-dvh flex-col bg-charcoal font-mono text-sm text-cream"
      onClick={() => {
        // Keep focus on the prompt unless the user is selecting text to copy.
        if (!window.getSelection()?.toString()) inputRef.current?.focus();
      }}
    >
      <div className="flex-1 overflow-y-auto px-4 py-4">
        {lines.map((line, i) => (
          <div
            key={i}
            className={`whitespace-pre-wrap break-words leading-relaxed ${lineColor[line.kind]}`}
          >
            {line.kind === "user" && (
              <span className="text-mint">{PROMPT} </span>
            )}
            {line.text}
          </div>
        ))}

        {isPending && (
          <div className="text-mist/70">
            thinking<span className="animate-pulse">...</span>
          </div>
        )}

        {!isPending && (
          <div className="flex items-baseline leading-relaxed">
            <label htmlFor="assistant-input" className="shrink-0 text-mint">
              {PROMPT}&nbsp;
            </label>
            <input
              id="assistant-input"
              ref={inputRef}
              data-slot="terminal-input"
              value={input}
              onChange={(event) => setInput(event.target.value)}
              onKeyDown={handleKeyDown}
              autoComplete="off"
              spellCheck={false}
              className="min-w-0 flex-1 bg-transparent text-cream caret-mint outline-none"
            />
          </div>
        )}

        <div ref={bottomRef} />
      </div>
    </div>
  );
}
