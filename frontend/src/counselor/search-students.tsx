import { useContext, useState} from "react";
import { useFetchStudents, useAssignCounselor} from "../api/endpoints";
import type {StudentListItem, FetchStudentsParams} from "../api/models";
import { AuthContext } from "../contexts/auth-context";
import { useDebounce } from "./use-debounce";
import { SearchBar } from "../../components/Searchbar";

export default function SearchStudents() {
    const auth = useContext(AuthContext)
    const counselorID = auth?.counselorID as string
    const [Input, setInput] = useState("");
    const params: FetchStudentsParams = {name: Input}
    const {data: response, error} = useFetchStudents(params);
    const {mutate: assignCounselor} = useAssignCounselor();
    useDebounce({Input, setInput})

  let students: StudentListItem[] = [];
  if (response?.students) {
    students = response.students;
  }

  const studentRows = students.map(student => 
    <div className="form-group">
        
        <li key={student.id}>
        {student.name}
         <button onClick={() => assignCounselor({ data: { counselor_id: counselorID, user_id: student.user_id } })}>
            {counselorID === student.counselor_id ? "Student Claimed" : "Claim Student"}
        </button>
        </li>
    </div>
  )

    return (
        <div className="mx-auto flex w-full max-w-[1024px] flex-col px-4 pt-16 pb-16">
        <div className="flex flex-col gap-1">
            <h1 className="font-inter text-2xl leading-8 font-semibold tracking-[-0.6px] text-[#171c20]">
                Claim &amp; Assign Students
            </h1>
            <p className="max-w-[768px] font-inter text-sm leading-5 text-[#45464c]">
                Claim existing students into your dashboard or invite new students.
            </p>
        </div>

        {/* Filter & query console (placeholder: search bar + two dropdowns) */}
        <div className="mt-4 flex items-center gap-4 border border-[#c6c6cd] bg-white p-4">
            <SearchBar
                value={Input}
                onChangeText={setInput}
                placeholder="Filter by student name"
                className="flex-1"
            />
            <div className="h-11 w-[150px] bg-gray-200" />
            <div className="h-11 w-[190px] bg-gray-200" />
        </div>

        {/* Student table (placeholder) */}
        <div className="mt-7 min-h-[326px] border border-[#c6c6cd] bg-gray-100 p-4">
            <ul>{studentRows}</ul>
            {error && <p>Could not load students.</p>}
        </div>

        {/* Invite / add student card (placeholder) */}
        <div className="mt-9 h-[92px] border border-dashed border-[#76777d] bg-gray-100" />

        {/* Footer navigation (placeholder: back + save buttons) */}
        <div className="mt-9 flex items-center justify-between border-t border-[#c6c6cd] pt-4">
            <div className="h-11 w-[170px] bg-gray-200" />
            <div className="h-11 w-[190px] bg-gray-300" />
        </div>
        </div>
        
    )
}