"use client";

import React, { useState } from "react";
import axios from "axios";
import { useRouter } from "next/navigation";
import {Button, CircularProgress, Stack, TextField} from "@mui/material";
import { ArrowBack } from "@mui/icons-material";
import Link from "next/link";

// Helper function to get the full URL with appropriate protocol
const getFullUrl = (baseUrl: string | undefined, path: string): string => {
    if (!baseUrl) return path;
    const protocol = baseUrl.startsWith('http') ? '' : 'https://';
    return `${protocol}${baseUrl}${path}`;
};

export default function Cluster() {
    const [data, setData] = useState<string>("");
    const [error, setError] = useState<string>("");
    const [loading, setLoading] = useState<boolean>(false);
    const [formData, setFormData] = useState({
        name: "",
        controlplaneMachineCount: 1,
        workerMachineCount: 1,
    });
    const router = useRouter();
    const backendUrl = process.env.NEXT_PUBLIC_BACKEND_URL


    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: name === "controlplaneMachineCount" || name === "workerMachineCount" ? parseInt(value) : value,
        });
    };

    const postCluster = async (e: React.FormEvent) => {
        setLoading(true);
        e.preventDefault();
        setError("");
        setData("");
        try {
            const response = await axios.post(getFullUrl(backendUrl, '/cluster'), formData);
            setData("Cluster created successfully!");
            router.push("/clusterlist");
        } catch (err) {
            setError("Failed to create cluster.");
        }
        setLoading(false);
    };

    return (
        <main
            style={{
                maxWidth: 600,
                margin: "2rem auto",
                fontFamily: "Arial, sans-serif",
                textAlign: "center",
            }}
        >
            <Stack spacing={6}>
            <Link href={"/clusterlist"}>
            <Button variant="outlined" startIcon={<ArrowBack />}>
            to Clusters
            </Button>
            </Link>

            {error && <p style={{ color: "red" }}>{error}</p>}
            {data && <p style={{ color: "green" }}>{data}</p>}
            <form onSubmit={postCluster} style={{ marginBottom: "1rem" }}>
                <h1 style={{ color: "#333", marginBottom: "1rem" }}>Post Cluster Form</h1>
                <Stack spacing={3}>

                    <TextField id="outlined-basic"
                               name="name"
                               value={formData.name}
                               onChange={handleInputChange}
                               label="Name"
                               variant="outlined" />
                    <TextField
                        type="number"
                        name="controlplaneMachineCount"
                        value={formData.controlplaneMachineCount}
                        onChange={handleInputChange}
                        label="Control plane Machine Number"
                        variant="outlined" />

                    <TextField
                        type="number"
                        name="workerMachineCount"
                        value={formData.workerMachineCount}
                        onChange={handleInputChange}
                        label="Worker plane Machine Number"
                        variant="outlined" />
                {loading ?
                        <CircularProgress /> :
                <Button variant="contained" type={"submit"}>
                    Submit
                </Button>}
                </Stack>
            </form>
            </Stack>
        </main>
    );
}
