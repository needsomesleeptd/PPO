import base64
import fitz
from io import BytesIO
import os
from PIL import Image
from base64 import b64encode

path_to_pdfs = "../reports/"
path_to_images = "../images/"
output_format = "png"
    #golang format
    #ID         uint64    `json:"id"`
	#PageData   []byte    `json:"page_data"`
	#ErrorBB    []float32 `json:"error_bb"`
	#ClassLabel uint64    `json:"class_label"`

class Anotattion:
    def __init__(self,page,bbs,cls) -> None:
        self.page = page
        self.bbs = bbs
        self.cls = cls
    def to_json_dict(self):
        return {
            "page_data" : self.page,
            "error_bb" : self.bbs,
            "class_label" : self.cls
        }


def get_anotattions(png_page,byte_page, model):
    predicts = model(png_page)
    #print(predicts[0])
    annots = []
    for predict in predicts:
        print(predict.boxes.xyxy.tolist(),predict.boxes.cls.tolist())
        bboxes = predict.boxes.xyxy.tolist()
        clses = predict.boxes.cls.tolist()
        if (len(bboxes) !=0):
            for i in range(len(bboxes)):
                annot = Anotattion(b64encode(byte_page).decode('utf-8'),bboxes[i],int(clses[i]))
                annots.append(annot)
    return annots
    

def extract_images_from_pdf(pdf_path, output_folder, prefix, dpi=180):  # make dpi lower if low
    pdf_document = fitz.open(pdf_path)

    for page_number in range(pdf_document.page_count):
        if page_number != 0:
            page = pdf_document[page_number]
            image = page.get_pixmap(matrix=fitz.Matrix(dpi / 72, dpi / 72), alpha=False)
            pil_image = Image.frombytes("RGB", [image.width, image.height], image.samples)
            pil_image.save(f"{output_folder}/{prefix}_{page_number + 1}.png", dpi=(dpi, dpi))

    pdf_document.close()




    
    


def extract_page_by_num(pdf_document,page_number,dpi=180):
    page = pdf_document[page_number]
    image = page.get_pixmap(matrix=fitz.Matrix(dpi / 72, dpi / 72), alpha=False)
    pil_image = Image.frombytes("RGB", [image.width, image.height], image.samples)
    image_data = image.tobytes()
    return pil_image,image_data
