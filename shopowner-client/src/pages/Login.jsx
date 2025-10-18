import { useState } from "react";
import Input from "../components/Input";
import Button from "../components/Button";
import { useMutation } from "@tanstack/react-query";
import { loginUser } from "../api/api";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "../store/authStore";

const Login = () => {
  const [formData, setFormData] = useState({ email: "", password: "" });
  const [err, setErr] = useState({ email: "", password: "" });
  const navigate = useNavigate();

  const { mutate, isPending, isError, isSuccess, error } = useMutation({
    mutationFn: loginUser,
    onSuccess: (data) => {
      // alert("Logged in ✅");
      localStorage.setItem("token", data.token);
      useAuthStore.setState({ isAuth: true });
      navigate("/menu");
    },
    onError: (error) => {
      console.log("Login failed: ", error);
    },
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const newError = { email: "", password: "" };

    if (!formData.email) newError.email = "Email is required";
    else if (!formData.email.includes("@"))
      newError.email = "Enter a valid email";
    if (!formData.password) newError.password = "Password is required";

    setErr(newError);
    if (newError.email || newError.password) return;

    mutate(formData);
  };

  return (
    <div className="flex h-screen">
      <div className="h-screen w-screen flex flex-col gap-6 items-center justify-center bg-white">
        <header className="flex flex-col gap-2 items-center">
          <h1 className="text-2xl font-semibold text-gray-800">
            Sign in to MyCafeApp
          </h1>
        </header>
        <div className="w-100 flex flex-col gap-2 bg-white p-5 rounded-md">
          <form onSubmit={handleSubmit} className="flex flex-col gap-2">
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
                className=" text-sm text-gray-800 hover:underline cursor-pointer"
              >
                Forgot Password?
              </button>
            </div>
            <Button type="submit" variant="primary">
              {isPending ? "Logging in..." : "Continue"}
            </Button>
            <div className="flex justify-center gap-1 w-full">
              <p className="text-gray-500">Don't have an account?</p>
              <button
                type="button"
                className=" text-sm text-gray-800 hover:underline cursor-pointer"
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
