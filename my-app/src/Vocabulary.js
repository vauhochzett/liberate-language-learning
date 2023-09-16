import React, { useState } from "react";
import Cookies from "js-cookie";
import "./Vocabulary.css"; // Import the CSS

const Vocabulary = ({ word }) => {
  const [correct, setCorrect] = useState(null);
  const [flipped, setFlipped] = useState(false);
  const [inputValue, setInputValue] = useState("");
  const [certificate, setCertificate] = useState(null); // New state for certificate
  const [correctWord, setCorrectWord] = useState(null); // New state for corrected word

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
        AccId: Cookies.get("accId"),
        Language: "de",
      }),
    });

    if (response.ok) {
      const data = await response.json(); // Parse the JSON response
      setCorrect(data?.Correct); // Set the 'correct' state based on the response
      setCertificate(data?.Certificate); // Set the 'certificate' state based on the response
      setCorrectWord(data?.CorrectWord); // Set the 'correctWord' state based on the response
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
          {correct === true && <span>✅</span>}
          {correct === false && <span>❌</span>}
          {correctWord && <p>Correct Word: {correctWord}</p>}{" "}
          {certificate && (
            <a href={`https://ipfs.io/ipfs/${certificate}`}>
              <img
                src={`https://ipfs.io/ipfs/${certificate}`}
                alt="Vocabulary Book"
                className="certificate"
              />
            </a>
          )}
        </div>
      </div>
    </div>
  );
};

export default Vocabulary;
