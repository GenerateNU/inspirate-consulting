import { NavLink } from "react-router-dom";

const links = [
  { to: "/essay", label: "Essays" },
  { to: "/college-list", label: "Colleges" },
  { to: "/create-college", label: "Add College" },
  { to: "/colleges", label: "Applications" },
  { to: "/create-application", label: "New Application" },
  { to: "/tasks", label: "Tasks" },
  { to: "/create-task", label: "New Task" },
  { to: "/videos", label: "Videos" },
  { to: "/upload-video", label: "Upload Video" },
  { to: "/student", label: "Student" },
  { to: "/counselor", label: "Counselor" },
];

export default function NavBar() {
  return (
    <nav className="sticky top-0 z-10 flex items-center gap-8 bg-white px-6 py-3 shadow-sm">
      <NavLink
        to="/"
        className="shrink-0 text-xl font-bold tracking-tight text-blue-700"
      >
        Inspirate
      </NavLink>
      <ul className="flex flex-row items-center gap-1 overflow-x-auto whitespace-nowrap">
        {links.map(({ to, label }) => (
          <li key={to}>
            <NavLink
              to={to}
              className={({ isActive }) =>
                `block rounded-md px-3 py-2 text-sm font-medium transition-colors ${
                  isActive
                    ? "bg-blue-50 text-blue-700"
                    : "text-gray-600 hover:bg-gray-100 hover:text-gray-900"
                }`
              }
            >
              {label}
            </NavLink>
          </li>
        ))}
      </ul>
    </nav>
  );
}
