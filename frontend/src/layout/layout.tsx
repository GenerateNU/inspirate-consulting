import { Outlet } from "react-router-dom";
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
  return (
    <SidebarProvider defaultOpen={readSidebarCookie()}>
      <NavBar />
      <SidebarInset>
        <header className="flex h-12 items-center border-b px-3 md:hidden">
          <SidebarTrigger />
        </header>
        <div className="flex flex-1 flex-col">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
