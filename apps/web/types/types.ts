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

export interface AdminCandidate {
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

export interface AdminStats {
    total_votes: number;
    web_votes: number;
    sms_votes: number;
    pending_sms: number;
    total_candidates: number;
}

export interface CreateCandidatePayload {
    code: string;
    name: string;
    party: string;
    bio: string;
    achievements: string;
    photo_url?: string;
}

export interface UpdateCandidatePayload {
    code?: string;
    name?: string;
    party?: string;
    bio?: string;
    achievements?: string;
    photo_url?: string;
    is_active?: boolean;
}
