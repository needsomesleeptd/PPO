#tle19u857
import math, cv2, sys
import numpy as np
from detection_scripts.schemes.scheme_err_detector import *


SCHEME_ARR_DEST_ERR_CLASS = 161

def get_filter_arrow_image(threslold_image):
    blank_image = np.zeros_like(threslold_image)

    # dilate image to remove self-intersections error
    kernel_dilate = cv2.getStructuringElement(cv2.MORPH_RECT, (2, 2))
    threslold_image = cv2.dilate(threslold_image, kernel_dilate, iterations=1)

    contours, hierarchy = cv2.findContours(threslold_image, cv2.RETR_TREE, cv2.CHAIN_APPROX_SIMPLE)

    if hierarchy is not None:

        threshold_distnace = 1000

        for cnt in contours:
            hull = cv2.convexHull(cnt, returnPoints=False)
            defects = cv2.convexityDefects(cnt, hull)

            if defects is not None:
                for i in range(defects.shape[0]):
                    start_index, end_index, farthest_index, distance = defects[i, 0]

                    # you can add more filteration based on this start, end and far point
                    # start = tuple(cnt[start_index][0])
                    # end = tuple(cnt[end_index][0])
                    # far = tuple(cnt[farthest_index][0])

                    if distance > threshold_distnace:
                        cv2.drawContours(blank_image, [cnt], -1, 255, -1)

        return blank_image
    else:
        return None


def get_length(p1, p2):
    line_length = ((p1[0] - p2[0]) ** 2 + (p1[1] - p2[1]) ** 2) ** 0.5
    return line_length


def get_max_distace_point(cnt):
    max_distance = 0
    max_points = None
    for [[x1, y1]] in cnt:
        for [[x2, y2]] in cnt:
            distance = get_length((x1, y1), (x2, y2))

            if distance > max_distance:
                max_distance = distance
                max_points = [(x1, y1), (x2, y2)]

    return max_points


def angle_beween_points(a, b):
    arrow_slope = (a[0] - b[0]) / (a[1] - b[1])
    arrow_angle = math.degrees(math.atan(arrow_slope))
    return arrow_angle


