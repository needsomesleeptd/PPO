# mia19u558
import cv2
import matplotlib.pyplot as plt
import numpy as np
from detection_scripts.graphs.graphs_err_detector import convert_pil_to_cv2_img,GraphsErrorDetector
from PIL  import Image

GRAPH_AXIS_SUBS_ERR_CLASS = 192
BIN_THRESHHOLD = 140

def contours_to_tuple4(contours):
    return [cv2.boundingRect(cnt) for cnt in contours]


def draw_contours_on_image(image, contours):
    for cnt in contours:
        x, y, w, h = cnt
        cv2.rectangle(image, (x, y), (x + w, y + h), (0, 255, 0), 2)


def find_similar_contours(contours, threshold=2):
    similar_y = []
    similar_x = []

    for contour in contours:
        x, y, w, h = contour

        found_similar_y = False
        for group in similar_y:
            ref_y = group[0][1]
            if abs(y - ref_y) <= threshold:
                group.append((x, y, w, h))
                found_similar_y = True
                break
        if not found_similar_y:
            similar_y.append([(x, y, w, h)])

        found_similar_x = False
        for group in similar_x:
            ref_x = group[0][0] + group[0][2]
            if abs((x + w) - ref_x) <= threshold:
                group.append((x, y, w, h))
                found_similar_x = True
                break
        if not found_similar_x:
            similar_x.append([(x, y, w, h)])

    similar_y = [group for group in similar_y if len(group) >= 3]
    similar_x = [group for group in similar_x if len(group) >= 3]

    return similar_y, similar_x


def unpack_contours(contours_list):
    res = []
    for ctrs in contours_list:
        for c in ctrs:
            res.append(c)
    return res


def area_filter(min_area, input_image):
    components_number, labeled_image, component_stats, component_centroids = \
        cv2.connectedComponentsWithStats(input_image, connectivity=4)

    remaining_component_labels = [i for i in range(1, components_number) if component_stats[i][4] >= min_area]

    filtered_image = np.where(np.isin(labeled_image, remaining_component_labels) == True, 255, 0).astype('uint8')

    return filtered_image


def remove_colored_lines(image, binary_thresh, min_area=4):
    img_float = image.astype(np.float64) / 255.

    k_channel = 1 - np.max(img_float, axis=2)

    k_channel = (255 * k_channel).astype(np.uint8)
    _, binary_image = cv2.threshold(k_channel, binary_thresh, 255, cv2.THRESH_BINARY)
    binary_image = area_filter(min_area, binary_image)
    kernel_size = 3
    op_iterations = 2
    morph_kernel = cv2.getStructuringElement(cv2.MORPH_RECT, (kernel_size, kernel_size))
    binary_image = cv2.morphologyEx(binary_image, cv2.MORPH_CLOSE, morph_kernel, None, None, op_iterations,
                                    cv2.BORDER_REFLECT101)
    return binary_image


def check_uniform_spacing(coordinates, threshold=2):
    sorted_by_y = sorted(coordinates, key=lambda c: c[1])
    sorted_by_x = sorted(coordinates, key=lambda c: c[0])

    y_diff = [sorted_by_y[i + 1][1] - (sorted_by_y[i][1] + sorted_by_y[i][3]) for i in range(len(sorted_by_y) - 1)]
    x_diff = [sorted_by_x[i + 1][0] - (sorted_by_x[i][0] + sorted_by_x[i][2]) for i in range(len(sorted_by_x) - 1)]

    is_eq_y_spaced = all(abs(diff - y_diff[0]) <= threshold for diff in y_diff)
    is_eq_x_spaced = all(abs(diff - x_diff[0]) <= threshold for diff in x_diff)

    return is_eq_y_spaced, is_eq_x_spaced


def filter_contour_subarrays(contours_list, threshold=2):
    res = []
    for ctrs in contours_list:
        w_min = min(ctrs, key=lambda c: c[2])[2]
        h_min = min(ctrs, key=lambda c: c[3])[3]

        w_max = max(ctrs, key=lambda c: c[2])[2]
        h_max = max(ctrs, key=lambda c: c[3])[3]

        count = len(ctrs)
        count_same_h = sum(1 for c in ctrs if abs(c[3] - h_min) <= threshold)
        count_same_w = sum(1 for c in ctrs if abs(c[2] - w_min) <= threshold)

        is_eq_y_spaced, is_eq_x_spaced = check_uniform_spacing(ctrs, 3)

        if (count_same_h == count - 1 and h_max > h_min * 2 and not is_eq_y_spaced) \
                or (count_same_w == count - 1 and w_max > w_min * 2 and not is_eq_x_spaced):
            res.append(ctrs)

    return res


def remove_bottom_text_contours(contours, image_h):
    return list(filter(lambda c: c[1] < 4 * image_h / 5, contours))


def image_has_wrong_axis_title(pil_img, bin_threshold, show_images=True):
    img = convert_pil_to_cv2_img(pil_img)

    bin_img = remove_colored_lines(img, bin_threshold)

    rect_kernel = cv2.getStructuringElement(cv2.MORPH_RECT, (8, 8))
    dilation = cv2.dilate(bin_img, rect_kernel, iterations=1)

    contours, _ = cv2.findContours(dilation, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_NONE)
    contours = contours_to_tuple4(contours)
    contours_without_bottom = remove_bottom_text_contours(contours, img.shape[0])

    im2 = img.copy()

    s_y, s_x = find_similar_contours(contours_without_bottom, 2)
    similar = [*s_y, *s_x]

    wrong_contours = filter_contour_subarrays(similar, 3)
    res = unpack_contours(wrong_contours)

    if len(res) > 0:
        draw_contours_on_image(im2, res)
        if show_images:
            draw_contours_on_image(im2, res)
            plt.imshow(dilation)
            plt.show()
            plt.imshow(im2)
            plt.show()
        return True,im2
    return False,im2



class AxisSubsErrorDetector(GraphsErrorDetector):
    def __init__(self):
        self.detected_image = None
  
    def detect_error(self, image: any) -> bool:
        res,img = image_has_wrong_axis_title(image, BIN_THRESHHOLD, show_images=False)
        self.detected_image = img
        return res
        
  
    def get_err_class(self) -> int:
        return GRAPH_AXIS_SUBS_ERR_CLASS

   
    def get_detected_image(self):
        return self.detected_image




if __name__ == "__main__":
    image = Image.open("image.png")
    detector = AxisSubsErrorDetector()
    print(detector.detect_error(image))