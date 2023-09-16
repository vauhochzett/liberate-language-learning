// Vocabulary.js
import React, { useState } from "react";

const Vocabulary = ({ word, translation }) => {
  const [correct, setCorrect] = useState(null);
  const [flipped, setFlipped] = useState(false);
  const [inputValue, setInputValue] = useState(""); // New state for input value

  const handleSubmit = async () => {
    setCorrect(true);
    setFlipped(true);

    const response = await fetch("https://httpbin.org/post", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ key: inputValue }),
    });

    if (response.ok) {
      console.log("Word marked as correct.");
    } else {
      console.error("Failed to mark word as correct.");
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
            onChange={handleInputChange} // Update input value on change
          />
          <button onClick={handleSubmit}>✔️</button>
        </div>
        <div className="card-back">
          <p>{translation}</p>
        </div>
      </div>
    </div>
  );
};

export default Vocabulary;
