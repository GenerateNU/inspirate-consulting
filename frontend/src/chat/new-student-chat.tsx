import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";

import {
  getGetStudentQueryKey,
  getStudent,
  useListCounselorStudents,
} from "../api/endpoints";
import type { CounselorStudent } from "../api/models";
import { Input } from "@/components/ui/input";
import { ErrorMessage } from "@/components/ui/error-message";
import UserAvatar from "@/components/user-avatar";

// most matches shown under the search box at once
const MAX_MATCHES = 6;

// counselor's "start a chat": search their own students by name and pick one
export default function NewStudentChat() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const studentsQuery = useListCounselorStudents();
  const [query, setQuery] = useState("");
  const [opening, setOpening] = useState(false);
  const [openError, setOpenError] = useState<unknown>(null);

  const trimmed = query.trim().toLowerCase();
  const matches =
    trimmed === ""
      ? []
      : (studentsQuery.data ?? [])
          .filter((student) => student.name.toLowerCase().includes(trimmed))
          .slice(0, MAX_MATCHES);

  // chats are addressed by user id, which the student list doesn't include, so look it up
  const openChat = async (student: CounselorStudent) => {
    setOpening(true);
    setOpenError(null);
    try {
      const { user_id } = await queryClient.fetchQuery({
        queryKey: getGetStudentQueryKey(student.student_id),
        queryFn: () => getStudent(student.student_id),
      });
      setQuery("");
      void navigate(`/chats/${user_id}`);
    } catch (err) {
      setOpenError(err);
    } finally {
      setOpening(false);
    }
  };

  return (
    <div className="mt-3 flex flex-col gap-1">
      <Input
        value={query}
        onChange={(e) => {
          setQuery(e.target.value);
          setOpenError(null);
        }}
        placeholder="Search your students to start a chat"
        aria-label="Search your students"
        className="bg-white"
      />

      <ErrorMessage
        error={studentsQuery.error}
        prefix="Couldn't load your students"
        size="xs"
      />
      <ErrorMessage
        error={openError}
        prefix="Couldn't open the chat"
        size="xs"
      />

      {trimmed !== "" && !studentsQuery.isLoading && (
        <ul className="mt-1 flex flex-col rounded-md border border-neutral-200 bg-white">
          {matches.length === 0 ? (
            <li className="px-3 py-2 text-sm text-neutral-500">
              No matching students.
            </li>
          ) : (
            matches.map((student) => (
              <li key={student.student_id}>
                <button
                  type="button"
                  data-slot="student-option"
                  disabled={opening}
                  onClick={() => void openChat(student)}
                  className="flex w-full items-center gap-2 px-3 py-2 text-left text-sm hover:bg-neutral-100 disabled:opacity-50"
                >
                  <UserAvatar
                    name={student.name}
                    className="size-6"
                    fallbackClassName="bg-mist text-[10px] text-neutral-600"
                  />
                  <span className="truncate">{student.name}</span>
                </button>
              </li>
            ))
          )}
        </ul>
      )}
    </div>
  );
}
