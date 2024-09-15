


import torch
from PIL import Image

class DetectionModelWrapper:
    def __init__(self, model, device='cpu'):
        self.model = model
        self.device = device

    def predict(self, image_path):
        return self.model(image_path)

    def load_weights(self, weights_path):
        self.model.load('./best_100_nirs.pt')






import torch
from PIL import Image

class DetectionModelWrapper:
    def __init__(self, model, device='cpu'):
        self.model = model
        self.device = device

    def predict(self, image_path):
        return self.model(image_path)

    def load_weights(self, weights_path):
        self.model.load('./best_100_nirs.pt')


