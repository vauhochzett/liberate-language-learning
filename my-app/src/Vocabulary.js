// Vocabulary.js
import React, { useState } from "react";

const Vocabulary = ({ word, translation }) => {
  const [correct, setCorrect] = useState(null);
  const [flipped, setFlipped] = useState(false);
  const [inputValue, setInputValue] = useState("");
  const [certificate, setCertificate] = useState(null); // New state for certificate

  const handleSubmit = async () => {
    setFlipped(true);

    const response = await fetch("/verifyWord", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        OriginalString: word,
        TranslatedString: inputValue,
      }),
    });

    if (response.ok) {
      const data = await response.json(); // Parse the JSON response
      setCorrect(data?.correct); // Set the 'correct' state based on the response
      setCertificate(data?.certificate); // Set the 'certificate' state based on the response
    } else {
      console.error("Failed to mark word as correct.");
      setCorrect(false);
    }
  };

  const handleInputChange = (e) => {
    setInputValue(e.target.value);
  };

  return (
    <div
      className={`vocabulary ${
        correct !== null ? (correct ? "correct" : "wrong") : ""
      } ${flipped ? "flipped" : ""}`}
    >
      <div className="card">
        <div className="card-front">
          <h2>{word}</h2>
          <input
            placeholder="insert translation"
            value={inputValue}
            onChange={handleInputChange}
          />
          <button onClick={handleSubmit}>Submit Input</button>
        </div>
        <div className="card-back">
          <p>{translation}</p>
          {correct === true && <span>✅</span>}
          {correct === false && <span>❌</span>}
          {certificate && (
            <a href={`https://ipfs.io/ipfs/${certificate}`}>View Certificate</a>
          )}
        </div>
      </div>
    </div>
  );
};

export default Vocabulary;
