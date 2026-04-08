import { Candidate } from "@/types/types";
import { getCandidates } from "../lib/api/api";
import CandidateCard from "@/components/CandidateCard";

// Revalidate every 5 minutes — candidates rarely change
export const revalidate = 300;

export default async function Home() {
  let candidates: Candidate[] = [];
  let error = false;

  try {
    candidates = await getCandidates();
  } catch {
    error = true;
  }

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">
          Meet the candidates
        </h1>
        <p className="mt-2 text-gray-500">
          Read about each candidate before casting your vote.
        </p>
      </div>
 
      {error && (
        <div className="card bg-red-50 border-red-200 text-red-700 mb-6">
          Could not load candidates. Please refresh the page.
        </div>
      )}
 
      {candidates.length === 0 && !error && (
        <div className="card text-center text-gray-400 py-16">
          No candidates have been listed yet.
        </div>
      )}
 
      <div className="grid gap-6 sm:grid-cols-2">
        {candidates.map((c) => (
          <CandidateCard key={c.id} candidate={c} />
        ))}
      </div>
    </div>
  );
}
