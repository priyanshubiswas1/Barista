import React, { useState } from "react";
import Input from "../components/Input";
import Button from "../components/Button";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState({ email: "", password: "" });

  const handleSubmit = (e) => {
    e.preventDefault();

    const newError = { email: "", password: "" };

    if (!email) newError.email = "Email is required";
    else if (!email.includes("@")) newError.email = "Enter a valid email";

    if (!password) newError.password = "Password is required";

    setError(newError);

    if (newError.email || newError.password) return;

    // Submit logic
    alert("Logged in ✅");
  };

  return (
    <div className="w-100 flex flex-col gap-2 bg-white p-5 rounded-md">
      <form onSubmit={handleSubmit} className="flex flex-col gap-2" noValidate>
        <Input
          label="Email"
          type="email"
          placeholder="Enter your email"
          helperText={error.email}
          error={!!error.email}
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <Input
          label="Password"
          type="password"
          placeholder="Enter your password"
          helperText={error.password}
          error={!!error.password}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
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
          Continue
        </Button>
      </form>
    </div>
  );
};

export default Login;
