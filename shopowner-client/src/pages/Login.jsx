import { useState } from "react";
import Input from "../components/Input";
import Button from "../components/Button";
import { useMutation } from "@tanstack/react-query";
import { loginUser } from "../api/api";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "../store/authStore";
import { toast } from "react-toastify";

const Login = () => {
  const [formData, setFormData] = useState({ email: "", password: "" });
  const [err, setErr] = useState({ email: "", password: "" });
  const navigate = useNavigate();

  const { mutate, isPending, isError, isSuccess, error } = useMutation({
    mutationFn: loginUser,
    onSuccess: (data) => {
      toast("Logged in ✅");
      localStorage.setItem("token", data.token);
      localStorage.setItem("shopId", data.shopId);
      localStorage.setItem("menuId", data.menuId);
      localStorage.setItem("name", data.name);
      localStorage.setItem("email", data.email);
      useAuthStore.setState({ isAuth: true });
      navigate("/menu");
    },
    onError: (error) => {
      toast.error("Login failed: " + error.message);
    },
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const newError = { email: "", password: "" };

    if (!formData.email) newError.email = "Please enter your email";
    else if (!formData.email.includes("@"))
      newError.email = "Enter a valid email";
    if (!formData.password) newError.password = "Please enter your password";

    setErr(newError);
    if (newError.email || newError.password) return;

    mutate(formData);
  };

  return (
    <div className="flex h-screen bg-gradient-to-br from-gray-50 via-emerald-50/30 to-gray-100">
      <div className="h-screen w-screen flex items-center justify-center p-4">
        <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md overflow-hidden relative">
          {/* Header with gradient */}
          <div className="bg-gradient-to-br from-emerald-600 to-emerald-700 text-white p-8 relative overflow-hidden">
            <div className="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full -mr-16 -mt-16" />
            <div className="absolute bottom-0 left-0 w-24 h-24 bg-white/10 rounded-full -ml-12 -mb-12" />
            
            <div className="relative z-10 flex flex-col items-center">
              <div className="bg-white rounded-full h-16 w-16 mb-3 flex items-center justify-center shadow-lg">
                <span className="text-3xl">☕</span>
              </div>
              <h2 className="font-bold text-2xl mb-1">Welcome to Barista</h2>
              <p className="text-emerald-100 text-sm">Sign in to manage your café</p>
            </div>
          </div>

          {/* Form Section */}
          <div className="p-8 bg-gradient-to-b from-gray-50 to-white">
            <form onSubmit={handleSubmit} className="flex flex-col gap-4">
              <Input
                label="Email"
                type="email"
                name="email"
                placeholder="Enter your email"
                helperText={err.email}
                error={!!err.email}
                value={formData.email}
                onChange={handleChange}
              />
              <Input
                label="Password"
                type="password"
                name="password"
                placeholder="Enter your password"
                autocomplete="current-password"
                helperText={err.password}
                error={!!err.password}
                value={formData.password}
                onChange={handleChange}
              />
              <div className="flex justify-end w-full -mt-2">
                <button
                  type="button"
                  className="text-sm text-emerald-700 hover:text-emerald-800 hover:underline cursor-pointer font-medium transition-colors"
                >
                  Forgot Password?
                </button>
              </div>
              <Button type="submit" variant="primary">
                {isPending ? "Signing in..." : "Sign In"}
              </Button>
              <div className="flex justify-center gap-1 w-full text-sm mt-2">
                <p className="text-gray-600">Don't have an account?</p>
                <button
                  type="button"
                  onClick={() => navigate("/register")}
                  className="text-emerald-700 hover:text-emerald-800 hover:underline cursor-pointer font-medium transition-colors"
                >
                  Sign Up
                </button>
              </div>
            </form>
            
            <div className="mt-6 text-center">
              <p className="text-gray-400 text-xs">Powered by Barista</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;
