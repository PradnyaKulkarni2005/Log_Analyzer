import { useState } from "react";

function App() {
  const [logs, setLogs] = useState("");
  const [result, setResult] = useState(null);

  const analyzeLogs = async () => {
    const res = await fetch("http://localhost:9090/analyze", {
      method: "POST",
      body: logs,
    });

    const data = await res.json();
    setResult(data);
  };

  return (
    <div style={{ padding: "20px" }}>
      <h1>🚀 Log Analyzer</h1>

      <textarea
        rows="10"
        cols="80"
        placeholder="Paste logs here..."
        value={logs}
        onChange={(e) => setLogs(e.target.value)}
      />

      <br /><br />

      <button onClick={analyzeLogs}>Analyze</button>

      {result && (
        <div style={{ marginTop: "20px" }}>
          <h2>Results</h2>
          <pre>{JSON.stringify(result, null, 2)}</pre>
        </div>
      )}
    </div>
  );
}

export default App;