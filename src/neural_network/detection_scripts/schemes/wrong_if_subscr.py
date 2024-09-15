# tnr19u668
# broken --  Unrecognized keyword arguments passed to Dense: {'weights': [array([[0., 0., 0., 0., 0., 0.], something with weights


import cv2 as cv
import numpy as np
import os
import keras_ocr as ko
from math import sqrt

def is_rhomb(figure):
    if len(figure) != 4:
        return False
    
    lens = []
    for i in range(4):
        lens.append(sqrt((figure[i][0][0] - figure[(i + 1) % 4][0][0]) ** 2 + (figure[i][0][1] - figure[(i + 1) % 4][0][1]) ** 2))
    
    avg_len = sum(lens) / 4
    diffs = [abs(length - avg_len) for length in lens]
    
    if max(diffs) > 10:
        return False
    
    if abs(figure[1][0][1] - figure[0][0][1]) < 5 or abs(figure[1][0][0] - figure[0][0][0]) < 5:
        return False
    
    if abs(figure[2][0][1] - figure[0][0][1]) > 5 and abs(figure[2][0][0] - figure[0][0][0]) > 5:
        return False
    
    return True

image = cv.imread("scheme_wrong.png")
orig_image = image.copy()

gray = cv.cvtColor(image, cv.COLOR_BGR2GRAY)
_, threshold = cv.threshold(gray, 215, 255, cv.THRESH_BINARY) 
contours, _ = cv.findContours(threshold, cv.RETR_TREE, cv.CHAIN_APPROX_SIMPLE)

good = True
    
for contour in contours:
    x, y, w, h = cv.boundingRect(contour)
    if (w <= 50 or w >= 300) and (h <= 50 or h >= 200): #or len(contour) < 500:
         continue
    
    approx = cv.approxPolyDP( 
        contour, 0.01 * cv.arcLength(contour, True), True) 
    
    if is_rhomb(approx):
        cv.drawContours(image, [contour], 0, (255,0,0), thickness=3)
        x_s = [approx[i][0][0] for i in range(4)]
        y_s = [approx[i][0][1] for i in range(4)]
        cv.rectangle(image, (min(x_s) - 150, min(y_s) - 50), (max(x_s) + 150, max(y_s) + 50), (0, 255, 0), 2)

        fragment = [orig_image[i][max(min(x_s) - 150, 0):max(x_s) + 150] for i in range(max(0, min(y_s)-40),max(y_s)+40)]
        
        images = [ko.tools.read(np.array(fragment))]
        pipeline = ko.pipeline.Pipeline()
        prediction_groups = pipeline.recognize(images)
        arr_recognized = [prediction_groups[0][i][0] for i in range(len(prediction_groups[0]))]
        this_good = False
        
        for word in arr_recognized:
            if word in ['Da', 'da', 'pa', 'ga', 'ca', 'la', 'aa']:
                this_good = True
        
        good = good and this_good
    
cv.imwrite('Contour.jpg', image)

if not good:
    print("В схеме есть ошибки")
else:
    print("В схеме нет ошибок.")
