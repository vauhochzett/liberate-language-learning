// App.js
import React, { useState } from "react";
import "./App.css";
import { vocabularies } from "./vocabularies";
import Vocabulary from "./Vocabulary";

function App() {
  const [myBlockchainkey, setBlockChainkey] = useState("");

  const requestKey = async () => {
    const response = await fetch("https://httpbin.org/get");

    if (response.ok) {
      setBlockChainkey();
      console.log("Got the key");
    } else {
      console.error("Getting key failed");
    }
  };
  return (
    <div className="App">
      <h1>Vocabulary App</h1>
      Get key:
      {myBlockchainkey}
      <button onClick={requestKey}>Hallo</button>
      <div className="vocab-list">
        {vocabularies.map((vocab, index) => {
          return (
            <Vocabulary
              key={index}
              word={vocab.word}
              translation={vocab.translation}
            />
          );
        })}
      </div>
    </div>
  );
}

export default App;
