import random
import cv2
import numpy as np
# broken -- no -- rt



def split_and_scale_image(image, num_parts, scale_factor):

    output_image = cv2.resize(image, None, fx=scale_factor, fy=scale_factor)

    return output_image



def preprocess(img):
    img_gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    img_blur = cv2.GaussianBlur(img_gray, (5, 5), 1)
    img_canny = cv2.Canny(img_blur, 50, 50)
    kernel = np.ones((3, 3))
    img_dilate = cv2.dilate(img_canny, kernel, iterations=2)
    img_erode = cv2.erode(img_dilate, kernel, iterations=1)
    return img_erode

def find_tip(points, convex_hull):
    print(points)
    print(convex_hull)
    length = len(points)
    indices = np.setdiff1d(range(length), convex_hull)
    for i in range(2):
        j = indices[i] + 2
        if j > length - 1:
            j = length - j
        print(points[j], points[indices[i - 1] - 2])
        if np.all(points[j] == points[indices[i - 1] - 2]):
            return tuple(points[j])





num_parts = 8  
scale_factor = 0.1  
img = split_and_scale_image(image, num_parts, scale_factor)

contours, hierarchy = cv2.findContours(preprocess(img), cv2.RETR_LIST, cv2.CHAIN_APPROX_TC89_KCOS)
flag = 0
for cnt in contours:
    peri = cv2.arcLength(cnt, True)
    approx = cv2.approxPolyDP(cnt, 0.01 * peri, True)
    hull = cv2.convexHull(approx, returnPoints=False)
    sides = len(hull)
    if 11 > sides > 3 and sides + 2 == len(approx):
        arrow_tip = find_tip(approx[:,0,:], hull.squeeze())
        if arrow_tip:
            cv2.drawContours(img, [cnt], -1, (0, 255, 0), 3)
            
            # cv2.drawContours(img, [cnt], -1, (0, 255, 0), 3)
            # for p in approx[:,0,:]:
            #     cv2.drawMarker(img, p,(255, 0, 0) )
            cv2.circle(img, arrow_tip, 3, (0, 0, 255), cv2.FILLED)
            flag = 1
            print("Arrow found")
    #else:
    #    cv2.drawContours(img, [cnt], -1, (random.randint(0, 255), random.randint(0, 255), random.randint(0, 255)), 3)

if not flag:
    print("arrows not found")
num_parts = 1  
scale_factor = 0.1  
img = split_and_scale_image(img, num_parts, scale_factor)

cv2.imshow("Image", img)
cv2.waitKey(0)
