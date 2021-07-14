import unittest
from custom_meta import CustomClass

class MetaTestCase(unittest.TestCase):

    def setUp(self):
        self.inst = CustomClass()

    def test_custom_x(self):
        self.assertEqual(self.inst.custom_x, 50)

    def test_custom_t(self):
        self.assertEqual(self.inst.t, 400)

    def test_custom_y(self):
        self.assertEqual(self.inst.custom_y, 150)

    def test_custom_z(self):
        self.assertEqual(self.inst.custom_z, 250)

    def test_custom_line(self):
        self.assertEqual(self.inst.custom_line(), 100)

    def test_custom_row(self):
        self.assertEqual(self.inst.custom_row(), 200)

    def test_custom_column(self):
        self.assertEqual(self.inst.custom_column(), 300)

    def test_x(self):
        with self.assertRaises(AttributeError):
            self.inst.x

    def test_y(self):
        with self.assertRaises(AttributeError):
            self.inst.y

    def test_z(self):
        with self.assertRaises(AttributeError):
            self.inst.z

    def test_line(self):
        with self.assertRaises(AttributeError):
            self.inst.line()

    def test_row(self):
        with self.assertRaises(AttributeError):
            self.inst.row()

    def test_column(self):
        with self.assertRaises(AttributeError):
            self.inst.column()


if __name__ == '__main__':
    unittest.main()
