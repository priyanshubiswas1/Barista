import { useState, useRef, useEffect } from "react";
import { ImageSquare, TrashSimple } from "phosphor-react";
import Button from "./Button";

const ImageUploader = ({
  label = "Image",
  details = "Max file size: 5MB",
  onChange,
  value = null,
}) => {
  const [image, setImage] = useState(null);
  const fileInputRef = useRef(null);

  useEffect(() => {
    if (value === null) {
      setImage(null);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    } else if (value instanceof File) {
      const url = URL.createObjectURL(value);
      setImage(url);
      return () => {
        URL.revokeObjectURL(url);
      };
    }
  }, [value]);

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      if (onChange) onChange(file);
    }
  };

  const handleImageClick = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  const handleRemove = (e) => {
    e.stopPropagation();
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
    if (onChange) onChange(null);
  };

  return (
    <div className="flex flex-col gap-2">
      {/* Label */}
      <span className="text-sm font-medium text-gray-800">{label}</span>

      {/* Image and Buttons */}
      <div className="flex items-center gap-6">
        {/* Image Placeholder */}
        <div
          onClick={handleImageClick}
          className={`w-32 h-32 flex items-center justify-center border border-gray-200 bg-gray-50 rounded-md cursor-pointer hover:border-gray-300 transition-colors ${
            image ? "relative" : ""
          }`}
          style={{
            borderStyle: "dashed",
            borderWidth: "1px",
            borderDasharray: "3 3",
          }}
        >
          {image ? (
            <>
              <img
                src={image}
                alt="uploaded"
                className="w-full h-full object-cover rounded-md"
              />
            </>
          ) : (
            <ImageSquare size={24} className="text-gray-300" />
          )}
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={handleFileChange}
            className="hidden"
          />
        </div>

        <div className="flex flex-col gap-3">
          <Button
            variant="outline"
            icon={TrashSimple}
            onClick={handleRemove}
            disabled={!image}
          >
            Remove
          </Button>
          <span className="text-xs text-gray-500">{details}</span>
        </div>
      </div>
    </div>
  );
};

export default ImageUploader;
