from flask import Flask, jsonify, request
import numpy as np
from ultralytics import YOLO
from model import DetectionModelWrapper
from io import BytesIO
from preprocess import extract_page_by_num,get_anotattions
import json
import fitz
import matplotlib
import matplotlib.pyplot as plt


app = Flask(__name__)
app.config['DEBUG'] = True
model = None

@app.route("/pred", methods=["POST", "GET"])
def image_post_request():
    file = request.files['document_data']
    print(file)
    pdf_bytes = file.read()
    #print(pdf_bytes)
    pdf_document = fitz.open(stream = pdf_bytes, filetype="pdf")
    page_count = pdf_document.page_count
    annots = []
    print('got request, starting to processs..,page_count {page_count}')
    for i in range(page_count):
        print(f'starting handling page:{i}')
        png_img,byte_img = extract_page_by_num(pdf_document,i)
        #plt.imshow(png_img)
        #plt.show()
        annots_page = get_anotattions(png_img,byte_img,model)
        annots.extend(annots_page)
    if len(annots) == 0:
        return jsonify([])
    annot_json_dict = [annot.to_json_dict() for annot in annots]
    res_json = jsonify(annot_json_dict)
    #print(annot_json_dict)
    return res_json 


if __name__ == "__main__":
    yolo_model=YOLO('./best.pt').to('cuda')
    model = yolo_model
    app.run(host="0.0.0.0", port=5000)