
from abc import ABC, abstractmethod
import numpy as np
import cv2
from PIL import Image

SCHEME_DETECTION_CLASS = 16 #from config.yaml


def convert_pil_to_cv2_img(image:Image):
    arr_image = np.array(image)
    image = cv2.cvtColor(arr_image, cv2.COLOR_RGB2BGR)
    return image


class ErrorDetector(ABC):
    @abstractmethod
    def detect_error(self, image: any) -> bool:
        pass
    @abstractmethod
    def get_err_class(self) -> int:
        pass
    @abstractmethod
    def get_detection_class(self) -> int:
        pass
    @abstractmethod
    def get_detected_image(self) -> int:
        pass

class SchemeErrorDetector(ErrorDetector):
    @abstractmethod
    def detect_error(self, image: any) -> bool:
        pass
    @abstractmethod
    def get_err_class(self) -> int:
        pass
    def get_detection_class(self) -> int:
        return SCHEME_DETECTION_CLASS
    @abstractmethod
    def get_detected_image(self):
        pass