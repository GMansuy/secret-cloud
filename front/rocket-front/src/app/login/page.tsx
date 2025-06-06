"use client";

import getUserManager from "@/app/auth-config";
import { Button } from "@mui/material";

export default function Login() {
        const handleLogin = async () => {
            try {
                const userManager = getUserManager();
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
            <Button variant="contained"
                    onClick={handleLogin}
            >Login with Microsoft</Button>
        </main>
    );
}