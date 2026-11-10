import { Outlet, useMatch } from "react-router-dom";
import { cn } from "@/lib/utils";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import NavBar from "./nav-bar";

// SidebarProvider persists the open/collapsed state in this cookie.
const readSidebarCookie = () =>
  !document.cookie.split("; ").includes("sidebar_state=false");

export default function Layout() {
  // Chat scrolls inside its own panes, so it needs the page pinned to the viewport.
  const fillViewport = useMatch("/chats/*") !== null;

  return (
    <SidebarProvider defaultOpen={readSidebarCookie()}>
      <NavBar />
      <SidebarInset className={cn(fillViewport && "h-dvh min-h-0")}>
        <header className="flex h-12 items-center border-b px-3 md:hidden">
          <SidebarTrigger />
        </header>
        <div className="flex min-h-0 flex-1 flex-col">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
