import React, { useEffect, useCallback } from "react";
import "./App.css";
import { vocabularies } from "./vocabularies";
import Vocabulary from "./Vocabulary";
import Cookies from "js-cookie";

function App() {
  const checkAccountData = useCallback(async () => {
    const accId = Cookies.get("accId");
    const pubKey = Cookies.get("pubKey");

    if (!accId || !pubKey) {
      const { AccId, PubKey, PrivKey } = await createKey();
      Cookies.set("accId", AccId);
      Cookies.set("pubKey", PubKey);
      console.log(`Private Key: ${PrivKey}`);
    }
  }, []); // Empty dependency array means this function is memoized and won't change on re-renders

  useEffect(() => {
    checkAccountData();
  }, [checkAccountData]); // Now checkAccountData is a dependency, but it's memoized so useEffect won't run repeatedly

  const createKey = async () => {
    // Call your backend API here
    const response = await fetch("/createKey", {
      method: "POST",
    });

    if (response.ok) {
      return await response.json();
    } else {
      throw new Error("Failed to create key");
    }
  };
  return (
    <div className="App">
      <h1>Vocabulary App</h1>
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
