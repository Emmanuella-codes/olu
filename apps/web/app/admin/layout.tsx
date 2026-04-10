export default function AdminLayout({ children }: { children: React.ReactNode }) {
    return (
        <div style={{ display: "flex", minHeight: "100vh", fontFamily: "var(--font-sans)" }}>
            <aside style={{ width: 210, background: "#111827", display: "flex", flexDirection: "column", flexShrink: 0 }}>
                <div style={{ padding: "16px 20px", borderBottom: "1px solid rgba(255,255,255,0.08)" }}>
                    <div style={{ fontSize: 16, fontWeight: 500, color: "#fff" }}>Olu</div>
                    <div style={{ fontSize: 11, color: "#6b7280", marginTop: 2 }}>Admin panel</div>
                </div>
            </aside>
        </div>
    );
}
