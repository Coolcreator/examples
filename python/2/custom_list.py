class CustomList(list):

    def __add__(self, other):
        new_list = [0] * max(len(self), len(other))
        for i in range(len(self)):
            new_list[i] += self[i]
        for i in range(len(other)):
            new_list[i] += other[i]
        return CustomList(new_list)

    __radd__ = __add__

    def __neg__(self):
        return CustomList(-i for i in self)

    def __sub__(self, other):
        return -(-self + other)

    def __rsub__(self, other):
        return -self + other

    def __eq__(self, other):
        return sum(self) == sum(other)

    def __ne__(self, other):
        return sum(self) != sum(other)

    def __lt__(self, other):
        return sum(self) < sum(other)

    def __gt__(self, other):
        return sum(self) > sum(other)
