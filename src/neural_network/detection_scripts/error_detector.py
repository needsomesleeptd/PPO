from abc import ABC, abstractmethod

NO_ERR_ERR_CLASS = 0


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