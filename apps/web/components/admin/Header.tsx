"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { ChartIcon, DashboardIcon, PeopleIcon, PlusIcon } from "../icons";

const NAV = [
  {
    href: "/admin",
    label: "Dashboard",
    icon: <DashboardIcon />,
  },
  {
    href: "/admin/candidates",
    label: "Candidates",
    icon: <PeopleIcon />,
  },
  {
    href: "/admin/candidates/new",
    label: "Add candidate",
    icon: <PlusIcon />,
  },
  {
    href: "/results",
    label: "Public results",
    icon: <ChartIcon />,
  },
];

interface Props {
  onNavigate?: () => void;
}

export default function AdminHeader({ onNavigate }: Props) {
  const pathname = usePathname();

  return (
    <nav className="flex-1 space-y-1 p-3">
      {NAV.map(({ href, label, icon }) => {
        const active = pathname === href;

        return (
          <Link
            href={href}
            key={href}
            onClick={onNavigate}
            className={`flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors ${
              active ? "bg-brand-500 text-white" : "text-slate-400 hover:bg-white/5 hover:text-white"
            }`}
          >
            {icon}
            {label}
          </Link>
        );
      })}
    </nav>
  );
}
