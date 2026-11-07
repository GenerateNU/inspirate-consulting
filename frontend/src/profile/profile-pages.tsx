import { Link, useSearchParams } from "react-router-dom";
import {
  ArrowLeft,
  Bell,
  ChevronRight,
  CircleUserRound,
  ClipboardCheck,
  ClipboardList,
  FileCheck,
  FilePlus,
  GraduationCap,
  Hourglass,
  IdCard,
  ImageUp,
  LayoutGrid,
  ListTodo,
  LogIn,
  LogOut,
  Palette,
  School,
  SquarePlay,
  Upload,
  UserCog,
  UserPlus,
  UserSearch,
  Users,
  Video,
  type LucideIcon,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { useAuthContext } from "@/contexts/use-auth-context";

type PageOption = {
  label: string;
  description: string;
  icon: LucideIcon;
  /** Static path, or a builder for routes that need the student's ID. */
  to: string | ((studentId: string) => string);
};

type Section = {
  id: string;
  title: string;
  description: string;
  icon: LucideIcon;
  pages: PageOption[];
};

const sections: Section[] = [
  {
    id: "account",
    title: "Account",
    description: "Sign in, sign out and notification settings",
    icon: CircleUserRound,
    pages: [
      {
        label: "Log In",
        description: "Sign in to your account",
        icon: LogIn,
        to: "/login",
      },
      {
        label: "Log Out",
        description: "Sign out of this device",
        icon: LogOut,
        to: "/logout",
      },
      {
        label: "Create Account",
        description: "Set up a new user",
        icon: UserPlus,
        to: "/create-account",
      },
      {
        label: "Notification Preferences",
        description: "Choose which emails you receive",
        icon: Bell,
        to: "/update-notification-preferences",
      },
    ],
  },
  {
    id: "student",
    title: "Student",
    description: "Your profile, tasks, reviews and videos",
    icon: GraduationCap,
    pages: [
      {
        label: "Student Profile",
        description: "View your student profile",
        icon: IdCard,
        to: (id) => `/students/${id}`,
      },
      {
        label: "Essay Reviews",
        description: "Request and track essay reviews",
        icon: FileCheck,
        to: "/student",
      },
      {
        label: "My Tasks",
        description: "Your to-do list",
        icon: ListTodo,
        to: "/student-tasks",
      },
      {
        label: "My Videos",
        description: "Videos shared with you",
        icon: SquarePlay,
        to: "/student-media-list",
      },
    ],
  },
  {
    id: "counselor",
    title: "Counselor",
    description: "Students, reviews and videos you manage",
    icon: UserCog,
    pages: [
      {
        label: "Complete Reviews",
        description: "Finish reviews and assign review balances",
        icon: ClipboardCheck,
        to: "/counselor",
      },
      {
        label: "My Students",
        description: "Students assigned to you",
        icon: Users,
        to: "/counselor-students",
      },
      {
        label: "Pending Reviews",
        description: "Essays waiting on your review",
        icon: Hourglass,
        to: "/pending-reviews",
      },
      {
        label: "Claim Students",
        description: "Find and assign students",
        icon: UserSearch,
        to: "/search",
      },
      {
        label: "All Videos",
        description: "Every uploaded video",
        icon: Video,
        to: "/counselor-media-list",
      },
    ],
  },
  {
    id: "more",
    title: "More",
    description: "Colleges, applications and uploads",
    icon: LayoutGrid,
    pages: [
      {
        label: "Add College",
        description: "Add a college to the catalog",
        icon: School,
        to: "/create-college",
      },
      {
        label: "College Applications",
        description: "Track your applications",
        icon: ClipboardList,
        to: "/college-applications",
      },
      {
        label: "New Application",
        description: "Start a college application",
        icon: FilePlus,
        to: "/create-application",
      },
      {
        label: "New Extracurricular",
        description: "Add an activity",
        icon: Palette,
        to: "/create-extracurricular",
      },
      {
        label: "Upload Video",
        description: "Add a new module video",
        icon: Upload,
        to: "/upload-video",
      },
      {
        label: "Upload Media",
        description: "Upload a file or video",
        icon: ImageUp,
        to: "/upload-media",
      },
    ],
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

export default function ProfilePages() {
  const { studentID } = useAuthContext();
  const [searchParams, setSearchParams] = useSearchParams();
  const section = sections.find((s) => s.id === searchParams.get("section"));


  // Pages that need a student ID only appear when one is available.
  const pages = section?.pages.filter(
    (page) => typeof page.to === "string" || studentID,
  );

  return (
    <div className="mx-auto flex w-full max-w-2xl flex-col gap-6 p-6 md:p-10">
      <div className="flex items-center gap-3">
        {section && (
          <Button
            variant="ghost"
            size="icon"
            aria-label="Back"
            onClick={() => setSearchParams({})}
          >
            <ArrowLeft className="size-5" />
          </Button>
        )}
        <h1 className="text-3xl font-bold tracking-tight">
          {section?.title ?? "Profile"}
        </h1>
      </div>

      <ul className="flex flex-col gap-3">
        {section
          ? pages?.map((page) => (
              <li key={page.label}>
                <Link
                  to={
                    typeof page.to === "string" ? page.to : page.to(studentID!)
                  }
                  className={optionClass}
                >
                  <OptionContent
                    icon={page.icon}
                    title={page.label}
                    description={page.description}
                  />
                </Link>
              </li>
            ))
          : sections.map((s) => (
              <li key={s.id}>
                <button
                  type="button"
                  data-slot="profile-section"
                  onClick={() => setSearchParams({ section: s.id })}
                  className={optionClass}
                >
                  <OptionContent
                    icon={s.icon}
                    title={s.title}
                    description={s.description}
                  />
                </button>
              </li>
            ))}
      </ul>
    </div>
  );
}
