class CustomMeta(type):
    def __new__(cls, name, bases, dct):
        custom_attrs = {}
        for key, value in dct.items():
            if key[:2] == '__' and key[-2:] == '__':
                custom_attrs[key] = value
            else:
                custom_attrs['custom_'+key] = value

        return super(CustomMeta, cls).__new__(cls, name, bases, custom_attrs)

    def __init__(cls, name, bases, dct):
        super().__init__(name, bases, dct)



class CustomClass(metaclass=CustomMeta):

    x = 50
    y = 150
    z = 250

    def __init__(self):
        self.t = 400

    def line(self):
        return 100

    def row(self):
        return 200

    def column(self):
        return 300
