import base64
import fitz
from io import BytesIO
import io
import os
from PIL import Image
from base64 import b64encode
import cv2
from detection_scripts.error_detector import ErrorDetector,NO_ERR_ERR_CLASS

path_to_pdfs = "../reports/"
path_to_images = "../images/"
output_format = "png"
    #golang format
    #ID         uint64    `json:"id"`
	#PageData   []byte    `json:"page_data"`
	#ErrorBB    []float32 `json:"error_bb"`
	#ClassLabel uint64    `json:"class_label"`

class Anotattion:
    def __init__(self,page,bbs,cls,err) -> None:
        self.page = page
        self.bbs = bbs
        self.cls = cls
        self.err =err
    def to_json_dict(self) -> dict[str, any]:
        return {
            "page_data" : self.page,
            "error_bb" : self.bbs,
            "class_label" : self.cls,
            "type_label" : self.err,
            "was_checked" : False
        }


def get_anotattions(png_page:Image,byte_page:bytes, model,detectors : list[ErrorDetector]) -> list[Anotattion]:
    predicts = model(png_page)
    #print(predicts[0])
    annots = []
    for predict in predicts:
        print(predict.boxes.xyxy.tolist(),predict.boxes.cls.tolist())
        bboxes = predict.boxes.xyxyn.tolist()
        clses = predict.boxes.cls.tolist()
      
        for i in range(len(bboxes)):
            for detector in detectors:
                error_cls = NO_ERR_ERR_CLASS
                cls = int(clses[i])
                png_page = png_page.crop(predict.boxes.xyxy.tolist()[i])
                if detector.get_detection_class() == cls:
                    has_err = detector.detect_error(png_page)
                    if has_err:
                        print(f"detected error for class {cls},{error_cls}")
                        error_cls = detector.get_err_class()
                        # cv2_img = detector.get_detected_image()
                        # if isinstance(cv2_img, Image.Image):
                        #     #byte_page = cv2_img.to_bytes()
                        #     pass
                        # else:
                        #     pil_image = Image.fromarray(cv2.cvtColor(cv2_img, cv2.COLOR_BGR2RGB))
                        #     #byte_page = pil_image.tobytes()


            
                        annot = Anotattion(b64encode(byte_page).decode('utf-8'),bboxes[i],error_cls,int(clses[i]))
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




    
    


def extract_page_by_num(pdf_document,page_number,dpi=180)  -> (Image, bytes):
    page = pdf_document[page_number]
    image = page.get_pixmap(matrix=fitz.Matrix(dpi / 72, dpi / 72), alpha=False)
    pil_image = Image.frombytes("RGB", [image.width, image.height], image.samples)
    image_data = image.tobytes()
    return pil_image,image_data
