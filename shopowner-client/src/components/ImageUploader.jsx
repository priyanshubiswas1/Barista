import { useState } from "react";
import { ImageSquare, UploadSimple, TrashSimple } from "phosphor-react";
import Button from "./Button"; // same folder

const ImageUploader = ({
  label = "Image",
  details = "Max file size: 5MB",
  onChange,
}) => {
  const [image, setImage] = useState(null);

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      const url = URL.createObjectURL(file);
      setImage(url);
      if (onChange) onChange(file);
    }
  };

  const handleRemove = () => {
    setImage(null);
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
          className={`w-32 h-32 flex items-center justify-center border border-gray-200 bg-gray-50 rounded-md`}
          style={{
            borderStyle: "dashed",
            borderWidth: "1px",
            borderDasharray: "3 3",
          }}
        >
          {image ? (
            <img
              src={image}
              alt="uploaded"
              className="w-full h-full object-cover rounded-md"
            />
          ) : (
            <ImageSquare size={24} className="text-gray-300" /> // gray-300
          )}
        </div>

        <div className="flex flex-col gap-3">
          <div className="flex gap-3">
            <label>
              <Button variant="outline" icon={UploadSimple}>
                Upload
                <input
                  type="file"
                  accept="image/*"
                  onChange={handleFileChange}
                  className="hidden"
                />
              </Button>
            </label>

            <Button
              variant="outline"
              icon={TrashSimple}
              onClick={handleRemove}
              disabled={!image} // disabled when no image
            >
              Remove
            </Button>
          </div>
          <span className="text-xs text-gray-500">{details}</span>
        </div>
      </div>
    </div>
  );
};

export default ImageUploader;
