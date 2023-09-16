import React, { useEffect, useCallback, useState } from "react";
import "./App.css";
import { vocabularies } from "./vocabularies";
import Vocabulary from "./Vocabulary";
import Cookies from "js-cookie";

function App() {
  const [currentIndex, setCurrentIndex] = useState(0); // New state for current index

  const checkAccountData = useCallback(async () => {
    const accId = Cookies.get("accId");
    const pubKey = Cookies.get("pubKey");

    if (!accId || !pubKey) {
      const { AccId, PubKey, PrivKey } = await createKey();
      Cookies.set("accId", AccId);
      Cookies.set("pubKey", PubKey);
      console.log(`Private Key: ${PrivKey}`);
    }
  }, []);

  useEffect(() => {
    checkAccountData();
  }, [checkAccountData]);

  const createKey = async () => {
    const response = await fetch("/createKey", {
      method: "POST",
    });

    if (response.ok) {
      return await response.json();
    } else {
      throw new Error("Failed to create key");
    }
  };

  const handleNextVocab = () => {
    if (currentIndex < vocabularies.length - 1) {
      setCurrentIndex(currentIndex + 1);
    }
  };

  return (
    <div className="App">
      <h1>Vocabulary App</h1>
      <div className="vocab-list">
        <Vocabulary
          key={currentIndex}
          word={vocabularies[currentIndex].word}
          onNext={handleNextVocab} // Pass the handleNextVocab function as a prop
        />
      </div>
    </div>
  );
}

export default App;
