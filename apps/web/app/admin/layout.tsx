"use client";

import { ReactNode, useEffect, useState, useSyncExternalStore } from "react";
import { usePathname, useRouter } from "next/navigation";
import AdminHeader from "@/components/admin/Header";
import { clearAdminToken, getAdminToken } from "@/lib/admin-api/api";

function subscribeToHydration() {
  return () => undefined;
}

function getClientHydrationSnapshot() {
  return true;
}

function getServerHydrationSnapshot() {
  return false;
}

export default function AdminLayout({ children }: { children: ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const [mobileOpen, setMobileOpen] = useState(false);
  const isHydrated = useSyncExternalStore(
    subscribeToHydration,
    getClientHydrationSnapshot,
    getServerHydrationSnapshot,
  );

  const isLoginPage = pathname === "/admin/login";
  const isAuthenticated = isHydrated ? Boolean(getAdminToken()) : false;

  // Auth guard runs after hydration because the token is stored in sessionStorage.
  useEffect(() => {
    if (isHydrated && !isLoginPage && !isAuthenticated) {
      router.replace("/admin/login");
    }
  }, [isHydrated, isAuthenticated, isLoginPage, router]);

  // Login page: no sidebar, plain background
  if (isLoginPage) {
    return <div className="fixed inset-0 z-50 overflow-auto bg-gray-50">{children}</div>;
  }

  // Block render until token is confirmed (prevents flash of sidebar for unauthenticated users)
  if (!isHydrated || !isAuthenticated) {
    return <div className="fixed inset-0 z-50 bg-slate-950" />;
  }

  return (
    <div className="fixed inset-0 z-50 flex flex-col bg-slate-950 text-slate-100">
      {/* Mobile top bar */}
      <div className="flex shrink-0 items-center justify-between border-b border-white/10 bg-slate-950 px-4 py-3 md:hidden">
        <div>
          <div className="text-base font-bold text-white">Olu</div>
          <div className="text-[10px] uppercase tracking-[0.24em] text-slate-500">Admin</div>
        </div>
        <button
          onClick={() => setMobileOpen((open) => !open)}
          className="rounded-lg p-2 text-slate-400 hover:bg-white/10 hover:text-white"
          aria-label="Toggle navigation"
        >
          {mobileOpen ? (
            <svg width="20" height="20" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
              <path d="M6 18L18 6M6 6l12 12" />
            </svg>
          ) : (
            <svg width="20" height="20" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
              <path d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          )}
        </button>
      </div>

      {/* Body */}
      <div className="flex flex-1 overflow-hidden">
        {/* Mobile overlay backdrop */}
        {mobileOpen && (
          <div
            className="fixed inset-0 z-40 bg-black/60 md:hidden"
            onClick={() => setMobileOpen(false)}
          />
        )}

        {/* Sidebar */}
        <aside
          className={`
            fixed inset-y-0 left-0 z-50 flex w-64 flex-col border-r border-white/10 bg-slate-950/95 transition-transform duration-200
            md:relative md:translate-x-0 md:flex
            ${mobileOpen ? "translate-x-0" : "-translate-x-full"}
          `}
        >
          <div className="border-b border-white/10 px-5 py-5">
            <div className="text-lg font-bold text-white">Olu</div>
            <div className="mt-1 text-xs uppercase tracking-[0.24em] text-slate-500">Admin panel</div>
          </div>

          <AdminHeader onNavigate={() => setMobileOpen(false)} />

          <div className="border-t border-white/10 p-3">
            <button
              onClick={() => {
                clearAdminToken();
                router.push("/admin/login");
              }}
              className="flex w-full items-center gap-2 rounded-xl px-3 py-2 text-left text-sm text-slate-400 transition-colors hover:bg-white/5 hover:text-white"
            >
              <svg width="16" height="16" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
                <path d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              Sign out
            </button>
          </div>
        </aside>

        {/* Main content */}
        <main className="flex-1 overflow-y-auto bg-slate-50 px-4 py-6 text-gray-950 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-5xl">{children}</div>
        </main>
      </div>
    </div>
  );
}
