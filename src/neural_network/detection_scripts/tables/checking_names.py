import cv2
import pytesseract
import re
import os
from matplotlib import pyplot as plt
from detection_scripts.tables.table_err_detector import * 

NO_TABLE_NAME_ERR_CLASS =  171

TABLE_WRONG_NAME_ERR_CLASS =  172


def find_table(thresh, image):
    contours, _ = cv2.findContours(thresh, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    x, y, w, h = None, None, None, None

    for contour in contours:
        x, y, w, h = cv2.boundingRect(contour)
        if w > 50 and h > 50:
            cv2.rectangle(image, (x, y), (x+w, y+h), (0, 255, 0), 2)
           
    return [x, y, w, h]

def split_sentences_from_word(text, word):
    sentences = re.split(r'(?<=[.!?\n])\s+', text)
    
    start_index = -1
    for i, sentence in enumerate(sentences):
        if word in sentence:
            start_index = i
            break
    
    if start_index != -1:
        new_text = '\n'.join(sentences[start_index:])
        return new_text
    else:
        return None

def find_table_name(thresh, image):
    text = pytesseract.image_to_string(thresh, lang='rus')
    text_data =  pytesseract.image_to_data(thresh, lang='rus', output_type=pytesseract.Output.DICT)
    

    text = split_sentences_from_word(text, 'Таблица')

    if text:
        text = text.split('\n')[0].split(' ')
    else:
        return None, None

    x, y, w, h = None, None, None, None
    for i in range(len(text_data['text'])):
        if text_data['text'][i] == text[0]:
            x, y = text_data['left'][i], text_data['top'][i]
        if text_data['text'][i] == text[-1]:
            w, h = text_data['left'][i] + text_data['width'][i], text_data['top'][i] + text_data['height'][i]
            cv2.rectangle(image, (x, y), (w, h), (0, 255, 0), 2)

    return text, [x, y, w, h]

def verification(text, text_border, table_border):
    errors = 0
    try:
        if abs(text_border[0] + text_border[2] - table_border[0] - table_border[2]) > 5: #threshold
            print("Не соблюдено выраванивание")
            errors += 1
        #if abs(text_border[0] - table_border[0]) > 5:
        #    print("Не соблюдено выравнивание")
        #    errors += 1
        if table_border[1] < text_border[1]:
            print("Неверное расположение")
            errors += 1
        if not text[1][-1].isdigit():
            print("Неверная нумерация")
            errors += 1
        if '—' not in text[2]:
            print("Неверное оформление -")
            errors += 1
        if '.' == text[-1][-1]:
            print("Неверное оформление .")
            errors += 1
        if text[0][0] != 'Т':
            print("Заглавная!")
            errors += 1
    except:
        return 'Неверный формат записи'

    return errors



class TableNameErrDetector(TableErrorDetector):
    def __init__(self):
        self.detected_image = None
        self.err_class = NO_TABLE_NAME_ERR_CLASS
        os.environ["TESSDATA_PREFIX"] = r"/usr/share/tesseract-ocr/4.00/tessdata"
        pytesseract.pytesseract.tesseract_cmd = r'/usr/bin/tesseract'
    def get_err_class(self) -> int:
        return self.err_class

    def detect_error(self, image: any) -> bool:
        image_cv = convert_pil_to_cv2_img(image)

        gray = cv2.cvtColor(image_cv, cv2.COLOR_BGR2GRAY)
        thresh = cv2.threshold(gray, 0, 255, cv2.THRESH_BINARY_INV + cv2.THRESH_OTSU)[1]
        table_border = find_table(thresh, image_cv)
        text, text_border = find_table_name(thresh, image_cv)
        self.detected_image = image_cv
        if not text :
            return True
        errors = verification(text, text_border, table_border)
        if errors > 0:
            self.err_class = TABLE_WRONG_NAME_ERR_CLASS
            return True
        return False
        
   
    def get_detected_image(self):
        return self.detected_image
    




def main():
    pytesseract.pytesseract.tesseract_cmd = r'/usr/bin/tesseract'
    
    image = cv2.imread('table.png')

    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    thresh = cv2.threshold(gray, 0, 255, cv2.THRESH_BINARY_INV + cv2.THRESH_OTSU)[1]

    table_border = find_table(thresh, image)
    text, text_border = find_table_name(thresh, image)

    if not text:
        print("Подписи не найдено или с ней вообще всё очень плохо")
        return
    
    errors = verification(text, text_border, table_border)

    if not errors:
        print("ВСё гуд")
    else:
        print(f'Количество ошибок: {errors}')

    cv2.imshow('Text Table', image)
    cv2.waitKey(0)
    cv2.destroyAllWindows()

if __name__ == "__main__":
   
    image = Image.open("table.png")
    detector = TableNameErrDetector()
    print(detector.detect_error(image))
    plt.imshow(detector.detected_image)
    plt.show()

