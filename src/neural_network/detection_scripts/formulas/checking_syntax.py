# spyu19u638
import pytesseract
from PIL import Image
import cv2
import numpy as np
from detection_scripts.formulas.formulas_err_detector import  FormulasErrorDetector,convert_pil_to_cv2_img

FORMULAS_BOUNDS_ERR = 153
FORMULAS_SYNTAX_ERR = 154



def check_brackets(text):
   stack = []
   brackets_map = {'(': ')', '[': ']', '{': '}'}
   for symbol in text:
       if symbol in brackets_map.keys():
           stack.append(symbol)
       elif symbol in brackets_map.values():
           if not stack:
               return False
           last_open = stack.pop()
           if brackets_map[last_open] != symbol:
               return False
   return stack == []

def process_image_and_extract_text(image):
   image_np = np.array(image)
   gray_image = cv2.cvtColor(image_np, cv2.COLOR_BGR2GRAY)
   pytesseract.pytesseract.tesseract_cmd = r'/usr/bin/tesseract'
   text = pytesseract.image_to_string(gray_image, config='--psm 3')
   return text

def check_image(pil_img):
   brackets_valid = True
   syntax_valid = True
   text = process_image_and_extract_text(pil_img)
   #print("text: ", text)
   if not check_brackets(text):
       brackets_valid = False
   text_lines = text.split("\n")
   text_lines.pop()
   if len(text_lines) == 0:
       syntax_valid = False
   else: 
        last_character = max(text_lines, key=len)[-1]
        if not last_character in ('.', ',', ';'):
            syntax_valid = False
   return syntax_valid,brackets_valid,pil_img


class CheckingFormulasSyntax(FormulasErrorDetector):
    def __init__(self):
        self.detected_image = None
        self.err_class = 0 #which means no error
    def detect_error(self, image: any) -> bool:
        valid_syn,valid_br,img = check_image(image)
        self.detected_image = img

        if not valid_br:
            self.err_class = FORMULAS_BOUNDS_ERR

        if not valid_syn:
            self.err_class = FORMULAS_SYNTAX_ERR

        return (not valid_br) or (not valid_syn)
        
  
    def get_err_class(self) -> int:
        return self.err_class

   
    def get_detected_image(self):
        return self.detected_image