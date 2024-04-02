import fitz
import io
import os
from PIL import Image



def extract_images_from_pdf(pdf_bytes : str, output_folder :str, prefix : str, dpi=180):  # make dpi lower if low
    pdf_document = fitz.open(stream=pdf_bytes, filetype="pdf")

    for page_number in range(pdf_document.page_count):
        if page_number != 0:
            page = pdf_document[page_number]
            image = page.get_pixmap(matrix=fitz.Matrix(dpi / 72, dpi / 72), alpha=False)
            pil_image = Image.frombytes("RGB", [image.width, image.height], image.samples)
            pil_image.save(f"{output_folder}/{prefix}_{page_number + 1}.png", dpi=(dpi, dpi))

    pdf_document.close()
