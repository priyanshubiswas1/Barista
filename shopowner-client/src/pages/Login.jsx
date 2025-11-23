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
    <div className="flex h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <div className="h-screen w-screen flex flex-col gap-6 items-center justify-center">
        <header className="flex flex-col gap-2 items-center">
          <div className="w-12 h-12 bg-gray-800 rounded-lg mb-2"></div>
          <h1 className="text-3xl font-bold text-gray-900">
            Sign in to Barista
          </h1>
          <p className="text-gray-600">Welcome back! Please enter your details</p>
        </header>
        <div className="w-96 flex flex-col gap-3 bg-white p-8 rounded-2xl shadow-xl border-2 border-gray-200">
          <form onSubmit={handleSubmit} className="flex flex-col gap-3">
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
            <div className="flex justify-end w-full">
              <button
                type="button"
                className="text-sm text-emerald-700 hover:text-emerald-800 hover:underline cursor-pointer font-medium"
              >
                Forgot Password?
              </button>
            </div>
            <Button type="submit" variant="primary">
              {isPending ? "Logging in..." : "Continue"}
            </Button>
            <div className="flex justify-center gap-1 w-full text-sm">
              <p className="text-gray-600">Don't have an account?</p>
              <button
                type="button"
                className="text-emerald-700 hover:text-emerald-800 hover:underline cursor-pointer font-medium"
              >
                Sign Up
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default Login;
