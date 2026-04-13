import Link from "next/link";

import { splitAchievements } from "@/lib/format/achievements";
import { Candidate } from "@/types/types";

interface Props {
  candidate: Candidate;
}

export default function CandidateProfileView({ candidate }: Props) {
  const achievements = splitAchievements(candidate.achievements);

  return (
    <div className="mx-auto max-w-2xl">
      <Link
        href="/"
        className="mb-6 inline-flex min-h-0 items-center gap-2 text-sm text-brand-600 hover:text-brand-800"
      >
        ← Back to candidates
      </Link>

      <div className="card">
        <div className="mb-6 flex items-start gap-4">
          <div className="flex h-20 w-20 shrink-0 items-center justify-center rounded-full bg-brand-100 text-2xl font-bold text-brand-700">
            {candidate?.name?.charAt(0)}
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{candidate.name}</h1>
            <p className="mt-0.5 font-medium text-brand-600">{candidate.party.toLocaleUpperCase()}</p>
            <span className="mt-2 inline-block rounded-full border border-brand-200 bg-brand-50 px-3 py-1 text-xs font-semibold text-brand-700">
              Code: {candidate.code}
            </span>
          </div>
        </div>

        <section className="mb-6">
          <h2 className="mb-2 text-sm font-semibold uppercase tracking-wider text-gray-400">
            About
          </h2>
          <p className="leading-relaxed text-gray-700">{candidate.bio}</p>
        </section>

        {achievements?.length > 0 && (
          <section className="mb-8">
            <h2 className="mb-3 text-sm font-semibold uppercase tracking-wider text-gray-400">
              Key achievements
            </h2>
            <ul className="space-y-2">
              {achievements.map((item) => (
                <li key={item} className="flex items-start gap-3 text-gray-700">
                  <span className="mt-1 h-2 w-2 shrink-0 rounded-full bg-brand-400" />
                  {item}
                </li>
              ))}
            </ul>
          </section>
        )}

        <div className="flex flex-col gap-3 border-t border-gray-100 pt-6 sm:flex-row">
          <Link href={`/vote?code=${candidate.code}`} className="btn-primary justify-center">
            Vote for {candidate?.name?.split(" ")[0]}
          </Link>
          <Link href="/" className="btn-secondary justify-center">
            See all candidates
          </Link>
        </div>
      </div>

      <p className="mt-4 text-center text-sm text-gray-400">
        No internet? Send <strong>VOTE {candidate.code}</strong> to <strong>****</strong>
      </p>
    </div>
  );
}
