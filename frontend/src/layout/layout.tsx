import { Outlet, useMatch } from "react-router-dom";
import NavBar from "./nav-bar";

export default function Layout() {
  // Chat scrolls inside its own panes, so it needs the page pinned to the viewport.
  const fillViewport = useMatch("/chats/*") !== null;

  return (
    <div className={`flex flex-col ${fillViewport ? "h-dvh" : "min-h-screen"}`}>
      <NavBar />
      <main className="flex min-h-0 flex-1 flex-col">
        <Outlet />
      </main>
    </div>
  );
}
