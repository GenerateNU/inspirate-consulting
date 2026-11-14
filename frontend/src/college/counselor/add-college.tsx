import { useNavigate } from "react-router-dom";

import { SchoolEditor, SchoolPageFrame } from "./college-details";

// "Add College": the same form as editing a school, starting empty
export default function AddCollege() {
  const navigate = useNavigate();
  return (
    <SchoolPageFrame editing>
      <SchoolEditor
        college={null}
        onSaved={(saved) =>
          navigate(`/counselor-college-list/${saved.id}`, { replace: true })
        }
        onDiscard={() => navigate("/counselor-college-list")}
      />
    </SchoolPageFrame>
  );
}
