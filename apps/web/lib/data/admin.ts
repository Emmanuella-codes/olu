import { AdminStats } from "@/types/types";

export const statCards = (s: AdminStats) => [
    { label: "Total votes",  value: s.total_votes.toLocaleString("en-NG") },
    { label: "Web votes",    value: s.web_votes.toLocaleString("en-NG") },
    { label: "SMS votes",    value: s.sms_votes.toLocaleString("en-NG") },
    { label: "Pending SMS",  value: s.pending_sms.toLocaleString("en-NG") },
    { label: "Candidates",   value: s.total_candidates.toLocaleString("en-NG") }
]; 

export function buildStatusRows(health: { status: string; error?: string }) {
    const ok = health.status === "ok";
    const dbError = health.error?.toLowerCase().includes("database");
    const redisError = health.error?.toLowerCase().includes("redis");

    return [
        { label: "API server",   value: ok ? "Healthy" : "Degraded" },
        { label: "Database",     value: dbError ? "Error" : "Connected" },
        { label: "Redis cache",  value: redisError ? "Error" : "Connected" },
        { label: "SMS provider", value: "Mock" },
    ];
}
