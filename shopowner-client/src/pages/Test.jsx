import React from "react";
import Button from "../components/Button";
import { Plus } from "phosphor-react";
import LeftSidebar from "../components/LeftSidebar.jsx";
import Login from "./Login.jsx";
import SearchableSelect from "../components/SearchableSelect.jsx";
import ImageUploader from "../components/ImageUploader.jsx";

const Test = () => {
  return (
    <div className="flex h-screen items-center justify-center">
      <ImageUploader label="Image" details="Supports JPEGs & PNGS under 5mb." />
    </div>
  );
};

export default Test;
