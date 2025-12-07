import React from "react";

const CircleLoading = ({
  size = "md",
  fullScreen = false,
  text = "Loading...",
}) => {
  const sizeClasses = {
    sm: "w-6 h-6 border-2",
    md: "w-10 h-10 border-4",
    lg: "w-16 h-16 border-4",
  };

  const spinner = (
    <div className="flex flex-col items-center justify-center gap-4">
      <div
        className={`${sizeClasses[size]} animate-spin rounded-full border-indigo-200 border-t-indigo-700`}
      />
      {text && (
        <p className="text-indigo-700 font-medium animate-pulse">{text}</p>
      )}
    </div>
  );

  if (fullScreen) {
    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm">
        <div className="bg-white rounded-2xl shadow-2xl p-8 flex flex-col items-center justify-center min-w-[200px]">
          {spinner}
        </div>
      </div>
    );
  }

  return spinner;
};

export default CircleLoading;
