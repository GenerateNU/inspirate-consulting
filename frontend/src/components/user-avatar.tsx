import { cn } from "cn";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

export function initials(name: string) {
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) return "?";
  return parts
    .slice(0, 2)
    .map((p) => p[0])
    .join("")
    .toUpperCase();
}

interface UserAvatarProps {
  name: string;
  src?: string | null;
  className?: string;
  fallbackClassName?: string;
}

// Shows the user's picture when one is available, otherwise their initials.
export default function UserAvatar({
  name,
  src,
  className,
  fallbackClassName,
}: UserAvatarProps) {
  return (
    <Avatar aria-hidden className={className}>
      {src && <AvatarImage src={src} alt="" />}
      <AvatarFallback className={cn("font-medium", fallbackClassName)}>
        {initials(name)}
      </AvatarFallback>
    </Avatar>
  );
}
