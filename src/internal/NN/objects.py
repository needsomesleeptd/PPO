class BoundingBoxYolo:
    def __init__(self, x, y, width, height):
        self.x = x
        self.y = y
        self.width = width
        self.height = height
    
    def get_top_left_corner(self):
        return (self.x, self.y)
    
    def get_bottom_right_corner(self):
        return (self.x + self.width, self.y + self.height)
    
    def get_area(self):
        return self.width * self.height
    
    def move(self, new_x, new_y):
        self.x = new_x
        self.y = new_y
        
    def resize(self, new_width, new_height):
        self.width = new_width
        self.height = new_height

