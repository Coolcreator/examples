import unittest
from custom_list import CustomList

class TestCase(unittest.TestCase):
    def setUp(self):
        self.list1 = CustomList([1, 2])
        self.list2 = CustomList([3, 4])
        self.list3 = CustomList([5, 6, 7])
        self.list4 = CustomList([8, 9, 10])
        self.list5 = CustomList([11, 12, 13, 14, 15])

    def test_add(self):
        self.assertEqual(self.list1 + self.list1, [2, 4])
        self.assertEqual(self.list1 + self.list2, [4, 6])
        self.assertEqual(self.list2 + self.list1, [4, 6])
        self.assertEqual(self.list2 + self.list3, [8, 10, 7])
        self.assertEqual(self.list3 + self.list2, [8, 10, 7])
        self.assertEqual(self.list4 + [], [8, 9, 10])
        self.assertEqual([] + self.list4, [8, 9, 10])
        self.assertEqual([0, 1] + self.list5, [11, 13, 13, 14, 15])
        self.assertEqual(self.list5 + [0, 1, 2], [11, 13, 15, 14, 15])
        self.assertEqual(CustomList([]) + CustomList([]), [])

    def test_sub(self):
        self.assertEqual(self.list1 - self.list1, [0, 0])
        self.assertEqual(self.list1 - self.list2, [-2, -2])
        self.assertEqual(self.list2 - self.list1, [2, 2])
        self.assertEqual(self.list2 - self.list3, [-2, -2, -7])
        self.assertEqual(self.list3 - self.list2, [2, 2, 7])
        self.assertEqual(self.list4 - [], [8, 9, 10])
        self.assertEqual([] - self.list4, [-8, -9, -10])
        self.assertEqual([0, 1] - self.list5, [-11, -11, -13, -14, -15])
        self.assertEqual(self.list5 - [0, 1, 2], [11, 11, 11, 14, 15])
        self.assertEqual(CustomList([]) - CustomList([]), [])

    def test_eq(self):
        self.assertTrue(self.list1 == self.list1)
        self.assertTrue(self.list2 == self.list1 + [2, 2])
        self.assertTrue(self.list3 == self.list4 - [3, 3, 3])
        self.assertTrue(self.list4 == [] + self.list4 == self.list4)
        self.assertTrue(self.list5 == [3, 3, 3, 14, 15] + self.list4 == self.list5)

        self.assertFalse(self.list1 == self.list2)
        self.assertFalse(self.list2 == self.list3)
        self.assertFalse(self.list3 == self.list4)
        self.assertFalse(self.list4 == self.list5)
        self.assertFalse(self.list5 == [])

    def test_lt(self):
        self.assertTrue(self.list1 < self.list2)
        self.assertTrue(self.list2 < self.list3)
        self.assertTrue(self.list3 < self.list4)
        self.assertTrue(self.list4 < self.list5)
        self.assertTrue(self.list5 < self.list5 + [0, 0, 0, 0, 1])

        self.assertFalse(self.list2 < self.list1)
        self.assertFalse(self.list3 < self.list2)
        self.assertFalse(self.list4 < self.list3)
        self.assertFalse(self.list5 < self.list4)
        self.assertFalse(self.list5 < self.list5)

    def test_gt(self):
        self.assertTrue(self.list2 > self.list1)
        self.assertTrue(self.list3 > self.list2)
        self.assertTrue(self.list4 > self.list3)
        self.assertTrue(self.list5 > self.list4)
        self.assertTrue(self.list5 + self.list1 > self.list5)

        self.assertFalse(self.list1 > self.list2)
        self.assertFalse(self.list2 > self.list3)
        self.assertFalse(self.list3 > self.list4)
        self.assertFalse(self.list4 > self.list5)
        self.assertFalse(self.list5 > self.list5 + self.list1)

    def test_le(self):
        self.assertTrue(self.list1 <= self.list2)
        self.assertTrue(self.list2 <= self.list3)
        self.assertTrue(self.list4 <= self.list5)
        self.assertTrue([1] <= [1, 1])
        self.assertTrue([] <= [])

    def test_ge(self):
        self.assertFalse(self.list1 >= self.list2)
        self.assertFalse(self.list2 >= self.list3)
        self.assertFalse(self.list3 >= self.list4)
        self.assertFalse(self.list4 >= [100])
        self.assertFalse(self.list5 >= [100])

    def test_type(self):
        self.assertTrue(type([1, 2] + CustomList([3, 4])), CustomList)
        self.assertTrue(type([3, 4] - CustomList([1, 2])), CustomList)
        self.assertTrue(type(CustomList([5, 6]) + CustomList([7, 8])), CustomList)
        self.assertTrue(type(CustomList([7, 8]) - CustomList([5, 6])), CustomList)
        self.assertTrue(type(CustomList([]) + []), CustomList)

if __name__ == '__main__':
    unittest.main()
