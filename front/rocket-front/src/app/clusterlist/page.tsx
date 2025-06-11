"use client";

import React, { useState, useEffect } from "react";
import axios from "axios";
import Link from "next/link";
import Alert from "@mui/material/Alert";
import {
    Button,
    CircularProgress,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow
} from "@mui/material";
import { saveAs } from 'file-saver';
import {border} from "@mui/system";

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

            saveAs(blob,`${clusterName}-kubeconfig.yaml`);
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
            {error && <Alert severity="error">{error}"</Alert>}
            {data && <Alert severity="success">{data}</Alert>}
            <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "1rem" }}>
                <Button
                    onClick={fetchClusters}
                    variant="contained"
                >
                    Refresh
                </Button>
                <Link href="/cluster">
                    <Button
                        variant="contained"
                        color={"success"}
                    >
                        Create a new cluster
                    </Button>
                </Link>
            </div>
            {error && <Alert severity="error">{error}"</Alert>}
            {loading ? (
                <CircularProgress />
            ) : clusters.length === 0 ? (
                <Alert severity="error">No clusters found</Alert>
            ) : (
                    <TableContainer style={{border: "1px solid #CCC", borderRadius: "4px", overflow: "hidden"}}>
                        <Table >
                            <TableHead>
                                <TableRow>
                                    <TableCell >Cluster Name </TableCell>
                                    <TableCell  >Status</TableCell>
                                    <TableCell  >Action</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {clusters.map((cluster, index) => (
                                    <TableRow key={index} style={{ borderBottom: "1px solid #ddd" }}>
                                        <TableCell >{cluster.name}</TableCell>
                                        <TableCell >{cluster.status}</TableCell>
                                        <TableCell >
                                            <Button
                                                onClick={() => deleteCluster(cluster.name)}
                                                color="error"
                                                variant="contained"
                                            >
                                                Delete
                                            </Button>
                                            <Button
                                                onClick={() => downloadKubeconfig(cluster.name)}
                                                variant="contained"
                                            >
                                                Kubeconfig
                                            </Button>
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
            )}
        </main>
    );
}