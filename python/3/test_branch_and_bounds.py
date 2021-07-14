import unittest
from branch_and_bounds import branch_and_bounds

class TestMethod(unittest.TestCase):
    '''
    Запускать тесты с помощью команды python3 test_branch_and_bounds.py
    Тесты должны находиться в одной директории с файлом branch_and_bounds.py
    '''
    def setUp(self):
        self.case_1 = [2, 1, [(3, 3), (4, 4)]]
        self.case_2 = [3, 3, [(1, 1), (2, 3), (3, 2)]]
        self.case_3 = [3, 5, [(2, 1), (3, 5), (4, 1)]]
        self.case_4 = [4, 3, [(1, 1), (2, 2), (3, 3), (4, 1), (5, 3)]]
        self.case_5 = [5, 2, [(2, 2), (3, 5), (4, 1), (1, 1), (3, 3)]]

    def test_branch_and_bounds(self):
        self.assertEqual((0, [0, 0]), branch_and_bounds(*self.case_1))
        self.assertEqual((4, [1, 1, 0]), branch_and_bounds(*self.case_2))
        self.assertEqual((6, [1, 1, 0]), branch_and_bounds(*self.case_3))
        self.assertEqual((3, [1, 1, 0, 0]), branch_and_bounds(*self.case_4))
        self.assertEqual((2, [1, 0, 0, 0, 0]), branch_and_bounds(*self.case_5))

if __name__ == '__main__':
    unittest.main()
