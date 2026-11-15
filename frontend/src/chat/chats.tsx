import { useParams } from "react-router-dom";
import { useFetchUser, useListChats } from "../api/endpoints";
import { errorMessage } from "../api/client";
import { useAuthContext } from "../contexts/use-auth-context";
import { cn } from "../lib/utils";
import ChatList from "./chat-list";
import Conversation from "./conversation";
import { POLL_INTERVAL_MS } from "./chat-utils";

export default function Chats() {
  const { userId } = useParams();
  const { userID: currentUserId } = useAuthContext();

  const chatsQuery = useListChats({
    query: { refetchInterval: POLL_INTERVAL_MS },
  });
  const chats = chatsQuery.data ?? [];
  const selected = chats.find((c) => c.other_user_id === userId);

  // A brand-new conversation isn't in the chat list yet, so look the name up.
  const userQuery = useFetchUser(userId ?? "", {
    query: { enabled: !!userId && !chatsQuery.isLoading && !selected },
  });
  const otherUserName = selected?.other_user_name ?? userQuery.data?.name ?? "";

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
        {userId ? (
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
