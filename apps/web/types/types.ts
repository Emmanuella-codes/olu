export type Step = "phone" | "otp" | "confirm" | "done";

export interface Candidate {
    id: string;
    code: string;
    name: string;
    party: string;
    bio: string;
    achievements: string;
    photo_url?: string;
    is_active: boolean;
    created_at: string;
}

export interface TallyRow {
    candidate_id: string;
    code: string;
    name: string;
    party: string;
    vote_count: number;
}

export interface Results {
    tally: TallyRow[];
    total_votes: number;
    is_tie: boolean;
    leaders: TallyRow[];
    cached_at: string;
}

export interface ApiError {
    error: string;
    code: string;
}
