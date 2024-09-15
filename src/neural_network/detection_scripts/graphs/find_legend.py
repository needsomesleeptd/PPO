# mmd19u555

from detection_scripts.graphs.graphs_err_detector import  GraphsErrorDetector
from PIL import Image
from imutils.object_detection import non_max_suppression
import numpy as np
import cv2
import os

_EAST_MODEL_PATH = '/home/andrew/uni/PPO/src/neural_network/detection_scripts/graphs/model/frozen_east_text_detection.pb'
_LINE_MAX_OFFSET = 80
_LINE_CENTER_OFFSET_PART = 3
_HORIZONTAL_KERNEL_SIZE_PART = 40

GRAPH_LEGEND_ERR_CLASS = 191

def _ShowWaitDestroy(winname, img):
    cv2.imshow(winname, img)
    cv2.moveWindow(winname, 500, 0)
    cv2.waitKey(0)
    cv2.destroyWindow(winname)

def _GetTextRegions(pil_image, minProb=0.7, debug=False):
    arr_image = np.array(pil_image)
    image = cv2.cvtColor(arr_image, cv2.COLOR_RGB2BGR)

    orig = np.copy(image)
    (H, W) = image.shape[:2]

    # set the new width and height and then determine the ratio in change
    # for both the width and height
    (newW, newH) = (320, 320)
    rW = W / float(newW)
    rH = H / float(newH)

    # resize the image and grab the new image dimensions
    image = cv2.resize(image, (newW, newH))
    (H, W) = image.shape[:2]

    # define the two output layer names for the EAST detector model that
    # we are interested -- the first is the output probabilities and the
    # second can be used to derive the bounding box coordinates of text
    layerNames = [
        "feature_fusion/Conv_7/Sigmoid",
        "feature_fusion/concat_3"]

    # load the pre-trained EAST text detector
    net = cv2.dnn.readNet(_EAST_MODEL_PATH)

    # construct a blob from the image and then perform a forward pass of
    # the model to obtain the two output layer sets
    blob = cv2.dnn.blobFromImage(image, 1.0, (W, H),(123.68, 116.78, 103.94), swapRB=True, crop=False)
    net.setInput(blob)
    (scores, geometry) = net.forward(layerNames)

    # grab the number of rows and columns from the scores volume, then
    # initialize our set of bounding box rectangles and corresponding
    # confidence scores
    (numRows, numCols) = scores.shape[2:4]
    rects = []
    confidences = []

    # loop over the number of rows
    for y in range(0, numRows):
        # extract the scores (probabilities), followed by the geometrical
        # data used to derive potential bounding box coordinates that
        # surround text
        scoresData = scores[0, 0, y]
        xData0 = geometry[0, 0, y]
        xData1 = geometry[0, 1, y]
        xData2 = geometry[0, 2, y]
        xData3 = geometry[0, 3, y]
        anglesData = geometry[0, 4, y]

        for x in range(0, numCols):
            # ignore probability values below 0.75
            if scoresData[x] < minProb:
                continue
            
            # compute the offset factor as our resulting feature maps will
            # be 4x smaller than the input image
            (offsetX, offsetY) = (x * 4.0, y * 4.0)

            # extract the rotation angle for the prediction and then
            # compute the sin and cosine
            angle = anglesData[x]
            cos = np.cos(angle)
            sin = np.sin(angle)

            # use the geometry volume to derive the width and height of
            # the bounding box
            h = xData0[x] + xData2[x]
            w = xData1[x] + xData3[x]

            # compute both the starting and ending (x, y)-coordinates for
            # the text prediction bounding box
            endX = int(offsetX + (cos * xData1[x]) + (sin * xData2[x]))
            endY = int(offsetY - (sin * xData1[x]) + (cos * xData2[x]))
            startX = int(endX - w)
            startY = int(endY - h)

            # add the bounding box coordinates and probability score to
            # our respective lists
            rects.append((startX, startY, endX, endY))
            confidences.append(scoresData[x])

    # apply non-maxima suppression to suppress weak, overlapping bounding
    # boxes
    boxes = non_max_suppression(np.array(rects), probs=confidences)

    # loop over the bounding boxes
    for box in boxes:
        # scale the bounding box coordinates based on the respective ratios
        box[0] = int(box[0] * rW)
        box[1] = int(box[1] * rH)
        box[2] = int(box[2] * rW)
        box[3] = int(box[3] * rH)

    if debug:
        for (startX, startY, endX, endY) in boxes:
            # draw the bounding box on the image
            cv2.rectangle(orig, (startX, startY), (endX, endY), (0, 255, 0), 2)
        _ShowWaitDestroy('boxes', orig)

    return boxes

def _GetHorizontals(pil_image, debug=False):
    arr_image = np.array(pil_image)
    img = cv2.cvtColor(arr_image, cv2.COLOR_RGB2BGR)
    

    # Convert the image to grayscale
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

    # Apply adaptiveThreshold at the bitwise_not of gray, notice the ~ symbol
    gray = cv2.bitwise_not(gray)
    bw = cv2.adaptiveThreshold(gray, 255, cv2.ADAPTIVE_THRESH_MEAN_C, \
                                cv2.THRESH_BINARY, 15, -2)

    horizontal = np.copy(bw)
    # Specify size on horizontal axis
    cols = horizontal.shape[1]
    horizontal_size = cols // _HORIZONTAL_KERNEL_SIZE_PART
    # Create structure element for extracting horizontal lines through morphology operations
    horizontalStructure = cv2.getStructuringElement(cv2.MORPH_RECT, (horizontal_size, 1))
    # Apply morphology operations
    horizontal = cv2.erode(horizontal, horizontalStructure)
    horizontal = cv2.dilate(horizontal, horizontalStructure)

    if debug:
        _ShowWaitDestroy('horizontal', horizontal)

    return horizontal

def HasLegend(image, debug=False):
    textBoxes = _GetTextRegions(image, debug=debug)
    horizontals = _GetHorizontals(image, debug=debug)

    possibleHorLineRegion = None
    for (startX, startY, endX, endY) in textBoxes:
        startLineX = startX - _LINE_MAX_OFFSET if startX - _LINE_MAX_OFFSET > 0 else 0
        centerOffset = (endY - startY) // _LINE_CENTER_OFFSET_PART
        possibleHorLineRegion = horizontals[startY+centerOffset:endY-centerOffset, startLineX:startX]

        if np.sum(possibleHorLineRegion) > 0:
            if debug:
                _ShowWaitDestroy('test', possibleHorLineRegion)
            return True,possibleHorLineRegion
    
    return False,possibleHorLineRegion


class LegendErrorDetector(GraphsErrorDetector):
    def __init__(self):
        self.detected_image = None
  
    def detect_error(self, image: any) -> bool:
        res,img = HasLegend(image)
        self.detected_image = img
        return not res
        
  
    def get_err_class(self) -> int:
        return GRAPH_LEGEND_ERR_CLASS

   
    def get_detected_image(self):
        return self.detected_image

if __name__ == "__main__":
    image = Image.open("image.png")
    detector = LegendErrorDetector()
    print(detector.detect_error(image))
    