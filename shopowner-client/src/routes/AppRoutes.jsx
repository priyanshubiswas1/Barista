import { BrowserRouter, Route, Routes } from "react-router-dom";
import Test from "../pages/Test.jsx";
import Menu from "../pages/Menu.jsx";
import Orders from "../pages/Orders.jsx";
import Tables from "../pages/Tables.jsx";
import History from "../pages/History.jsx";
import Settings from "../pages/Settings.jsx";
import Community from "../pages/Community.jsx";
import Help from "../pages/Help.jsx";
import Login from "../pages/Login.jsx";
import Register from "../pages/Register.jsx";

function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/menu" element={<Menu />} />
        <Route path="/orders" element={<Orders />} />
        <Route path="/tables" element={<Tables />} />
        <Route path="/history" element={<History />} />
        <Route path="/settings" element={<Settings />} />
        <Route path="/community" element={<Community />} />
        <Route path="/help" element={<Help />} />
        <Route path="/test" element={<Test />} />
      </Routes>
    </BrowserRouter>
  );
}

export default AppRoutes;
