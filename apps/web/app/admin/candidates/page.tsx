"use client";

import Link from "next/link";
import { useCallback, useEffect, useState } from "react";

import { deactivateCandidate, getAdminCandidates, updateCandidate } from "@/lib/admin-api/api";
import { AdminCandidate } from "@/types/types";

function getErrorMessage(error: unknown, fallback: string) {
  return error instanceof Error && error.message ? error.message : fallback;
}

export default function AdminCandidatePage() {
  const [candidates, setCandidates] = useState<AdminCandidate[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [togglingId, setTogglingId] = useState<string | null>(null);

  const load = useCallback(async () => {
    setError("");
    try {
      setCandidates(await getAdminCandidates());
    } catch (error) {
      setError(getErrorMessage(error, "Failed to load candidates."));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const toggleActive = async (candidate: AdminCandidate) => {
    setTogglingId(candidate.id);
    setError("");

    try {
      if (candidate.is_active) {
        await deactivateCandidate(candidate.id);
      } else {
        await updateCandidate(candidate.id, { is_active: true });
      }
      await load();
    } catch (error) {
      setError(getErrorMessage(error, "Failed to toggle candidate status."));
    } finally {
      setTogglingId(null);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.2em] text-brand-600">Admin</p>
          <h1 className="mt-2 text-3xl font-bold text-gray-950">Candidates</h1>
        </div>
        <Link href="/admin/candidates/new" className="btn-primary justify-center">
          + Add candidate
        </Link>
      </div>

      {error && <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>}
      {loading && <div className="card text-sm text-gray-500">Loading candidates...</div>}

      {!loading && candidates.length === 0 && !error && (
        <div className="card space-y-4 text-center">
          <p className="text-gray-500">No candidates yet.</p>
          <Link href="/admin/candidates/new" className="btn-primary justify-center">
            Add the first one
          </Link>
        </div>
      )}

      <div className="space-y-3">
        {candidates.map((candidate) => (
          <div
            key={candidate.id}
            className={`card flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between ${
              candidate.is_active ? "" : "opacity-60"
            }`}
          >
            <div className="flex min-w-0 items-center gap-4">
              <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-brand-100 text-lg font-bold text-brand-700">
                {candidate.name.charAt(0)}
              </div>
              <div className="min-w-0">
                <div className="flex flex-wrap items-center gap-2">
                  <span className="font-semibold text-gray-950">{candidate.name}</span>
                  <span className="rounded-full bg-gray-100 px-2 py-0.5 font-mono text-xs text-gray-500">
                    {candidate.code}
                  </span>
                  {!candidate.is_active && (
                    <span className="rounded-full bg-red-50 px-2 py-0.5 text-xs font-medium text-red-700">
                      Inactive
                    </span>
                  )}
                </div>
                <div className="mt-1 text-sm text-gray-500">{candidate.party.toLocaleUpperCase()}</div>
              </div>
            </div>

            <div className="flex gap-2">
              <Link href={`/admin/candidates/${candidate.id}`} className="btn-secondary h-auto min-h-0 px-4 py-2">
                Edit
              </Link>
              <button
                className="btn-secondary h-auto min-h-0 px-4 py-2"
                onClick={() => toggleActive(candidate)}
                disabled={togglingId === candidate.id}
              >
                {togglingId === candidate.id ? "Saving..." : candidate.is_active ? "Deactivate" : "Activate"}
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
