import React from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import Menu from "../pages/Menu";

function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Menu />} />
      </Routes>
    </BrowserRouter>
  );
}

export default AppRoutes;
