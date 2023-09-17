import Cookies from "js-cookie";
import React, { useCallback, useEffect, useState } from "react";
import "./App.css";
import Vocabulary from "./Vocabulary";
import { vocabularies } from "./vocabularies";

import Card from 'react-bootstrap/Card';
import Col from 'react-bootstrap/Col';
import Container from "react-bootstrap/Container";
import Row from 'react-bootstrap/Row';

function App() {
  const [currentIndex, setCurrentVocabIndex] = useState(0); // New state for current index

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

  const moveToNextVocab = () => {
    if (currentIndex < vocabularies.length - 1) {
      setCurrentVocabIndex(currentIndex + 1); // Move to the next vocab
    }
  };

  return (
    <Container className="App">
      <Row>
        <Col>
          <h1>Chain the Gain!</h1>
        </Col>
      </Row>
      <Row>
        <Col></Col>
        <Col xs={12} md={8}>
          <Card className="vocab-list">
            <Vocabulary
              key={currentIndex}
              word={vocabularies[currentIndex].word}
              onNext={moveToNextVocab} // Pass the handleNextVocab function as a prop
            />
          </Card>
        </Col>
        <Col></Col>
      </Row>
    </Container>
  );
}

export default App;
