import React from "react";
import Button from "../components/Button";
import { Plus } from "phosphor-react";
import LeftSidebar from "../components/LeftSidebar.jsx";
import Login from "./Login.jsx";

const Test = () => {
  return (
    <div className="flex h-screen">
      {/* <LeftSidebar /> */}
      <div className="h-screen w-screen flex flex-col gap-6 items-center justify-center bg-white">
        {/* <div className="inline-flex flex-col gap-4 p-6 bg-white rounded-lg">
          <h3 className="text-base font-semibold text-gray-800">Text Only :</h3>
          <div className="flex gap-2">
            <Button>Primary</Button>
            <Button variant="outline">Outline</Button>
            <Button variant="subtle">Subtle</Button>
          </div>
          <h3 className="text-base font-semibold text-gray-800">
            Text & Icon :
          </h3>
          <div className="flex gap-2">
            <Button variant="primary" icon={Plus}>
              Add Item
            </Button>
            <Button variant="outline" icon={Plus}>
              Add Item
            </Button>
            <Button variant="subtle" icon={Plus}>
              Add Item
            </Button>
          </div>
          <h3 className="text-base font-semibold text-gray-800">Icon Only :</h3>
          <div className="flex gap-2">
            <Button variant="primary" icon={Plus}></Button>
            <Button variant="outline" icon={Plus}></Button>
            <Button variant="subtle" icon={Plus}></Button>
          </div>
          <h3 className="text-base font-semibold text-gray-800">Disabled :</h3>
          <div className="flex gap-2">
            <Button variant="primary" icon={Plus} disabled>
              Add Item
            </Button>
            <Button variant="outline" icon={Plus} disabled>
              Add Item
            </Button>
            <Button variant="subtle" icon={Plus} disabled>
              Add Item
            </Button>
          </div>
        </div> */}
        <header className="flex flex-col gap-2 items-center">
          <h1 className="text-2xl font-semibold text-gray-800">
            Sign in to MyCafeApp
          </h1>
          {/* <h2 className="text-gray-400 text-base font-normal">
            Enter your details to continue
          </h2> */}
        </header>
        <Login></Login>
      </div>
    </div>
  );
};

export default Test;
