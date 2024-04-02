from detectionModel import ObjectDetectionNeuralNetwork
import cv2
import numpy as np
import matplotlib.pyplot as plt

class YOLOObjectDetection(ObjectDetectionNeuralNetwork):
    def __init__(self,model, confidence_threshold=0.5, nms_threshold=0.4):
        self.confidence_threshold = confidence_threshold
        self.nms_threshold = nms_threshold
        self.model = model

    
    def load_model(self, model_path):
        self.model = cv2.dnn.readNetFromDarknet(model_path + '.cfg', model_path + '.weights')

    def detect_objects(self, image):
        blob = cv2.dnn.blobFromImage(image, 1/255, (416, 416), swapRB=True, crop=False)
        self.model.setInput(blob)
        outputs = self.model.forward(self.get_output_layers(self.model))
        return self.post_process(image, outputs)

    def draw_boxes(self, image, boxes, labels):
        for box in boxes:
            x, y, w, h = box
            cv2.rectangle(image, (x, y), (x + w, y + h), (255, 0, 0), 2)
            cv2.putText(image, labels[box], (x, y - 5), cv2.FONT_HERSHEY_SIMPLEX, 0.5, (255, 0, 0), 2)
        return image

    def get_output_layers(self, net):
        layer_names = net.getLayerNames()
        output_layers = [layer_names[i[0] - 1] for i in net.getUnconnectedOutLayers()]
        return output_layers

    def post_process(self, image, outputs):
        frame_height, frame_width = image.shape[:2]
        boxes = []
        confidences = []
        class_ids = []

        for output in outputs:
            for detection in output:
                scores = detection[5:]
                class_id = np.argmax(scores)
                confidence = scores[class_id]
                if confidence > self.confidence_threshold:
                    center_x, center_y, width, height = list(map(int, detection[0:4] * [frame_width, frame_height, frame_width, frame_height]))
                    x = center_x - width // 2
                    y = center_y - height // 2
                    boxes.append([x, y, width, height])
                    confidences.append(float(confidence))
                    class_ids.append(class_id)

        indices = cv2.dnn.NMSBoxes(boxes, confidences, self.confidence_threshold, self.nms_threshold)
        return [boxes[i[0]] for i in indices], [str(class_ids[i]) for i in indices]