import { Candidate } from "@/types/types";

interface Props {
    candidate: Candidate;
    onConfirm: () => void;
    onCancel: () => void;
    loading: boolean;
}

export default function ConfirmationModal({ candidate, onConfirm, onCancel, loading }: Props) {
    return (
        <div className="card border-brand-200 space-y-6">
            <div className="text-center">
                <p className="text-sm text-gray-400 uppercase tracking-wider font-semibold">
                    Confirm your vote
                </p>
                <p className="mt-2 text-gray-500 text-sm">
                    This action is <strong>irreversible</strong>. You may only vote once.
                </p>
            </div>
 
            {/* candidate summary */}
            <div className="bg-brand-50 border border-brand-200 rounded-xl p-4 flex items-center gap-4">
                <div className="w-14 h-14 rounded-full bg-brand-200 flex items-center justify-center text-xl font-bold text-brand-700 flex-shrink-0">
                    {candidate.name.charAt(0)}
                </div>
                    <div>
                        <p className="font-bold text-gray-900 text-lg">{candidate.name}</p>
                        <p className="text-brand-600 text-sm">{candidate.party}</p>
                        <p className="text-xs font-mono text-gray-400 mt-0.5">
                            Code: {candidate.code}
                        </p>
                </div>
            </div>
 
            <p className="text-xs text-center text-gray-400">
                By clicking &quot;Confirm vote&quot; you agree that your vote is final and cannot be withdrawn.
            </p>
 
            <div className="flex flex-col gap-3">
                <button
                    className="btn-primary justify-center"
                    onClick={onConfirm}
                    disabled={loading}
                >
                    {loading ? "Submitting..." : "Confirm vote"}
                </button>
                <button
                    className="btn-secondary justify-center"
                    onClick={onCancel}
                    disabled={loading}
                >
                    Go back
                </button>
            </div>
    </div>
    )
}
