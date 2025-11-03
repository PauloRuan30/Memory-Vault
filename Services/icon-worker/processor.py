from PIL import Image, ImageDraw, ImageFilter
import random

def generate_texture(input_path, output_path):
    """
    Generate a PS2-style texture for the given file
    """
    # Open the original image
    try:
        original = Image.open(input_path)
        # Resize to a small size to simulate PS2 memory card texture
        texture_size = (64, 64)  # PS2-style small texture
        original.thumbnail(texture_size, Image.Resampling.LANCZOS)
        
        # Create a new image with PS2-style dimensions
        texture = Image.new('RGB', texture_size, (30, 30, 50))  # Dark blue background like PS2 menu
        
        # Paste the original image in the center
        x = (texture_size[0] - original.width) // 2
        y = (texture_size[1] - original.height) // 2
        texture.paste(original, (x, y))
        
    except Exception:
        # If original image can't be processed, create a generic texture
        texture = Image.new('RGB', (64, 64), (50, 50, 80))
        draw = ImageDraw.Draw(texture)
        
        # Add some PS2-style patterns
        for _ in range(10):
            x1 = random.randint(0, 63)
            y1 = random.randint(0, 63)
            x2 = random.randint(0, 63)
            y2 = random.randint(0, 63)
            color = (random.randint(100, 200), random.randint(100, 200), random.randint(150, 220))
            draw.line([x1, y1, x2, y2], fill=color, width=1)
    
    # Apply a slight blur to simulate PS2 graphics
    texture = texture.filter(ImageFilter.GaussianBlur(radius=0.5))
    
    # Add a border to simulate PS2 selection frame
    draw = ImageDraw.Draw(texture)
    draw.rectangle([0, 0, 63, 63], outline=(100, 150, 255), width=1)
    
    texture.save(output_path, 'PNG')