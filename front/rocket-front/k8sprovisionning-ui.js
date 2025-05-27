// A simple Next.js page that includes:
// - A form to provision Kubernetes clusters
// - A list to view provisioned clusters and their statuses
// - API integration via Axios

'use client';

import { useState, useEffect } from 'react';
import axios from 'axios';

export default function ClusterDashboard() {
    const [clusters, setClusters] = useState([]);
    const [name, setName] = useState('');
    const [controlPlaneNodes, setControlPlaneNodes] = useState(1);
    const [workerNodes, setWorkerNodes] = useState(2);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;

    const fetchClusters = async () => {
        try {
            const res = await axios.get('${API_BASE_URL}/clusters');
            setClusters(res.data);
        } catch (err) {
            console.error('Failed to fetch clusters', err);
        }
    };

    useEffect(() => {
        fetchClusters();
    }, []);

    const handleSubmit = async () => {
        setLoading(true);
        setError('');
        try {
            await axios.post('${API_BASE_URL}/clusters', { name, controlPlaneNodes, workerNodes });
            await fetchClusters();
            setName('');
        } catch (err) {
            setError('Provisioning failed');
        } finally {
            setLoading(false);
        }
    };

    return (
        <main className="max-w-2xl mx-auto py-10 px-4">
            <h1 className="text-2xl font-bold mb-6">Kubernetes Cluster Provisioning</h1>

            <div className="bg-white shadow-md rounded-lg p-6 mb-8">
                <h2 className="text-xl font-semibold mb-4">New Cluster</h2>
                <input
                    className="border p-2 mb-2 w-full"
                    placeholder="Cluster Name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />
                <input
                    className="border p-2 mb-2 w-full"
                    type="number"
                    min={1}
                    placeholder="Control Plane Nodes"
                    value={controlPlaneNodes}
                    onChange={(e) => setControlPlaneNodes(Number(e.target.value))}
                />
                <input
                    className="border p-2 mb-4 w-full"
                    type="number"
                    min={1}
                    placeholder="Worker Nodes"
                    value={workerNodes}
                    onChange={(e) => setWorkerNodes(Number(e.target.value))}
                />
                <button
                    className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
                    onClick={handleSubmit}
                    disabled={loading}
                >
                    {loading ? 'Provisioning...' : 'Provision Cluster'}
                </button>
                {error && <p className="text-red-500 mt-2">{error}</p>}
            </div>

            <div>
                <h2 className="text-xl font-semibold mb-4">Provisioned Clusters</h2>
                <ul className="space-y-3">
                    {clusters.map((cluster) => (
                        <li
                            key={cluster.id}
                            className="border p-4 rounded-md flex justify-between items-center"
                        >
                            <div>
                                <p className="font-medium">{cluster.name}</p>
                                <p className="text-sm text-gray-600">
                                    Control Plane: {cluster.controlPlaneNodes} | Workers: {cluster.workerNodes}
                                </p>
                            </div>
                            <span className={`text-sm px-2 py-1 rounded-full ${
                                cluster.status === 'ready' ? 'bg-green-200 text-green-800' : 'bg-yellow-200 text-yellow-800'
                            }`}>
                {cluster.status}
              </span>
                        </li>
                    ))}
                </ul>
            </div>
        </main>
    );
}
