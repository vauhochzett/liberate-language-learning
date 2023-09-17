import Cookies from "js-cookie";
import React, { useState } from "react";
import "./Vocabulary.css"; // Import the CSS

const Vocabulary = ({ word, onNext }) => {
  const [correct, setCorrect] = useState(null);
  const [flipped, setFlipped] = useState(false);
  const [inputValue, setInputValue] = useState("");
  const [certificate, setCertificate] = useState(null); // New state for certificate
  const [clickCount, setClickCount] = useState(0); // New state for click count
  const [correctWord, setCorrectWord] = useState(null); // New state for corrected word

  const handleSubmit = async () => {
    setFlipped(true);

    if (clickCount === 0) {
      // First click: Verify the word
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
        setCorrectWord(data?.CorrectWord);
      } else {
        console.error("Failed to mark word as correct.");
        setCorrect(false);
      }

      setClickCount(1); // Increment click count
    } else {
      // Second click: Move to the next word
      onNext(); // Call the onNext function to move to the next word
      setClickCount(0); // Reset click count
      setFlipped(false); // Reset the flipped state
      setInputValue(""); // Clear the input value
      setCorrect(null); // Reset the correct state
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
          <button onClick={handleSubmit}>
            {!clickCount ? "Submit" : "Next Word"}
          </button>
        </div>
        <div className="card-back">
          {correct === true && <span>✅</span>}
          {correct === false && <span>❌</span>}
          {correctWord && <p>Correct Word: {correctWord}</p>}{" "}
          {certificate && (
            <div>
              <h2>Congratulations!</h2>
              <p>
                You just earned the <i>English: Daily Conversation</i>{" "}
                certificate!
              </p>
              <p>Your NFT was minted and is available via IPFS:</p>
              <a href={`https://ipfs.io/ipfs/${certificate}`}>
                <img
                  src={`https://ipfs.io/ipfs/${certificate}`}
                  alt="Vocabulary Book"
                  className="certificate"
                />
              </a>
            </div>
          )}
        </div>
        <p>
          Track your learning progress via micro-certificates on the blockchain.
        </p>
      </div>
    </div>
  );
};

export default Vocabulary;
