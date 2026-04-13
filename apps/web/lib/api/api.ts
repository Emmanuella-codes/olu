import { ApiError, Candidate, Results } from "@/types/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export class ApiRequestError extends Error {
    constructor(
        public readonly status: number,
        message: string
    ) {
        super(message);
        this.name = "ApiRequestError";
    }
}

export async function handleResponse<T>(path: string, options?: RequestInit): Promise<T> {
    let res: Response;
    try {
        res = await fetch(`${API_URL}/api/v1${path}`, {
            headers: { "Content-Type": "application/json", ...options?.headers },
            ...options,
        });
    } catch {
        throw new ApiRequestError(0, "Network error. Please check your connection.");
    }

    const json = await res.json();

    if (!res.ok) {
        const err = json as ApiError;
        throw new ApiRequestError(res.status, err.error ?? `HTTP ${res.status}`);
    }

    return json;
}

// candidates

export async function getCandidates(): Promise<Candidate[]> {
    const data = await handleResponse<{ data: Candidate[]; count: number }>("/candidates");
    return data.data;
}

export async function getCandidate(id: string): Promise<Candidate> {
    const data = await handleResponse<{ data: Candidate }>(`/candidates/${id}`);
    return data.data;
}

// otp

export async function sendOTP(phone: string): Promise<void> {
    await handleResponse<void>("/auth/send-otp", {
        method: "POST",
        body: JSON.stringify({ phone }),
    });
}

export async function verifyOTP(phone: string, code: string): Promise<{ token: string }> {
    return handleResponse<{ token: string }>("/auth/verify-otp", {
        method: "POST",
        body: JSON.stringify({ phone, code }),
    });
}

// vote

export async function castVote(
    candidateCode: string,
    token: string
): Promise<{ confirmation_id: string; candidate_name: string }> {
    const res = await handleResponse<{ data: { confirmation_id: string; candidate_name: string } }>("/vote", {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: JSON.stringify({ candidate_code: candidateCode }),
    });
    return res.data;
}

// health

export async function getHealth(): Promise<{ status: string; error?: string }> {
    let res: Response;
    try {
        res = await fetch(`${API_URL}/health`);
    } catch {
        return { status: "degraded", error: "unreachable" };
    }
    const json = await res.json();
    if (!res.ok) return { status: "degraded", error: json.error };
    return json;
}

// results

export async function getResults(): Promise<Results> {
    const data = await handleResponse<{ data: Results }>("/results");
    return data.data;
}
