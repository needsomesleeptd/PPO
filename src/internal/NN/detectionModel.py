from abc import ABC, abstractmethod

class ObjectDetectionNeuralNetwork(ABC):
    def __init__(self):
        pass

    @abstractmethod
    def load_model(self, model_path):
        pass

    @abstractmethod
    def detect_objects(self, image):
        pass

    @abstractmethod
    def draw_boxes(self, image, boxes, labels):
        pass