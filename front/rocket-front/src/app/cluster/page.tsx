"use client";

import React, { useState } from "react";
import axios from "axios";

export default function Cluster() {
    const [data, setData] = useState<string>("");
    const [error, setError] = useState<string>("");
    const [formData, setFormData] = useState({
        name: "",
        controlplaneMachineCount: 1,
        workerMachineCount: 2,
    });

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: name === "controlplaneMachineCount" || name === "workerMachineCount" ? parseInt(value) : value,
        });
    };

    const postCluster = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        setData("");
        try {
            const response = await axios.post("http://localhost:8080/cluster", formData);
            setData("Cluster created successfully!");
            console.log(response.data);
        } catch (err) {
            setError("Failed to create cluster.");
            console.error(err);
        }
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
            <h1 style={{ color: "#333", marginBottom: "1rem" }}>Post Cluster Form</h1>
            {error && <p style={{ color: "red" }}>{error}</p>}
            {data && <p style={{ color: "green" }}>{data}</p>}
            <form onSubmit={postCluster} style={{ marginBottom: "1rem" }}>
                <div style={{ marginBottom: "1rem" }}>
                    <label style={{ display: "block", marginBottom: "0.5rem" }}>Name:</label>
                    <input
                        type="text"
                        name="name"
                        value={formData.name}
                        onChange={handleInputChange}
                        style={{ padding: "0.5rem", width: "100%" }}
                    />
                </div>
                <div style={{ marginBottom: "1rem" }}>
                    <label style={{ display: "block", marginBottom: "0.5rem" }}>Control Plane Machine Count:</label>
                    <input
                        type="number"
                        name="controlplaneMachineCount"
                        value={formData.controlplaneMachineCount}
                        onChange={handleInputChange}
                        style={{ padding: "0.5rem", width: "100%" }}
                    />
                </div>
                <div style={{ marginBottom: "1rem" }}>
                    <label style={{ display: "block", marginBottom: "0.5rem" }}>Worker Machine Count:</label>
                    <input
                        type="number"
                        name="workerMachineCount"
                        value={formData.workerMachineCount}
                        onChange={handleInputChange}
                        style={{ padding: "0.5rem", width: "100%" }}
                    />
                </div>
                <button
                    type="submit"
                    style={{
                        padding: "0.5rem 1rem",
                        fontSize: "1rem",
                        cursor: "pointer",
                        backgroundColor: "#28a745",
                        color: "#fff",
                        border: "none",
                        borderRadius: "5px",
                        transition: "background-color 0.3s",
                    }}
                    onMouseOver={(e) => (e.currentTarget.style.backgroundColor = "#218838")}
                    onMouseOut={(e) => (e.currentTarget.style.backgroundColor = "#28a745")}
                >
                    Submit
                </button>
            </form>
        </main>
    );
}