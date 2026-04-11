"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

import CandidateForm from "@/components/admin/CandidateForm";
import { createCandidate } from "@/lib/admin-api/api";
import { CreateCandidatePayload } from "@/types/types";

function getErrorMessage(error: unknown, fallback: string) {
  return error instanceof Error && error.message ? error.message : fallback;
}

export default function NewCandidatePage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (candidate: CreateCandidatePayload) => {
    setError("");
    setLoading(true);

    try {
      await createCandidate(candidate);
      router.push("/admin/candidates");
    } catch (error) {
      setError(getErrorMessage(error, "Failed to create candidate."));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <Link href="/admin/candidates" className="text-sm font-medium text-brand-700 hover:text-brand-800">
          ← Back to candidates
        </Link>
        <h1 className="mt-3 text-3xl font-bold text-gray-950">Add candidate</h1>
        <p className="mt-1 text-sm text-gray-500">Candidate codes must use the A1 format.</p>
      </div>

      <div className="card">
        <CandidateForm onSubmit={handleSubmit} submitLabel="Create candidate" loading={loading} error={error} />
      </div>
    </div>
  );
}
