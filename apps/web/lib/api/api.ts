import { ApiError, Candidate, Results } from "@/types/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

async function handleResponse<T>(path: string, options?: RequestInit): Promise<T> {
    const res = await fetch(`${API_URL}/api/v1${path}`, {
        headers: { "Content-Type": "application/json", ...options?.headers },
        ...options,
    });

    const json = await res.json();

    if (!res.ok) {
        const err = json as ApiError;
        throw new Error(err.error ?? `HTTP ${res.status}`);
    }

    return json
}

// candidates

export async function getCandidates(): Promise<Candidate[]> {
    const data = await handleResponse<{ data: Candidate[]; count: number }>("/candidates");
    return data.data;
}

export async function getCandidate(id: string): Promise<Candidate> {
    const data = await handleResponse<Candidate>(`/candidates/${id}`);
    return data;
}

// otp

export async function sendOTP(phone: string): Promise<void> {
    await handleResponse<void>("/auth/send-otp", {
        method: "POST",
        body: JSON.stringify({ phone }),
    });
}

export async function verifyOTP(phone: string, code: string): Promise<{ token: string; message: string }> {
    return handleResponse<{ token: string; message: string }>("/auth/verify-otp", {
        method: "POST",
        body: JSON.stringify({ phone, code }),
    });
}

// vote

export async function castVote(
    candidateCode: string,
    token: string
): Promise<{ message: string; confirmation_id: string; candidate: string }> {
    return handleResponse("/vote", {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: JSON.stringify({ candidate_code: candidateCode }),
    })
}

// results

export async function getResults(): Promise<Results> {
    const data = await handleResponse<{ data: Results }>("/results");
    return data.data;
}
