import argparse
import json
import os
from gradio_client import Client, handle_file


def main() -> int:
    parser = argparse.ArgumentParser(description="Run Hunyuan3D and print JSON result.")
    parser.add_argument("--image", required=True, help="Path to input image file")
    parser.add_argument("--space", default=os.environ.get("HUNYUAN3D_SPACE", "Jbowyer/Hunyuan3D-2.1"))
    parser.add_argument("--steps", type=int, default=int(os.environ.get("HUNYUAN3D_STEPS", "30")))
    parser.add_argument("--guidance-scale", type=float, default=float(os.environ.get("HUNYUAN3D_GUIDANCE_SCALE", "5")))
    parser.add_argument("--octree-resolution", type=int, default=int(os.environ.get("HUNYUAN3D_OCTREE_RESOLUTION", "256")))
    parser.add_argument("--seed", type=int, default=int(os.environ.get("HUNYUAN3D_SEED", "1234")))
    parser.add_argument("--remove-bg", action="store_true", default=os.environ.get("HUNYUAN3D_REMBG", "true").lower() == "true")
    args = parser.parse_args()

    client = Client(args.space)
    result = client.predict(
        image=handle_file(args.image),
        mv_image_front=None,
        mv_image_back=None,
        mv_image_left=None,
        mv_image_right=None,
        steps=args.steps,
        guidance_scale=args.guidance_scale,
        seed=args.seed,
        octree_resolution=args.octree_resolution,
        check_box_rembg=args.remove_bg,
        num_chunks=8000,
        randomize_seed=True,
        api_name="/generation_all",
    )

    # The space returns a list-like response that typically contains file paths/urls.
    print(json.dumps({"result": result}))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())