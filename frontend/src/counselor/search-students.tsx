import { useContext, useState} from "react";
import type {StudentListItem, FetchStudentsParams} from "../api/models";
import { AuthContext } from "../contexts/auth-context";
import { useDebounce } from "../../hooks/use-debounce";
import { ClaimStudentSearchbar, type Grade } from "./claim-student-searchbar";
import { ClaimStudentTable } from "./claim-student-table";
import { Link } from "react-router-dom";
import { ArrowLeft, ArrowRight, UserPlus } from "lucide-react";
import { useInfiniteStudents } from "./useInfiniteStudents";

// Not built yet, placeholder paths
const EDIT_PROFILE_PATH = "/counselor/profile";
const DASHBOARD_PATH = "/counselor/dashboard";

export default function SearchStudents() {
    const auth = useContext(AuthContext)
    const counselorID = auth?.counselorID as string
    const [Input, setInput] = useState("");

    const [grade, setGrade] = useState<Grade | "">("");
    const params: FetchStudentsParams = {name: Input}
    const {students: loadedStudents, error, hasMore, isLoadingMore, loadMore} = useInfiniteStudents(params);
    
    useDebounce({Input, setInput})

  let students: StudentListItem[] = loadedStudents;
  // only narrow by grade when one is picked; "All Grades" ("") keeps every result
  if (grade !== "") {
    students = students.filter((student) => student.year === grade);
  }

    return (
        <div className="mx-auto flex w-full max-w-[1056px] flex-col px-4 pt-16 pb-16">
        <div className="flex flex-col gap-1">
            <h1 className="font-inter text-2xl leading-8 font-semibold tracking-[-0.6px] text-[#171c20]">
                Claim &amp; Assign Students
            </h1>
            <p className="max-w-[768px] font-inter text-sm leading-5 text-[#45464c]">
                Claim existing students into your dashboard or invite new students.
            </p>
        </div>

        <div className="mt-4">
            <ClaimStudentSearchbar
                search={Input}
                onSearchChange={setInput}
                grade={grade}
                onGradeChange={setGrade}
            />
        </div>

        <div className="mt-7">
            <ClaimStudentTable
                input={Input}
                students={students}
                counselorID={counselorID}
                error={error}
                hasMore={hasMore}
                isLoadingMore={isLoadingMore}
                onLoadMore={loadMore}
            />
        </div>

        <div className="mt-9 flex items-center justify-between gap-4 border border-dashed border-[#76777d] bg-[rgba(240,244,250,0.5)] p-[25px]">
            <div className="flex items-start gap-4">
                <div className="flex size-10 shrink-0 items-center justify-center border border-[#c6c6cd] bg-white">
                    <UserPlus size={20} strokeWidth={1.5} aria-hidden="true" className="text-[#45464c]" />
                </div>
                <div className="flex flex-col">
                    <p className="font-inter text-base leading-6 font-semibold text-[#171c20]">
                        Don&apos;t see a student in the directory?
                    </p>
                    <p className="font-inter text-sm leading-5 text-[#45464c]">
                        Invite a student to the Inspirate Consulting Hub!
                    </p>
                </div>
            </div>
            {/* TODO: open the Send Student Account Invitation page once it exists (Figma node 5142:1603) */}
            <button
                type="button"
                className="inline-flex h-11 shrink-0 items-center justify-center rounded-none border border-[#76777d] bg-white px-[17px] font-inter text-sm leading-5 font-medium tracking-[0.14px] text-[#171c20]"
            >
                + Add New Student / Send Invite
            </button>
        </div>


        <div className="mt-9 flex items-center justify-between border-t border-[#c6c6cd] pt-[17px]">
            <Link
                to={EDIT_PROFILE_PATH}
                className="inline-flex h-11 items-center justify-center gap-2 border border-[#76777d] bg-white px-[17px] font-inter text-sm leading-5 font-medium tracking-[0.14px] text-[#171c20]"
            >
                <ArrowLeft size={12} aria-hidden="true" />
                Back / Edit Profile
            </Link>
            <Link
                to={DASHBOARD_PATH}
                className="inline-flex h-11 items-center justify-center gap-2 border border-black bg-black px-[25px] font-inter text-sm leading-5 font-medium tracking-[0.14px] text-white"
            >
                Save &amp; Proceed to Dashboard
                <ArrowRight size={12} aria-hidden="true" />
            </Link>
        </div>
        </div>
        
    )
}