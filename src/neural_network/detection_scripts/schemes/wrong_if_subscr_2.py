# vea19u756

from detection_scripts.schemes.scheme_err_detector import *
import cv2
import numpy as np
import os
#import pytesseract
import os
os.environ["TESSDATA_PREFIX"] = "."
import easyocr
from matplotlib import pyplot as plt

SCHEME_IF_SUBS_ERR_CLASS =  162

def DrawContours(image, contours, name):
    cont = image.copy()
    for contour in contours:
        cv2.drawContours(cont, [contour], 0, (0, 255, 0), 3)
    return cont
    #cv2.imwrite(name, cont)
def FindContours(image):
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    edges = cv2.Canny(gray, 50, 150)
    contours, _ = cv2.findContours(edges, cv2.RETR_TREE, cv2.CHAIN_APPROX_SIMPLE)
    return contours
def Find4AngleContours(image):
    contours = FindContours(image)
    angle4_contours = []
    for contour in contours:
        epsilon = 0.01 * cv2.arcLength(contour, True)
        approx = cv2.approxPolyDP(contour, epsilon, True)
        area = cv2.contourArea(approx)
        image_area = image.shape[0] * image.shape[1]
        if len(approx) == 4 and area > 1e-3 * image_area:
            angle4_contours.append(approx)
    return angle4_contours
            
def FindIfContours(pil_image):
    image = convert_pil_to_cv2_img(pil_image)
    angle4_contours = Find4AngleContours(image)
    parallelograms = []
    for approx in angle4_contours:
        check_90 = False
        for i in range(len(approx) - 1):
            vector1 = np.reshape(approx[i] - approx[i - 1], (2))
            vector2 = np.reshape(approx[i + 1] - approx[i], (2))
            cosine_angle = np.dot(vector1, vector2) / (np.linalg.norm(vector1) * np.linalg.norm(vector2))
            if abs(cosine_angle) < 0.1 :
                check_90 = True
                break
        if not check_90:
            parallelograms.append(approx)
    return parallelograms
def AnalyzeText(pil_image, parallelograms):
    image = convert_pil_to_cv2_img(pil_image)
    found = False
    for r in parallelograms:
        x, y, w, h = cv2.boundingRect(r)
        to_test_areas = []
        to_test_areas.append(image[y:y+h, x+w:x+2*w]) # right
        to_test_areas.append(image[y:y+h, x-w:x]) # left
        to_test_areas.append(image[y+h:y+h+int(h/2), x:x+w]) # bottom
        for roi in to_test_areas:
            if roi.any():
                gray = cv2.cvtColor(roi, cv2.COLOR_BGR2GRAY)
                
                reader = easyocr.Reader(['ru', 'en'])
                text = reader.readtext(gray, detail=0, paragraph=True)
                if 'да' in text or 'Да' in text or 'нет' in text or 'Нет' in text:
                    found = True
    if len(parallelograms) != 0 and not found:
        return True    
    return False



class IfSubscriptionErrDetector(SchemeErrorDetector):
    def __init__(self):
        self.detected_image = None
  
    def detect_error(self, image: any) -> bool:
        parallelograms = FindIfContours(image)
        img_cv = convert_pil_to_cv2_img(image)
        img_cv= DrawContours(img_cv,parallelograms,"detected")
        error_found = AnalyzeText(image, parallelograms)
        self.detected_image = img_cv
        return error_found



    def get_err_class(self):
        return  SCHEME_IF_SUBS_ERR_CLASS
   
    def get_detected_image(self):
        return self.detected_image

if __name__ == "__main__":
    #image_file = get_args_console()
    image = Image.open("scheme_trig.png")
    detector = IfSubscriptionErrDetector()
    print(detector.detect_error(image))
    plt.imshow(detector.detected_image)
    plt.show()

