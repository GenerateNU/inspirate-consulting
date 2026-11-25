import { Files } from "lucide-react";

import { Card, CardContent } from "@/components/ui/card";
import CardIconHeader from "../shared/card-icon-header";
import { panelClassName } from "../shared/media-styles";

interface DescriptionCardProps {
  description: string;
}

// read-only "Description" card under the student's video
export default function DescriptionCard({ description }: DescriptionCardProps) {
  return (
    <Card className={panelClassName}>
      <CardIconHeader icon={Files} title="Description" />
      <CardContent className="px-5">
        <p className="min-h-[106px] rounded-lg border border-greyish-brown px-3 py-3 font-host text-xs whitespace-pre-wrap text-charcoal">
          {description || "No description."}
        </p>
      </CardContent>
    </Card>
  );
}
