import { Tabs } from "../../constants/constants";
import { TabHeader } from "../TabHeader";

const Orders = ({ toggleSidebar }) => {
  return (
    <div className="p-5">
      <TabHeader toggleSidebar={toggleSidebar} tab={Tabs.ORDERS} />
    </div>
  );
};

export default Orders;
