import { useEffect } from "react";
import { IoClose } from "react-icons/io5";

// Registers the <model-viewer> web component.
import "@google/model-viewer";

const stop = (e) => e.stopPropagation();

export default function ARModelModal({ open, title, src, onClose }) {
  useEffect(() => {
    if (!open) return;
    const onKeyDown = (e) => {
      if (e.key === "Escape") onClose?.();
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-[100] bg-black/60 backdrop-blur-sm flex items-center justify-center p-4"
      onClick={onClose}
      role="dialog"
      aria-modal="true"
      aria-label="3D viewer"
    >
      <div
        className="w-full max-w-3xl bg-white rounded-2xl shadow-2xl overflow-hidden border border-white/20"
        onClick={stop}
      >
        <div className="flex items-center justify-between px-4 py-3 border-b border-gray-100">
          <div className="min-w-0">
            <div className="text-sm font-bold text-gray-900 truncate">
              {title || "View in AR"}
            </div>
            <div className="text-xs text-gray-500 truncate">
              Tap AR to open camera (supported devices).
            </div>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="p-2 rounded-xl hover:bg-gray-100 active:scale-95 transition"
            aria-label="Close"
          >
            <IoClose size={22} />
          </button>
        </div>

        <div className="bg-gray-50">
          <model-viewer
            src={src}
            ar
            ar-modes="webxr scene-viewer quick-look"
            camera-controls
            autoplay
            shadow-intensity="1"
            exposure="1"
            style={{ width: "100%", height: "70vh" }}
          >
            <button
              slot="ar-button"
              className="px-5 py-3 rounded-xl font-bold bg-emerald-600 text-white shadow-lg hover:bg-emerald-700 active:scale-95 transition"
              style={{ position: "absolute", bottom: 16, right: 16 }}
            >
              View in AR
            </button>
          </model-viewer>
        </div>
      </div>
    </div>
  );
}

