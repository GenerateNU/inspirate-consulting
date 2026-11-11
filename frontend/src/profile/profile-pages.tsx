import { Link } from "react-router-dom";
import {
  Bell,
  ChevronRight,
  IdCard,
  LogIn,
  LogOut,
  UserPlus,
  type LucideIcon,
} from "lucide-react";
import { useAuthContext } from "@/contexts/use-auth-context";

type PageOption = {
  label: string;
  description: string;
  icon: LucideIcon;
  to: string;
};

const accountPages: PageOption[] = [
  {
    label: "Notification Preferences",
    description: "Choose which emails you receive",
    icon: Bell,
    to: "/update-notification-preferences",
  },
  {
    label: "Log In",
    description: "Sign in to your account",
    icon: LogIn,
    to: "/login",
  },
  {
    label: "Create Account",
    description: "Set up a new user",
    icon: UserPlus,
    to: "/create-account",
  },
  {
    label: "Log Out",
    description: "Sign out of this device",
    icon: LogOut,
    to: "/logout",
  },
];

const optionClass =
  "group flex w-full items-center gap-4 rounded-xl border bg-card p-4 text-left shadow-xs outline-none transition-[background-color,border-color,box-shadow,transform] duration-200 hover:border-primary/40 hover:bg-accent hover:shadow-sm focus-visible:ring-[3px] focus-visible:ring-ring/50 active:scale-[0.99]";

function OptionContent({
  icon: Icon,
  title,
  description,
}: {
  icon: LucideIcon;
  title: string;
  description: string;
}) {
  return (
    <>
      <span className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-secondary/60 text-secondary-foreground">
        <Icon className="size-6" strokeWidth={1.75} />
      </span>
      <span className="flex min-w-0 flex-1 flex-col">
        <span className="font-medium">{title}</span>
        <span className="text-sm text-muted-foreground">{description}</span>
      </span>
      <ChevronRight className="size-5 shrink-0 text-muted-foreground transition-transform duration-200 group-hover:translate-x-0.5" />
    </>
  );
}

/**
 * Account pages only. Other pages are linked from their sidebar section,
 * matching the Figma sitemap.
 */
export default function ProfilePages() {
  const { studentID } = useAuthContext();

  // The student profile route needs the signed-in student's ID.
  const pages: PageOption[] = studentID
    ? [
        {
          label: "Student Profile",
          description: "View your student profile",
          icon: IdCard,
          to: `/students/${studentID}`,
        },
        ...accountPages,
      ]
    : accountPages;

  return (
    <div className="mx-auto flex w-full max-w-2xl flex-col gap-6 p-6 md:p-10">
      <h1 className="text-3xl font-bold tracking-tight">Profile</h1>
      <ul className="flex flex-col gap-3">
        {pages.map((page) => (
          <li key={page.label}>
            <Link to={page.to} className={optionClass}>
              <OptionContent
                icon={page.icon}
                title={page.label}
                description={page.description}
              />
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
