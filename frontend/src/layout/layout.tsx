import { Outlet } from "react-router-dom";
import NavBar from "./nav-bar";

export default function Layout() {
  return (
    <div className="flex min-h-screen flex-col">
      <NavBar />
      <main className="flex flex-1 flex-col">
        <Outlet />
      </main>
    </div>
  );
}
