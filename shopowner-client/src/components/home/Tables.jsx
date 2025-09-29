import { TabHeader } from "../TabHeader";
import { Tabs } from "../../constants/constants";

const Tables = ({ toggleSidebar }) => {
  return (
    <div className="p-5">
      <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.TABLES} />{" "}
    </div>
  );
};

export default Tables;
