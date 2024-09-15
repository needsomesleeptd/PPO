#pas23um188 

from detection_scripts.formulas.formulas_err_detector import  FormulasErrorDetector,convert_pil_to_cv2_img

from pathlib import Path
import os

# Import OpenCV module 
import PIL.Image
import cv2
# Import pyplot from matplotlib as pltd 
from matplotlib import pyplot as pltd
import PIL

import random


FORMULAS_POSITION_ERR_CLASS = 151


current_dir = "./checks" # for debugging 
image_name = random.randint(1,100)

def handle_image(pil_imging : PIL.Image, debug_value: bool = False):
    cv2_image = convert_pil_to_cv2_img(pil_imging)
    # Opening the image from files 
    # Altering properties of image with cv2 
    img_gray = cv2.cvtColor(cv2_image, cv2.COLOR_BGR2GRAY)
    imaging_rgb = cv2.cvtColor(cv2_image, cv2.COLOR_BGR2RGB)
    if debug_value:
        # Plotting image with subplot() from plt 
        pltd.subplot(1, 1, 1)
        # Displaying image in the output 
        pltd.imshow(imaging_rgb)
        pltd.show()
    gray = cv2.cvtColor(cv2_image, cv2.COLOR_BGR2GRAY)
    ret, thresh1 = cv2.threshold(gray, 0, 255, cv2.THRESH_OTSU |
                                 cv2.THRESH_BINARY_INV)
    if debug_value:
        cv2.imwrite('threshold_image.jpg',thresh1)
    rect_kernel = cv2.getStructuringElement(cv2.MORPH_RECT, (12, 12))
    dilation = cv2.dilate(thresh1, rect_kernel, iterations=2)
    if debug_value:
        cv2.imwrite('dilation_image.jpg',dilation)
    contours, hierarchy = cv2.findContours(
        dilation, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_NONE)
    
    im2 = cv2_image.copy()
    # cv2.imwrite('test.jpg',cv2.rectangle(im2, (x, y), (x + w, y + h), (0, 255, 0), 2))
    centering_flag = False
    for cnt in contours:
        # Находим координаты
        x, y, w, h = cv2.boundingRect(cnt)
        y_img, x_img, _ = im2.shape
        # Правило № 1
        # Изображение по ширине находится в правой части (0.85 от ширины изображения)
        if not (x >= 85 * x_img / 100):
            if debug_value:
                print('Изображение по ширине НЕ находится в правой части (0.85 от ширины)')
            continue
        # print(x, y, w, h)
        if (h / y_img) < 0.6:
            if debug_value:
                print(h / y_img)
            # Правило № 2
            # Изображение по высоте находится примерно по середине (0.75 от ширины изображения)
            if not (y + h <= 85 * y_img / 100):
                
                if debug_value:
                    print('Изображение по высоте НЕ находится примерно по середине (0.75)')
                    rect=cv2.rectangle(im2, (x, y), (x + w, y + h), (0, 0, 255), 2)
                    cv2.imwrite(str(current_dir / 'marked_image' / image_name), rect)
                continue
            # Правило № 3
            # Изображение по высоте находится примерно по середине (0.22 от ширины изображения)
            if not (y >= 22 * y_img / 100):
                
                if debug_value:
                    print(
                        'Изображение по высоте НЕ находится примерно по середине (0.22)',
                        f'{y} >= {25 * y_img / 100}'
                    )
                    rect=cv2.rectangle(im2, (x, y), (x + w, y + h), (0, 0, 255), 2)
                    cv2.imwrite(str(current_dir / 'marked_image' / image_name), rect)
                continue
        # Правило № 4
        # Площадь прямоугольника не больше 120 * 100
        if w * h >= 120 * 100:
            if debug_value:
                print('Площадь слишком большая')
            continue
        # Правило № 5
        # Площадь прямоугольника не менее 60 * 50
        if w * h <= 60 * 50:
            if debug_value:
                print('Площадь слишком маленькая')
            continue
        # Правило № 6
        # Ширина должна быть больше высоты
        if w < h:
            if debug_value:
                print('Ширина меньше высоты')
            continue
        
        centering_flag = True
        rect=cv2.rectangle(im2, (x, y), (x + w, y + h), (0, 255, 0), 2)
        if debug_value:
            cv2.imwrite(str(current_dir / 'marked_image' / image_name), rect)
     

    return centering_flag,im2


class CheckingFormulasPositions(FormulasErrorDetector):
    def __init__(self):
        self.detected_image = None
  
    def detect_error(self, image: any) -> bool:
        res,img = handle_image(image)
        self.detected_image = img
        return not res
        
  
    def get_err_class(self) -> int:
        return FORMULAS_POSITION_ERR_CLASS

   
    def get_detected_image(self):
        return self.detected_image
