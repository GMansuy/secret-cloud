"use client";

import { useState } from "react";
import axios from "axios";

export default function Home() {
  const [data, setData] = useState<string>("");
  const [error, setError] = useState<string>("");

  const callBackend = async () => {
    setError("");
    setData("");
    try {
      const response = await axios.get("http://localhost:8080/");
      setData(JSON.stringify(response.data, null, 2));
    } catch (err) {
      setError("Failed to fetch from backend");
      console.error(err);
    }
  };

    const postCluster = async () => {
        setError("");
        setData("");
        axios.post('http://localhost:8080/cluster', {
            name: 'test-cluster',
            controlplaneMachineCount: 1,
            workerMachineCount: 2
        })
            .then(function (response) {
                console.log(response);
            })
            .catch(function (error) {
                console.log(error);
            });
     //   try {
     //       const response = await axios.post("http://localhost:8080/cluster",
     //           //{
     //           //    name: "example-cluster",
     //           //    controlplaneMachineCount: 1,
     //           //    workerMachineCount: 1,
     //           //}
     //       );
     //       setData(JSON.stringify(response.data, null, 2));
     //   } catch (err) {
     //       setError("Failed to fetch from backend");
     //       console.error(err);
     //   }
    };

  return (
      <main style={{ maxWidth: 600, margin: "2rem auto", fontFamily: "Arial, sans-serif" }}>
        <h1>Call Backend Example</h1>
        <button
            onClick={callBackend}
            style={{
              padding: "0.5rem 1rem",
              fontSize: "1rem",
              cursor: "pointer",
              marginBottom: "1rem",
            }}
        >
          Call Backend
        </button>

        {error && <p style={{ color: "red" }}>{error}</p>}

        {data && (
            <pre
                style={{
                  background: "#000000",
                  padding: "1rem",
                  borderRadius: "5px",
                  whiteSpace: "pre-wrap",
                  maxHeight: "300px",
                  overflow: "auto",
                }}
            >
          {data}
        </pre>
        )}
          <button
              onClick={postCluster}
              style={{
                  padding: "0.5rem 1rem",
                  fontSize: "1rem",
                  cursor: "pointer",
                  marginBottom: "1rem",
              }}
          >
              post Cluster
          </button>
          {error && <p style={{ color: "red" }}>{error}</p>}

          {data && (
              <pre
                  style={{
                      background: "#000000",
                      padding: "1rem",
                      borderRadius: "5px",
                      whiteSpace: "pre-wrap",
                      maxHeight: "300px",
                      overflow: "auto",
                  }}
              >
          {data}
        </pre>
          )}
      </main>
  );
}
