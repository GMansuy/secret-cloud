"use client";

import { useEffect } from "react";
import userManager from "@/app/auth-config";

export default function Callback() {
    useEffect(() => {
        userManager.signinRedirectCallback()
            .then((user) => {
                localStorage.setItem("access_token", user.access_token);
                window.location.href = "/clusterlist"; // Redirect to the cluster list page after successful login
            })
            .catch((err) => {
                console.error("Login failed:", err);
                document.body.innerHTML = "Login failed. "+ err.message;
            });
    }, []);

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
            <h2>Processing login...</h2>
        </main>
    );
}