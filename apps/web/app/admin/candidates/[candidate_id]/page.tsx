"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import CandidateForm from "@/components/admin/CandidateForm";
import { getAdminCandidate, updateCandidate } from "@/lib/admin-api/api";
import { AdminCandidate, CreateCandidatePayload } from "@/types/types";

function getErrorMessage(error: unknown, fallback: string) {
  return error instanceof Error && error.message ? error.message : fallback;
}

export default function EditCandidatePage() {
  const params = useParams<{ candidate_id: string }>();
  const router = useRouter();
  const [candidate, setCandidate] = useState<AdminCandidate | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    getAdminCandidate(params.candidate_id)
      .then(setCandidate)
      .catch((error) => setError(getErrorMessage(error, "Failed to load candidate.")))
      .finally(() => setLoading(false));
  }, [params.candidate_id]);

  const handleSubmit = async (payload: CreateCandidatePayload) => {
    if (!candidate) {
      return;
    }

    setError("");
    setSaving(true);

    try {
      await updateCandidate(candidate.id, payload);
      router.push("/admin/candidates");
    } catch (error) {
      setError(getErrorMessage(error, "Failed to update candidate."));
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <Link href="/admin/candidates" className="text-sm font-medium text-brand-700 hover:text-brand-800">
          ← Back to candidates
        </Link>
        <h1 className="mt-3 text-3xl font-bold text-gray-950">Edit candidate</h1>
        <p className="mt-1 text-sm text-gray-500">Update profile details and voting code.</p>
      </div>

      {loading && <div className="card text-sm text-gray-500">Loading candidate...</div>}

      {!loading && candidate && (
        <div className="card">
          <CandidateForm
            initial={candidate}
            onSubmit={handleSubmit}
            submitLabel="Save changes"
            loading={saving}
            error={error}
          />
        </div>
      )}

      {!loading && !candidate && error && (
        <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>
      )}
    </div>
  );
}
