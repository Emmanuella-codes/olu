import { Candidate } from "@/types/types";
import Link from "next/link";

interface Props {
    candidate: Candidate
}

export default function CandidateCard({ candidate }: Props) {
    const firstAchievement = candidate.achievements
        .split("\n")
        .map((a) => a.trim())
        .find(Boolean);

    return (
        <div className="card flex flex-col hover:border-brand-300 hover:shadow-md transition-all duration-150">
            {/* header */}
            <div className="flex items-center gap-3 mb-4">
                <div className="w-12 h-12 rounded-full bg-brand-100 flex items-center justify-center shrink-0 text-lg font-bold text-brand-700">
                    {candidate.name.charAt(0)}
                </div>
                <div className="min-w-0">
                    <h2 className="font-semibold text-gray-900 truncate">{candidate.name}</h2>
                    <p className="text-sm text-brand-600 truncate">{candidate.party.toLocaleUpperCase()}</p>
                </div>
                <span className="ml-auto shrink-0 bg-brand-50 text-brand-700 text-xs font-bold px-2 py-1 rounded border border-brand-200 font-mono">
                    {candidate.code}
                </span>
            </div>
            {/* bio */}
            <p className="text-sm text-gray-500 leading-relaxed line-clamp-3 flex-1">
                {candidate.bio}
            </p>
        
            {/* achievement */}
            {firstAchievement && (
                <div className="mt-3 flex items-start gap-2 text-xs text-gray-400">
                    <span className="w-1.5 h-1.5 rounded-full bg-brand-400 mt-1 shrink-0" />
                    <span className="line-clamp-2">{firstAchievement}</span>
                </div>
            )}
            {/* actions */}
            <div className="mt-4 pt-4 border-t border-gray-100 flex gap-2">
                <Link
                    href={`/candidates/${candidate.id}`}
                    className="btn-secondary text-sm px-4 py-2 min-h-0 h-auto flex-1 justify-center"
                >
                    Full profile
                </Link>
                <Link
                    href={`/vote?code=${candidate.code}`}
                    className="btn-primary text-sm px-4 py-2 min-h-0 h-auto flex-1 justify-center"
                >
                    Vote
                </Link>
            </div>
        </div>
    );
}
