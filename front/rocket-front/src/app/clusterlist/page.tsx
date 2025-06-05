"use client";

import React, { useState, useEffect } from "react";
import axios from "axios";
import Link from "next/link";

interface ClusterData {
    name: string;
    status: string;
}

interface clusterStatus {
    type: string;
    status: string;
    reason?: string;
}

export default function ClustersList() {
    const [clusters, setClusters] = useState<ClusterData[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>("");
    const [data, setData] = useState<string>("");

    useEffect(() => {
        fetchClusters();
    }, []);

    const fetchClusters = async () => {
        setLoading(true);
        try {
            const response = await axios.get("http://localhost:8080/list");

            // Parse the response data - this depends on the actual format returned by your API
            const clusterData = parseClustersData(response.data.clusters);
            setClusters(clusterData);
            setError("");

            console.log("api clusters:", response.data.clusters);
            console.log( "Fetched clusters:", clusterData);
        } catch (err) {
            setError("Failed to fetch clusters");
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    // This function will need to be adapted based on the actual format of your API response
    const parseClustersData = (rawData: string): ClusterData[] => {
        if (!rawData) return [];

        try {
            // Find all matches of cluster name followed by JSON array of conditions
            const clusterPattern = /([a-zA-Z0-9-]+)(\[.*?\])/g;
            const clusters: ClusterData[] = [];
            let match;

            while ((match = clusterPattern.exec(rawData)) !== null) {
                const clusterName = match[1];
                const conditionsStr = match[2];

                try {
                    // Parse the conditions array
                    const conditions = JSON.parse(conditionsStr);

                    // Find specific conditions
                    const controlPlaneReady = conditions.find(
                        (c: clusterStatus) => c.type === "ControlPlaneReady"
                    );

                    const infrastructureReady = conditions.find(
                        (c: clusterStatus) => c.type === "InfrastructureReady"
                    );

                    // Generate status text based on conditions
                    let status = "Unknown";

                    if (controlPlaneReady && infrastructureReady) {
                        if (controlPlaneReady.status === "True" && infrastructureReady.status === "True") {
                            status = "Ready";
                        } else if (controlPlaneReady.status === "False" || infrastructureReady.status === "False") {
                            status = "Not Ready";
                            // Add reason if available
                            const reason = controlPlaneReady.reason || infrastructureReady.reason;
                            if (reason) status += ` (${reason})`;
                        }
                    } else if (controlPlaneReady) {
                        status = controlPlaneReady.status === "True" ? "Control Plane Ready" : "Control Plane Not Ready";
                    } else if (infrastructureReady) {
                        status = infrastructureReady.status === "True" ? "Infrastructure Ready" : "Infrastructure Not Ready";
                    }

                    clusters.push({
                        name: clusterName,
                        status: status
                    });
                } catch (e) {
                    console.error(`Error parsing conditions for cluster ${clusterName}:`, e);
                    clusters.push({
                        name: clusterName,
                        status: "Error parsing status"
                    });
                }
            }

            return clusters;
        } catch (e) {
            console.error("Error parsing cluster data:", e);
            return [];
        }
    };

    const deleteCluster = async (clusterName: string) => {
        setData("");
        try {
            await axios.delete("http://localhost:8080/cluster", {
                data: { name: clusterName }
            });
            await fetchClusters(); // Refresh the list
            setData(`Cluster ${clusterName} deleted successfully!`);
        } catch (err) {
            setError("Failed to delete cluster");
            console.error(err);
        }
    };
    const downloadKubeconfig = async (clusterName: string) => {
        try {
            const response = await axios.get(`http://localhost:8080/cluster/${clusterName}/kubeconfig`);

            // Create a blob with the kubeconfig content
            const blob = new Blob([response.data.kubeconfig], { type: 'text/yaml' });

            // Create a URL for the blob
            const url = window.URL.createObjectURL(blob);

            // Create a temporary anchor element to trigger the download
            const link = document.createElement('a');
            link.href = url;
            link.download = `${clusterName}-kubeconfig.yaml`;

            // Append to the document, click it, and clean up
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);

            // Release the blob URL
            window.URL.revokeObjectURL(url);

            setData(`Kubeconfig for ${clusterName} downloaded successfully!`);
        } catch (err) {
            setError(`Failed to download kubeconfig for ${clusterName}`);
            console.error(err);
        }
    };

    return (
        <main style={{
            maxWidth: 800,
            margin: "2rem auto",
            fontFamily: "Arial, sans-serif",
        }}>
            <h1 style={{ color: "#333", marginBottom: "1rem", textAlign: "center" }}>Clusters</h1>
            {error && <p style={{ color: "red" }}>{error}</p>}
            {data && <p style={{ color: "green" }}>{data}</p>}
            <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "1rem" }}>
                <button
                    onClick={fetchClusters}
                    style={{
                        padding: "0.5rem 1rem",
                        backgroundColor: "#0078d4",
                        color: "white",
                        border: "none",
                        borderRadius: "4px",
                        cursor: "pointer",
                    }}
                >
                    Refresh
                </button>


                <Link href="/cluster">
                    <button style={{
                        padding: "0.5rem 1rem",
                        backgroundColor: "#28a745",
                        color: "white",
                        border: "none",
                        borderRadius: "4px",
                        cursor: "pointer",
                    }}>
                        Create New Cluster
                    </button>
                </Link>
            </div>

            {error && <p style={{ color: "red", textAlign: "center" }}>{error}</p>}

            {loading ? (
                <p style={{ textAlign: "center" }}>Loading clusters...</p>
            ) : clusters.length === 0 ? (
                <p style={{ textAlign: "center" }}>No clusters found</p>
            ) : (
                <div style={{
                    border: "1px solid #ddd",
                    borderRadius: "5px",
                    overflow: "hidden"
                }}>
                    <table style={{ width: "100%", borderCollapse: "collapse" }}>
                        <thead>
                        <tr style={{ backgroundColor: "#f8f9fa" }}>
                            <th style={{ padding: "12px 15px", textAlign: "left", borderBottom: "1px solid #ddd" }}>Cluster Name</th>
                            <th style={{ padding: "12px 15px", textAlign: "left", borderBottom: "1px solid #ddd" }}>Status</th>
                            <th style={{ padding: "12px 15px", textAlign: "center", borderBottom: "1px solid #ddd" }}>Actions</th>
                        </tr>
                        </thead>
                        <tbody>
                        {clusters.map((cluster, index) => (
                            <tr key={index} style={{ borderBottom: "1px solid #ddd" }}>
                                <td style={{ padding: "12px 15px" }}>{cluster.name}</td>
                                <td style={{ padding: "12px 15px" }}>{cluster.status}</td>
                                <td style={{ padding: "12px 15px", textAlign: "center" }}>
                                    <button
                                        onClick={() => deleteCluster(cluster.name)}
                                        style={{
                                            padding: "5px 10px",
                                            backgroundColor: "#dc3545",
                                            color: "white",
                                            border: "none",
                                            borderRadius: "4px",
                                            cursor: "pointer",
                                        }}
                                    >
                                        Delete
                                    </button>
                                    <button
                                        onClick={() => downloadKubeconfig(cluster.name)}
                                        style={{
                                            padding: "5px 10px",
                                            backgroundColor: "#0078d4",
                                            color: "white",
                                            border: "none",
                                            borderRadius: "4px",
                                            cursor: "pointer",
                                        }}
                                    >
                                        Kubeconfig
                                    </button>
                                </td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                </div>
            )}
        </main>
    );
}