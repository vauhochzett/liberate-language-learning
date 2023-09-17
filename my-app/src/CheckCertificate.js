import React, { useState } from "react";

const CheckCertificate = () => {
  const [accId, setAccId] = useState("");
  const [certId, setCertId] = useState("");
  const [serial, setSerial] = useState("");
  const [isValid, setIsValid] = useState(null);

  const mobileStyle = {
    container: {
      padding: "20px",
      textAlign: "center",
    },
    input: {
      width: "80%",
      padding: "10px",
      margin: "10px 0",
    },
    button: {
      padding: "10px 20px",
      fontSize: "16px",
    },
  };

  const handleCheckCertificate = async () => {
    const payload = {
      AccId: accId,
      CertId: certId,
      Serial: serial,
    };

    try {
      const response = await fetch("/checkCert", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        const data = await response.json();
        setIsValid(data.Valid);
      } else {
        console.error("Failed to check certificate.");
      }
    } catch (error) {
      console.error("An error occurred:", error);
    }
  };

  return (
    <div style={mobileStyle.container}>
      <h1>Check Certificate</h1>
      <div>
        <label>Account ID:</label>
        <input
          type="text"
          value={accId}
          onChange={(e) => setAccId(e.target.value)}
          style={mobileStyle.input}
        />
      </div>
      <div>
        <label>Certificate ID:</label>
        <input
          type="text"
          value={certId}
          onChange={(e) => setCertId(e.target.value)}
          style={mobileStyle.input}
        />
      </div>
      <div>
        <label>Serial:</label>
        <input
          type="text"
          value={serial}
          onChange={(e) => setSerial(e.target.value)}
          style={mobileStyle.input}
        />
      </div>
      <button onClick={handleCheckCertificate} style={mobileStyle.button}>
        Check
      </button>

      {isValid !== null && (
        <div>
          {isValid ? (
            <p>The certificate is valid.</p>
          ) : (
            <p>The certificate is not valid.</p>
          )}
        </div>
      )}
    </div>
  );
};

export default CheckCertificate;
