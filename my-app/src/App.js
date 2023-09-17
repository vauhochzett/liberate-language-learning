import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import CheckCertificate from "./CheckCertificate"; // Import your subpage
import Main from "./Main"; // Import your existing Vocabulary component or page

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Main />} />
        <Route path="/check-certs" element={<CheckCertificate />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
