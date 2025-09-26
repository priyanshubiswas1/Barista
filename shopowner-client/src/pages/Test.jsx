import React from "react";
import Button from "../components/Button";
import { Plus } from "phosphor-react";
import Sidebar from "../components/Sidebar";

const Test = () => {
  return (
    <div className="flex h-screen">
      <Sidebar></Sidebar>
      <div className="h-screen w-screen flex items-center justify-center bg-gray-100">
        <div className="inline-flex flex-col gap-4 p-6 bg-white rounded-lg">
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
        </div>
      </div>
    </div>
  );
};

export default Test;
