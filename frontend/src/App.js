import React, { useState } from "react";
import axios from "axios";

function App() {
  const [file, setFile] = useState(null);
  const [uploadQuery, setUploadQuery] = useState("");
  const [chatQuery, setChatQuery] = useState("");
  const [uploadResponses, setUploadResponses] = useState([]); 
  const [chatResponses, setChatResponses] = useState([]);
  const [activeSection, setActiveSection] = useState("upload");
  const [responseVisible, setResponseVisible] = useState(false);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleUploadQueryChange = (e) => {
    setUploadQuery(e.target.value);
  };

  const handleChatQueryChange = (e) => {
    setChatQuery(e.target.value);
  };

  const handleUpload = async () => {
    if (!file || !uploadQuery) {
      alert("Please select a file and enter a question.");
      return;
    }
  
    const formData = new FormData();
    formData.append("file", file);
    formData.append("question", uploadQuery);
  
    try {
      const res = await axios.post("http://localhost:8080/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      setUploadResponses(prev => [...prev, { question: uploadQuery, answer: res.data.answer || "No response received." }]); 
      setResponseVisible(true); 
      setUploadQuery(""); 
    } catch (error) {
      console.error("Error uploading file:", error);
      setUploadResponses(prev => [...prev, { question: uploadQuery, answer: "Error uploading file." }]);
      setResponseVisible(true); 
    }
  };
  
  const handleChat = async () => {
    if (!chatQuery) {
      alert("Please enter a question.");
      return;
    }
  
    try {
      const res = await axios.post("http://localhost:8080/chat", { query: chatQuery });
      setChatResponses(prev => [...prev, { question: chatQuery, answer: res.data.answer || "No response received." }]); 
      setResponseVisible(true); 
      setChatQuery(""); 
    } catch (error) {
      console.error("Error querying chat:", error);
      setChatResponses(prev => [...prev, { question: chatQuery, answer: "Error querying chat." }]); 
      setResponseVisible(true); 
    }
  };
  
  const handleReplaceFile = () => {
    setFile(null); 
  };

  return (
    <div
      style={{
        maxWidth: "1000px",
        margin: "0 auto",
        padding: "20px",
        fontFamily: "Arial, sans-serif",
        display: "flex",
        flexDirection: "column",
        height: "75vh",
        boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
        borderRadius: "8px",
        backgroundColor: "#0B192C",
      }}
    >
      {/* Navigation */}
<div
  style={{
    display: "flex",
    justifyContent: "center",
    marginBottom: "20px",
    gap: "15px",
  }}
>
  <button
    onClick={() => setActiveSection("upload")}
    style={{
      flex: 1,
      padding: "12px 20px",
      backgroundColor: activeSection === "upload" ? "#5A72A0" : "#f1f1f1",
      color: activeSection === "upload" ? "#fff" : "#333",
      border: "none",
      borderRadius: "8px",
      cursor: "pointer",
      fontSize: "16px",
      fontWeight: "600",
      transition: "background-color 0.3s ease, transform 0.2s ease-in-out",
    }}
    onMouseLeave={(e) => e.target.style.transform = "scale(1)"}
  >
    <i className="fas fa-upload" style={{ marginRight: "8px" }}></i>
    EcoSphere Analyze
  </button>
  <button
    onClick={() => setActiveSection("chat")}
    style={{
      flex: 1,
      padding: "12px 20px",
      backgroundColor: activeSection === "chat" ? "#5A72A0" : "#f1f1f1",
      color: activeSection === "chat" ? "#fff" : "#333",
      border: "none",
      borderRadius: "8px",
      cursor: "pointer",
      fontSize: "16px",
      fontWeight: "600",
      transition: "background-color 0.3s ease, transform 0.2s ease-in-out",
    }}
    onMouseEnter={(e) => e.target.style.transform = "scale(1.05)"}
    onMouseLeave={(e) => e.target.style.transform = "scale(1)"}
  >
    <i className="fas fa-comments" style={{ marginRight: "8px" }}></i>
    Chat With AI
  </button>
</div>

      {/* Response Section */}
      <div
        style={{
          flex: 1,
          overflowY: "auto",
          padding: "10px",
          border: "1px solid #ccc",
          borderRadius: "4px",
          backgroundColor: "white",
        }}
      >
        {activeSection === "upload" ? (
          <div>
            {!uploadResponses.length && (
              <p
                style={{
                  color: "#BCCCDC",
                  fontSize: "24px",
                  fontWeight: "bold",
                  textAlign: "center",
                  marginBottom: "25px",
                  marginTop: "200px",
                }}
              >
                Start Upload and Analyze....
              </p>
            )}
            {uploadResponses.map((response, index) => (
              <div
                key={index}
                style={{
                  marginBottom: "10px",
                  padding: "10px",
                  border: "1px solid #021526",
                  borderRadius: "4px",
                  backgroundColor: "#E2E2B6",
                  color: "black",
                  opacity: responseVisible ? 1 : 0,
                  transition: "opacity 1s ease-in-out",
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "flex-end",
                  marginRight: "0",
                }}
              >
                <div
                  style={{
                    marginBottom: "5px",
                    textAlign: "right",
                    fontWeight: "bold",
                    color: "#333",
                    marginLeft: "auto", 
                  }}
                >
                  <strong>Question:</strong> {response.question}
                </div>
                <div
                  style={{
                    textAlign: "left",
                    backgroundColor: "#A3C9F1",
                    padding: "10px",
                    borderRadius: "8px",
                    color: "#333",
                    minWidth: "50%",
                  }}
                >
                  <strong>AI Response:</strong> {response.answer}
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div>
            {!chatResponses.length && (
              <p
                style={{
                  color: "#BCCCDC",
                  fontSize: "24px",
                  fontWeight: "bold",
                  textAlign: "center",
                  marginBottom: "20px",
                  marginTop: "200px",
                }}
              >
                Start Talk With AI....
              </p>
            )}
            {chatResponses.map((response, index) => (
              <div
                key={index}
                style={{
                  marginBottom: "10px",
                  padding: "10px",
                  border: "1px solid #021526",
                  borderRadius: "4px",
                  backgroundColor: "#E2E2B6",
                  color: "black",
                  opacity: responseVisible ? 1 : 0,
                  transition: "opacity 1s ease-in-out",
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "flex-end",
                  marginRight: "0",
                }}
              >
                <div
                  style={{
                    marginBottom: "5px",
                    textAlign: "right",
                    fontWeight: "bold",
                    color: "#333",
                    marginLeft: "auto", 
                  }}
                >
                  <strong>Question:</strong> {response.question}
                </div>
                <div
                  style={{
                    textAlign: "left",
                    backgroundColor: "#A3C9F1",
                    padding: "10px",
                    borderRadius: "8px",
                    color: "#333",
                    minWidth: "50%",
                  }}
                >
                  <strong>AI Response:</strong> {response.answer}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Input Section */}
      <div
        style={{
          marginTop: "10px",
          display: "flex",
          gap: "10px",
          alignItems: "center",
        }}
      >
        {activeSection === "upload" && (
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: "10px",
            }}
          >
            {file ? (
              <span
                style={{
                  flex: 1,
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                  whiteSpace: "nowrap",
                  fontSize: "14px",
                  color: "white",
                }}
              >
                {file.name}
              </span>
            ) : (
              <label
                style={{
                  display: "flex",
                  alignItems: "center",
                  padding: "10px",
                  border: "1px solid #ccc",
                  borderRadius: "4px",
                  backgroundColor: "#f0f0f0",
                  cursor: "pointer",
                }}
              >
                <i className="fas fa-upload" style={{ fontSize: "20px" }}></i>
                <input
                  type="file"
                  onChange={handleFileChange}
                  style={{
                    display: "none",
                  }}
                />
              </label>
            )}

            {file && (
              <button
                onClick={handleReplaceFile}
                style={{
                  background: "none",
                  border: "none",
                  cursor: "pointer",
                }}
              >
                <i className="fas fa-redo" style={{ fontSize: "12px", color: "white" }}></i>
              </button>
            )}
          </div>
        )}

        <input
          type="text"
          value={activeSection === "upload" ? uploadQuery : chatQuery}
          onChange={activeSection === "upload" ? handleUploadQueryChange : handleChatQueryChange}
          placeholder={activeSection === "upload" ? "Enter Your Question Here..." : "Enter Your Question Here..."}
          style={{
            flex: 1,
            padding: "10px",
            border: "1px solid #ccc",
            borderRadius: "4px",
          }}
        />
        <button
          onClick={activeSection === "upload" ? handleUpload : handleChat}
          style={{
            padding: "10px",
            backgroundColor: "#5A72A0",
            color: "white",
            border: "none",
            borderRadius: "4px",
            cursor: "pointer",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
        <i
          className="fas fa-paper-plane"
          style={{
            fontSize: "20px",
          }}
        ></i>
        </button>
      </div>
    </div>
  );
}

export default App;
