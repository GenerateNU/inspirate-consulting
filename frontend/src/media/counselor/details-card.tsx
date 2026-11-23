import { useId } from "react";

import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { cn } from "@/lib/utils";
import { fieldClassName, labelClassName, panelClassName } from "./media-styles";
import type { MediaEditForm } from "./hooks/use-media-edit-form";
import { countWords, descriptionWordLimit } from "./utils/word-count";

interface DetailsCardProps {
  form: MediaEditForm;
}

// the "Details" heading and card: module title and description
export default function DetailsCard({ form }: DetailsCardProps) {
  const titleId = useId();
  const descriptionId = useId();
  const words = countWords(form.description);

  return (
    <div className="flex flex-col gap-2">
      <h2 className="font-jakarta text-lg leading-[normal] text-charcoal">
        Details
      </h2>
      <Card className={panelClassName}>
        <CardContent className="flex flex-col gap-6 px-5">
          <div className="flex flex-col gap-2">
            <Label htmlFor={titleId} className={labelClassName}>
              Module Title
            </Label>
            <Input
              id={titleId}
              value={form.title}
              placeholder="Name"
              required
              disabled={form.saving}
              onChange={(e) => form.setTitle(e.target.value)}
              className={cn(fieldClassName, "h-auto")}
            />
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor={descriptionId} className={labelClassName}>
              Description
            </Label>
            <Textarea
              id={descriptionId}
              value={form.description}
              placeholder="What will students learn? What should they do next? What do they need to know?"
              disabled={form.saving}
              aria-invalid={words > descriptionWordLimit}
              onChange={(e) => form.setDescription(e.target.value)}
              className={cn(fieldClassName, "h-[150px] resize-none")}
            />
            <span
              className={cn(
                "self-end font-host text-[11px] leading-[normal]",
                words > descriptionWordLimit
                  ? "text-danger"
                  : "text-dark-greyish-brown",
              )}
            >
              {words} / {descriptionWordLimit} words
            </span>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
