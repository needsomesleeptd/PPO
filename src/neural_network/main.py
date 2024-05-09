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
<<<<<<< HEAD
app.config['DEBUG'] = True
=======

>>>>>>> d3ec582 (It is ALIVE)
model = None

@app.route("/pred", methods=["POST", "GET"])
def image_post_request():
    file = request.files['document_data']
<<<<<<< HEAD
    print(file)
=======
    #print(file)
>>>>>>> d3ec582 (It is ALIVE)
    pdf_bytes = file.read()
    #print(pdf_bytes)
    pdf_document = fitz.open(stream = pdf_bytes, filetype="pdf")
    page_count = pdf_document.page_count
    annots = []
<<<<<<< HEAD
    print('got request, starting to processs..,page_count {page_count}')
=======
    print('got request, starting to processs..')
>>>>>>> d3ec582 (It is ALIVE)
    for i in range(page_count):
        print(f'starting handling page:{i}')
        png_img,byte_img = extract_page_by_num(pdf_document,i)
        #plt.imshow(png_img)
        #plt.show()
        annots_page = get_anotattions(png_img,byte_img,model)
        annots.extend(annots_page)
    if len(annots) == 0:
<<<<<<< HEAD
        return jsonify([])
    annot_json_dict = [annot.to_json_dict() for annot in annots]
    res_json = jsonify(annot_json_dict)
    #print(annot_json_dict)
=======
        return jsonify({})
    annot_json_dict = [annot.to_json_dict() for annot in annots]
    res_json = jsonify(annot_json_dict)
>>>>>>> d3ec582 (It is ALIVE)
    return res_json 


if __name__ == "__main__":
<<<<<<< HEAD
    yolo_model=YOLO('./best.pt').to('cuda')
=======
    yolo_model=YOLO('./best.pt')
>>>>>>> d3ec582 (It is ALIVE)
    model = yolo_model
    app.run(host="0.0.0.0", port=5000)