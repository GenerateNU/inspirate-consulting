import { useNavigate } from "react-router-dom";

import MediaEditor from "./media-editor";
import MediaPageFrame from "./media-page-frame";

// "Create Module": the same form as editing a module, starting empty
export default function AddMedia() {
  const navigate = useNavigate();
  return (
    <MediaPageFrame>
      <MediaEditor media={null} onSaved={() => navigate(-1)} />
    </MediaPageFrame>
  );
}
