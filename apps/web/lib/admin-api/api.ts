import {
  AdminCandidate,
  AdminStats,
  ApiError,
  CreateCandidatePayload,
  UpdateCandidatePayload,
} from "@/types/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL;
const ADMIN_TOKEN_KEY = process.env.NEXT_PUBLIC_ADMIN_TOKEN_KEY ?? "olu_admin_token";

export class AdminApiError extends Error {
  constructor(
    public readonly status: number,
    message: string
  ) {
    super(message);
    this.name = "AdminApiError";
  }
}

function getAdminStorage(): Storage | null {
  if (typeof window === "undefined") {
    return null;
  }

  return window.sessionStorage;
}

export function getAdminToken(): string | null {
  return getAdminStorage()?.getItem(ADMIN_TOKEN_KEY) ?? null;
}

export function setAdminToken(token: string): void {
  getAdminStorage()?.setItem(ADMIN_TOKEN_KEY, token);
}

export function clearAdminToken(): void {
  getAdminStorage()?.removeItem(ADMIN_TOKEN_KEY);
}

async function adminRequest<T>(path: string, options?: RequestInit): Promise<T> {
  let response: Response;

  try {
    response = await fetch(`${API_URL}/api/v1/admin${path}`, {
      headers: {
        "Content-Type": "application/json",
        ...options?.headers,
      },
      ...options,
    });
  } catch {
    throw new AdminApiError(0, "Network error. Please check your connection.");
  }

  const json = (await response.json()) as T | ApiError;

  if (!response.ok) {
    const apiError = json as ApiError;
    throw new AdminApiError(response.status, apiError.error ?? `HTTP ${response.status}`);
  }

  return json as T;
}

function adminAuthHeaders(token: string) {
  return { Authorization: `Bearer ${token}` };
}

function requireAdminToken(token?: string): string {
  const resolvedToken = token ?? getAdminToken();
  if (!resolvedToken) {
    throw new AdminApiError(401, "Admin token is missing. Please sign in again.");
  }

  return resolvedToken;
}

export async function adminLogin(email: string, password: string): Promise<{ token: string; expires_in: number }> {
  return adminRequest<{ token: string; expires_in: number }>("/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export async function createAdmin(
  payload: { email: string; password: string },
  token?: string
): Promise<{ id: string; email: string }> {
  const adminToken = requireAdminToken(token);
  const response = await adminRequest<{ data: { id: string; email: string } }>("/create", {
    method: "POST",
    headers: adminAuthHeaders(adminToken),
    body: JSON.stringify(payload),
  });

  return response.data;
}

export async function getAdminCandidate(id: string, token?: string): Promise<AdminCandidate> {
  const adminToken = requireAdminToken(token);
  const response = await adminRequest<{ data: AdminCandidate }>(`/candidates/${id}`, {
    headers: adminAuthHeaders(adminToken),
  });
  return response.data;
}

export async function getAdminCandidates(token?: string): Promise<AdminCandidate[]> {
  const adminToken = requireAdminToken(token);
  const response = await adminRequest<{ data: AdminCandidate[]; count: number }>("/candidates", {
    headers: adminAuthHeaders(adminToken),
  });

  return response.data;
}

export async function createCandidate(payload: CreateCandidatePayload, token?: string): Promise<AdminCandidate> {
  const adminToken = requireAdminToken(token);
  const response = await adminRequest<{ data: AdminCandidate }>("/candidates", {
    method: "POST",
    headers: adminAuthHeaders(adminToken),
    body: JSON.stringify(payload),
  });

  return response.data;
}

export async function updateCandidate(
  candidateId: string,
  payload: UpdateCandidatePayload,
  token?: string
): Promise<AdminCandidate> {
  const adminToken = requireAdminToken(token);
  const response = await adminRequest<{ data: AdminCandidate }>(`/candidates/${candidateId}`, {
    method: "PUT",
    headers: adminAuthHeaders(adminToken),
    body: JSON.stringify(payload),
  });

  return response.data;
}

export async function deactivateCandidate(candidateId: string, token?: string): Promise<string> {
  const adminToken = requireAdminToken(token);
  const response = await adminRequest<{ data: string }>(`/candidates/${candidateId}`, {
    method: "DELETE",
    headers: adminAuthHeaders(adminToken),
  });

  return response.data;
}

export async function getAdminStats(token?: string): Promise<AdminStats> {
  const adminToken = requireAdminToken(token);
  const response = await adminRequest<{ data: AdminStats }>("/stats", {
    headers: adminAuthHeaders(adminToken),
  });

  return response.data;
}
