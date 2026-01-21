import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
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
import Home from "../pages/Home.jsx";
import { useQuery } from "@tanstack/react-query";
import { verifyUser } from "../api/api.js";
import { useAuthStore } from "../store/authStore.js";
import CircleLoading from "../components/CircleLoading";

function Protected({ children }) {
  const token = localStorage.getItem("token");
  if (!token) return <Navigate to="/register" />;
  const { data, isLoading } = useQuery({
    queryKey: [token],
    queryFn: () => verifyUser(),
    onSuccess: (data) => {
      useAuthStore.setState({ isAuth: data });
      console.log("verified");
    },
  });

  if (isLoading) return <CircleLoading fullScreen text="Verifying..." />;
  if (!data) {
    return <Navigate to="/register" />;
  }

  return <>{children}</>;
}

function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route
          path="/menu"
          element={
            <Protected>
              <Menu />
            </Protected>
          }
        />
        <Route
          path="/orders"
          element={
            <Protected>
              <Orders />
            </Protected>
          }
        />
        <Route
          path="/tables"
          element={
            <Protected>
              <Tables />
            </Protected>
          }
        />
        <Route
          path="/history"
          element={
            <Protected>
              <History />
            </Protected>
          }
        />
        <Route
          path="/settings"
          element={
            <Protected>
              <Settings />
            </Protected>
          }
        />
        <Route
          path="/community"
          element={
            <Protected>
              <Community />
            </Protected>
          }
        />
        <Route
          path="/help"
          element={
            <Protected>
              <Help />
            </Protected>
          }
        />
        <Route path="/" element={<Home />} />
        <Route path="/*" element={<Navigate to="/" />} />
        <Route path="/test" element={<Test />} />
      </Routes>
    </BrowserRouter>
  );
}

export default AppRoutes;
