import AppRoutes from "./routes/AppRoutes.jsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ToastContainer, Slide } from "react-toastify";

const App = () => {
  const queryClient = new QueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <AppRoutes />
      <ToastContainer
        position="bottom-center"
        autoClose={3000}
        limit={2}
        hideProgressBar
        draggable
        pauseOnHover
        toastClassName={() =>
          "flex px-3 py-2 gap-8 bg-gray-950 rounded-md shadow-sm items-center text-sm text-white font-medium"
        }
        closeButton={false}
        transition={Slide}
        style={{ bottom: "2rem" }}
      />
    </QueryClientProvider>
  );
};

export default App;
