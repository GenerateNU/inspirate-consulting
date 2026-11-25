import { cn } from "cn";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { usePfpSrc } from "@/lib/pfp-cache";

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
  userId?: string;
  pfpKey?: string | null;
  className?: string;
  fallbackClassName?: string;
}

// Shows the user's picture when one is available, otherwise their initials.
export default function UserAvatar({
  name,
  src,
  userId,
  pfpKey,
  className,
  fallbackClassName,
}: UserAvatarProps) {
  const cached = usePfpSrc(userId, pfpKey, src);
  if (userId) src = cached;
  return (
    <Avatar aria-hidden className={className}>
      {src && <AvatarImage src={src} alt="" />}
      <AvatarFallback className={cn("font-medium", fallbackClassName)}>
        {initials(name)}
      </AvatarFallback>
    </Avatar>
  );
}
