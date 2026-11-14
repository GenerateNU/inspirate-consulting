import { useParams } from "react-router-dom";

import { useGetGlobalCollege } from "../../api/endpoints";

// loads the college whose id is in the page's URL (/counselor-college-list/:id)
export function useCollegeFromRoute() {
  const id = Number(useParams().id);
  return useGetGlobalCollege(id);
}
