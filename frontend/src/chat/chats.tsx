import { Navigate, useParams } from "react-router-dom";
import { useFetchUser, useListChats } from "../api/endpoints";
import { CurrentUserRole } from "../api/models";
import { errorMessage } from "../api/client";
import { useAuthContext } from "../contexts/use-auth-context";
import { cn } from "../lib/utils";
import AssistantConversation from "./assistant-conversation";
import ChatList from "./chat-list";
import Conversation from "./conversation";
import {
  AGENT_ENABLED,
  ASSISTANT_CHAT_ID,
  POLL_INTERVAL_MS,
} from "./chat-utils";

export default function Chats() {
  const { userId } = useParams();
  const { userID: currentUserId, role } = useAuthContext();
  // The assistant is counselor-only and behind the AGENT_ENABLED flag.
  const showAssistant = AGENT_ENABLED && role === CurrentUserRole.counselor;
  const isAssistant = userId === ASSISTANT_CHAT_ID;

  const chatsQuery = useListChats({
    query: { refetchInterval: POLL_INTERVAL_MS },
  });
  const chats = chatsQuery.data ?? [];
  const selected = chats.find((c) => c.other_user_id === userId);

  // A brand-new conversation isn't in the chat list yet, so look the name up.
  const userQuery = useFetchUser(userId ?? "", {
    query: {
      enabled: !!userId && !isAssistant && !chatsQuery.isLoading && !selected,
    },
  });
  const otherUserName = selected?.other_user_name ?? userQuery.data?.name ?? "";

  if (isAssistant && !showAssistant) return <Navigate to="/chats" replace />;

  return (
    <div className="flex min-h-0 flex-1 bg-paper">
      <aside
        className={cn(
          "min-h-0 w-full flex-col border-r border-neutral-200 md:flex md:w-80 md:shrink-0",
          userId ? "hidden" : "flex",
        )}
      >
        <ChatList
          chats={chats}
          currentUserId={currentUserId}
          showAssistant={showAssistant}
          isLoading={chatsQuery.isLoading}
          error={chatsQuery.error ? errorMessage(chatsQuery.error) : null}
        />
      </aside>

      <section
        className={cn(
          "min-h-0 min-w-0 flex-1 flex-col md:flex",
          userId ? "flex" : "hidden",
        )}
      >
        {isAssistant ? (
          <AssistantConversation />
        ) : userId ? (
          <Conversation
            key={userId}
            currentUserId={currentUserId}
            otherUserId={userId}
            otherUserName={otherUserName}
          />
        ) : (
          <div className="flex flex-1 items-center justify-center text-neutral-500">
            Select a conversation to start messaging.
          </div>
        )}
      </section>
    </div>
  );
}
