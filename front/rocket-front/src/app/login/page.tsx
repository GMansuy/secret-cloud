"use client";

import userManager from "@/app/auth-config";

export default function Login() {
    const handleLogin = async () => {
        try {
            await userManager.signinRedirect();
        } catch (error) {
            console.error("Login failed:", error);
        }
    };

    return (
        <main
            style={{
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                justifyContent: "center",
                height: "100vh",
                fontFamily: "Arial, sans-serif",
            }}
        >
            <h2 style={{ marginBottom: "1rem" }}>Login with Microsoft</h2>
            <button
                onClick={handleLogin}
                style={{
                    padding: "0.5rem 1rem",
                    fontSize: "1rem",
                    cursor: "pointer",
                    backgroundColor: "#0078D4",
                    color: "#fff",
                    border: "none",
                    borderRadius: "5px",
                    transition: "background-color 0.3s",
                }}
                onMouseOver={(e) => (e.currentTarget.style.backgroundColor = "#005A9E")}
                onMouseOut={(e) => (e.currentTarget.style.backgroundColor = "#0078D4")}
            >
                Login
            </button>
        </main>
    );
}