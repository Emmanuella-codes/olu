"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { getAdminStats } from "@/lib/admin-api/api";
import { getHealth } from "@/lib/api/api";
import { buildStatusRows, statCards } from "@/lib/data/admin";
import { AdminStats } from "@/types/types";

export default function AdminDashboard() {
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [health, setHealth] = useState<{ status: string; error?: string } | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    Promise.all([getAdminStats(), getHealth()])
      .then(([statsData, healthData]) => {
        setStats(statsData);
        setHealth(healthData);
      })
      .catch(() => setError(true))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-8">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.2em] text-brand-600">Admin</p>
          <h1 className="mt-2 text-3xl font-bold text-gray-950">Dashboard</h1>
        </div>
        <Link href="/results" className="btn-secondary h-auto min-h-0 px-4 py-2">
          View public results ↗
        </Link>
      </div>

      {error && (
        <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
          Could not load stats. Check the API connection.
        </div>
      )}

      {loading ? (
        <div className="card text-gray-500">Loading stats...</div>
      ) : stats ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
          {statCards(stats).map(({ label, value }) => (
            <div key={label} className="card">
              <div className="text-3xl font-bold text-brand-700">{value}</div>
              <div className="mt-1 text-sm text-gray-500">{label}</div>
            </div>
          ))}
        </div>
      ) : null}

      <section className="space-y-4">
        <h2 className="text-lg font-semibold text-gray-950">Quick actions</h2>
        <div className="grid gap-4 sm:grid-cols-2">
          <Link href="/admin/candidates" className="card block transition hover:border-brand-300 hover:shadow-md">
            <div className="font-semibold text-gray-950">Manage candidates</div>
            <div className="mt-1 text-sm text-gray-500">Edit bios, activate or deactivate candidates.</div>
          </Link>
          <Link href="/admin/candidates/new" className="card block transition hover:border-brand-300 hover:shadow-md">
            <div className="font-semibold text-gray-950">Add candidate</div>
            <div className="mt-1 text-sm text-gray-500">Register a new candidate before voting opens.</div>
          </Link>
        </div>
      </section>

      <section className="space-y-4">
        <h2 className="text-lg font-semibold text-gray-950">System status</h2>
        <div className="card divide-y divide-gray-100">
          {buildStatusRows(health ?? { status: "unknown" }).map(({ label, value }) => (
            <div key={label} className="flex items-center justify-between py-3 first:pt-0 last:pb-0">
              <div className="text-sm text-gray-500">{label}</div>
              <div className={`text-sm font-medium ${value === "Error" || value === "Degraded" ? "text-red-600" : "text-gray-950"}`}>
                {loading ? "—" : value}
              </div>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
