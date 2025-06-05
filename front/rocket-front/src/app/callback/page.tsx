"use client";

import { useEffect } from "react";
import getUserManager from "@/app/auth-config";
import { useRouter } from 'next/navigation';

export default function Callback() {
    const router = useRouter();
    useEffect(() => {
        const userManager = getUserManager();
        userManager.signinRedirectCallback()
            .then((user) => {
                localStorage.setItem("access_token", user.access_token);
                router.push("/clusterlist");
            })
            .catch((err) => {
                console.error("Login failed:", err);
                document.body.innerHTML = "Login failed. "+ err.message;
            });
    }, [router]);

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