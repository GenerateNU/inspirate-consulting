import { NavLink, matchPath, useLocation } from "react-router-dom";
import {
  CircleUserRound,
  FileText,
  ListTodo,
  MessageCircle,
  Palette,
  PanelLeftClose,
  PanelLeftOpen,
  School,
  SquarePlay,
  type LucideIcon,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
  SidebarRail,
  useSidebar,
} from "@/components/ui/sidebar";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import inspirateLogo from "@/assets/inspirate-logo.svg";
import inspirateIcon from "@/assets/inspirate-icon.svg";
import { useActiveIndicator } from "./use-active-indicator";

type NavItemConfig = {
  to: string;
  label: string;
  icon: LucideIcon;
};

const mainLinks: NavItemConfig[] = [
  { to: "/student-tasks", label: "Task Dashboard", icon: ListTodo },
  { to: "/college-applications", label: "Colleges", icon: School },
  { to: "/essay", label: "Essays", icon: FileText },
  { to: "/extracurriculars", label: "Extracurriculars", icon: Palette },
  { to: "/videos", label: "Modules", icon: SquarePlay },
  { to: "/chats", label: "Chats", icon: MessageCircle },
];

const footerLinks: NavItemConfig[] = [
  { to: "/profile", label: "Profile", icon: CircleUserRound },
];

const isActivePath = (to: string, pathname: string) =>
  matchPath({ path: to, end: false }, pathname) !== null;

function NavItem({
  item: { to, label, icon: Icon },
  active,
  showTooltip,
  onNavigate,
}: {
  item: NavItemConfig;
  active: boolean;
  showTooltip: boolean;
  onNavigate: () => void;
}) {
  return (
    <SidebarMenuItem>
      {/* disableHoverableContent: close as soon as the pointer leaves the
          trigger. Otherwise Radix waits for the (unmounted, while expanded)
          content to close it, and stale tooltips pop up on collapse. */}
      <Tooltip disableHoverableContent>
        <TooltipTrigger asChild>
          <NavLink
            to={to}
            onClick={onNavigate}
            data-nav-path={to}
            aria-current={active ? "page" : undefined}
            className={cn(
              "relative z-10 flex h-11 items-center gap-3 overflow-hidden rounded-md px-2.5 text-base outline-hidden ring-sidebar-ring",
              "transition-[padding,background-color] duration-300 ease-in-out focus-visible:ring-2",
              "group-data-[collapsible=icon]:px-4",
              active
                ? "text-sidebar-accent-foreground"
                : "text-sidebar-foreground hover:bg-white/10",
            )}
          >
            <Icon
              className={cn(
                "size-6 shrink-0 transition-transform duration-300 ease-out",
                active && "scale-105",
              )}
              strokeWidth={1.75}
            />
            {/* Fade out fast on collapse, fade in slower on expand */}
            <span className="whitespace-nowrap transition-opacity duration-300 ease-in-out group-data-[collapsible=icon]:opacity-0 group-data-[collapsible=icon]:duration-150">
              {label}
            </span>
          </NavLink>
        </TooltipTrigger>
        {showTooltip && <TooltipContent side="right">{label}</TooltipContent>}
      </Tooltip>
    </SidebarMenuItem>
  );
}

export default function NavBar() {
  const { state, isMobile, setOpenMobile, toggleSidebar } = useSidebar();
  const { pathname } = useLocation();
  const collapsed = state === "collapsed" && !isMobile;

  const activePath = [...mainLinks, ...footerLinks].find((item) =>
    isActivePath(item.to, pathname),
  )?.to;
  const { containerRef, indicator, ready } = useActiveIndicator(activePath);

  const closeMobile = () => {
    if (isMobile) setOpenMobile(false);
  };

  const renderItems = (items: NavItemConfig[]) =>
    items.map((item) => (
      <NavItem
        key={item.to}
        item={item}
        active={item.to === activePath}
        showTooltip={collapsed}
        onNavigate={closeMobile}
      />
    ));

  return (
    <Sidebar collapsible="icon" className="border-sidebar-border">
      <SidebarHeader className="relative h-24 shrink-0 overflow-hidden p-0">
        {/* Expanded: full wordmark + collapse button */}
        <NavLink
          to="/"
          onClick={closeMobile}
          tabIndex={collapsed ? -1 : undefined}
          aria-hidden={collapsed}
          className="absolute top-6 left-5 transition-opacity duration-300 ease-in-out group-data-[collapsible=icon]:pointer-events-none group-data-[collapsible=icon]:opacity-0 group-data-[collapsible=icon]:duration-150"
        >
          <img src={inspirateLogo} alt="Inspirate" className="h-9 w-auto" />
        </NavLink>
        <Button
          variant="ghost"
          size="icon"
          onClick={toggleSidebar}
          aria-label="Collapse sidebar"
          tabIndex={collapsed ? -1 : undefined}
          className="absolute top-6 right-4 text-sidebar-foreground transition-opacity duration-300 hover:bg-white/10 hover:text-sidebar-foreground group-data-[collapsible=icon]:pointer-events-none group-data-[collapsible=icon]:opacity-0 group-data-[collapsible=icon]:duration-150 [&_svg]:size-6!"
        >
          <PanelLeftClose strokeWidth={1.75} />
        </Button>

        {/* Collapsed: bracket icon that turns into an expand button on hover */}
        <button
          type="button"
          data-slot="sidebar-expand"
          onClick={toggleSidebar}
          aria-label="Expand sidebar"
          tabIndex={collapsed ? undefined : -1}
          aria-hidden={!collapsed}
          className="group/expand pointer-events-none absolute top-5 left-[18px] flex size-11 items-center justify-center rounded-md text-sidebar-foreground opacity-0 ring-sidebar-ring outline-hidden transition-[opacity,background-color] duration-150 hover:bg-white/10 focus-visible:ring-2 group-data-[collapsible=icon]:pointer-events-auto group-data-[collapsible=icon]:opacity-100 group-data-[collapsible=icon]:delay-100 group-data-[collapsible=icon]:duration-300"
        >
          <img
            src={inspirateIcon}
            alt=""
            className="absolute h-10 w-auto rounded-md transition-[opacity,transform] duration-200 group-hover/expand:scale-90 group-hover/expand:opacity-0 group-focus-visible/expand:opacity-0"
          />
          <PanelLeftOpen
            className="size-6 scale-90 opacity-0 transition-[opacity,transform] duration-200 group-hover/expand:scale-100 group-hover/expand:opacity-100 group-focus-visible/expand:opacity-100"
            strokeWidth={1.75}
          />
        </button>
      </SidebarHeader>

      <div ref={containerRef} className="relative flex min-h-0 flex-1 flex-col">
        <div
          aria-hidden
          className={cn(
            "pointer-events-none absolute inset-x-3 top-0 rounded-md bg-sidebar-accent shadow-md shadow-black/30",
            ready &&
              "transition-[transform,height,opacity] duration-300 ease-[cubic-bezier(0.22,1,0.36,1)] motion-reduce:transition-none",
          )}
          style={{
            height: indicator.height,
            transform: `translateY(${indicator.top}px)`,
            opacity: indicator.visible ? 1 : 0,
          }}
        />

        <SidebarContent className="px-3 pt-6">
          <SidebarMenu className="gap-5">{renderItems(mainLinks)}</SidebarMenu>
        </SidebarContent>

        <SidebarFooter className="px-3 pt-3 pb-6">
          <SidebarMenu>{renderItems(footerLinks)}</SidebarMenu>
        </SidebarFooter>
      </div>

      <SidebarRail />
    </Sidebar>
  );
}