def get_arrow_info(arrow_image):
    arrow_info_image = cv2.cvtColor(arrow_image.copy(), cv2.COLOR_GRAY2BGR)
    contours, hierarchy = cv2.findContours(arrow_image, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    arrow_info = []
    if hierarchy is not None:

        for cnt in contours:
            # draw single arrow on blank image
            blank_image = np.zeros_like(arrow_image)
            cv2.drawContours(blank_image, [cnt], -1, 255, -1)

            point1, point2 = get_max_distace_point(cnt)

            angle = angle_beween_points(point1, point2)
            lenght = get_length(point1, point2)

            cv2.line(arrow_info_image, point1, point2, (0, 255, 255), 1)

            cv2.circle(arrow_info_image, point1, 2, (255, 0, 0), 3)
            cv2.circle(arrow_info_image, point2, 2, (255, 0, 0), 3)

            cv2.putText(arrow_info_image, "angle : {0:0.2f}".format(angle),
                        point2, cv2.FONT_HERSHEY_PLAIN, 0.8, (0, 0, 255), 1)
            cv2.putText(arrow_info_image, "lenght : {0:0.2f}".format(lenght),
                        (point2[0], point2[1] + 20), cv2.FONT_HERSHEY_PLAIN, 0.8, (0, 0, 255), 1)

        return arrow_info_image, arrow_info
    else:
        return None, None

def find_tip(points, convex_hull):
    length = len(points)
    indices = np.setdiff1d(range(length), convex_hull)

    for i in range(2):
        j = indices[i] + 2
        if j > length - 1:
            j = length - j
        if np.all(points[j] == points[indices[i - 1] - 2]):
            return tuple(points[j])
        
def get_args_console():
    args = sys.argv
    if len(args) != 2:
        raise('Не хватает аргументов комндной строки: название программы, название файла со схемой')
    image_file = args[1]

    return image_file

def preprocess_image(pil_image):
    '''
    Предобработка изображения (удаление элементов схемы алгоритма)
    '''
    image = convert_pil_to_cv2_img(pil_image)
    gray_image = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    gray_blur = cv2.GaussianBlur(gray_image, (3, 3), 1)
    edged = cv2.Canny(gray_blur, 10, 250)
    #cv2.imshow("thresh_image", edged)
    kernel = cv2.getStructuringElement(cv2.MORPH_RECT, (5, 5))
    closed = cv2.morphologyEx(edged, cv2.MORPH_CLOSE, kernel)
    _, thresh_image = cv2.threshold(edged, 100, 255, cv2.THRESH_BINARY_INV)

    contours, hierarchy = cv2.findContours(edged, cv2.RETR_TREE, cv2.CHAIN_APPROX_SIMPLE)
    for cnt in contours:
        peri = cv2.arcLength(cnt, True)
        approx = cv2.approxPolyDP(cnt, 0.025 * peri, True)
        hull = cv2.convexHull(approx, returnPoints=False)
        sides = len(approx)
        if sides == 4:
            x, y, w, h = cv2.boundingRect(approx)
            cv2.rectangle(image, (x, y), (x + w, y + h), (255, 255, 255), -2)

    #cv2.imshow("thresh_image", image)
    #cv2.imwrite('thresh_image.png', image)
    #cv2.imwrite('thresh_image.png', image)
    #cv2.waitKey(0)
    #cv2.destroyAllWindows()

    return image

def recognize_arrow(image):
    arrow_tips = list(); approxes = list()
    gray_image = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    contours, _ = cv2.findContours(gray_image, cv2.RETR_TREE, cv2.CHAIN_APPROX_SIMPLE)
    for cnt in contours:
        peri = cv2.arcLength(cnt, True)
        approx = cv2.approxPolyDP(cnt, 0.0275 * peri, True)
        hull = cv2.convexHull(approx, returnPoints=False)
        sides = len(hull)
        if 6 > sides > 3 and sides + 2 == len(approx):
            arrow_tip = find_tip(approx[:,0,:], hull.squeeze())
            if arrow_tip:
                arrow_tips.append(arrow_tip)
                approxes.append(approx[:,0,:])
                cv2.drawContours(image, [cnt], -1, (0, 255, 0), 3)
                cv2.circle(image, arrow_tip, 3, (0, 0, 255), cv2.FILLED)

    #cv2.imshow("thresh_image2", image)
    #cv2.imwrite('thresh_image_2.png', image)
    #cv2.waitKey(0)
    #cv2.destroyAllWindows()

    print(len(arrow_tips), len(approxes))

    return arrow_tips, approxes

def is_valid_arrow(pil_image, arrow_tips, approxes):
    image = convert_pil_to_cv2_img(pil_image)
    is_changed = False; is_valid_arrow = True
    for i, arrow_tip in enumerate(arrow_tips):
        min_index = np.argmin(approxes[i][:,1], axis=0)

        arrow_begin_point = approxes[i][min_index]
        arrow_end_point = arrow_tip

        if arrow_begin_point[1] < arrow_end_point[1] and (arrow_end_point[0] - arrow_begin_point[0]) < 5:
            cv2.drawContours(image, [approxes[i]], -1, (0, 0, 255), 2)
            is_changed = True
        elif (arrow_end_point[1] - arrow_begin_point[1]) < 5 and arrow_begin_point[0] < arrow_end_point[0]:
            cv2.drawContours(image, [approxes[i]], -1, (0, 0, 255), 2)
            is_changed = True

    #if is_changed:
    #    cv2.imshow("thresh_image3", image)
    #    cv2.imwrite('thresh_image_3.png', image)
    #    cv2.waitKey(0)
    #    cv2.destroyAllWindows()
    return is_valid_arrow




class ArrowsDestinationErrDetector(SchemeErrorDetector):
    def __init__(self):
        self.detected_image = None
  
    def detect_error(self, image: any) -> bool:
        image_work = image.copy()
        image_preprocess = preprocess_image(image_work)
        arrow_tips, approxes = recognize_arrow(image_preprocess)
        is_valid_schema = is_valid_arrow(image, arrow_tips, approxes)
        self.detected_image = image_preprocess
        return not is_valid_schema



    def get_err_class(self):
        return  SCHEME_ARR_DEST_ERR_CLASS
   
    def get_detected_image(self):
        return self.detected_image


if __name__ == "__main__":
    #image_file = get_args_console()
    image = Image.open("ex3.png")
    detector = ArrowsDestinationErrDetector()
    print(detector.detect_error(image))

#    image_input = cv2.imread(image_file)
#    image_work = image_input.copy()
#    image_preprocess = preprocess_image(image_work)
#    arrow_tips, approxes = recognize_arrow(image_preprocess)
#    is_valid_schema = is_valid_arrow(image_input, arrow_tips, approxes)

#    print(f'{image_file}: {"OK" if is_valid_schema==True else "ERROR"}') 
