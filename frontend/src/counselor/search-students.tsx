import { useContext, useEffect, useState } from "react";
import { useFetchStudents, useAssignCounselor} from "../api/endpoints";
import type { FetchStudentsParams, StudentListItem } from "../api/models";
import { AuthContext } from "../contexts/auth-context";

export default function SearchStudents() {

    type Timer = ReturnType<typeof setTimeout>;
    const [Input, setInput] = useState("");
    const [timer, setTimer] = useState<Timer>();
    const auth = useContext(AuthContext)
    const counselorID = auth?.counselorID as string

    const delay: number = 8000;
    const params: FetchStudentsParams = {name: Input}
    
    const {data: response, error} = useFetchStudents(params);
    const {mutate: assignCounselor} = useAssignCounselor();


    useEffect(() => {
        const newTimer = setTimeout(() => {
            setInput(Input)
        }, delay)
    
        clearTimeout(timer);
        setTimer(newTimer);
    }, [Input]);


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
        <div>
        <div className="form-group">
            <h2>Claim and Assign Students</h2>
        </div>
        <div className="form-group">
        <input type="search" value={Input} name="student-search" placeholder="Filter by name" onChange={(e) => setInput(e.target.value)}/>
        </div>

        <ul>{studentRows}</ul>
        {error && <p>Could not load students.</p>}
        
        </div>
        
    )
}