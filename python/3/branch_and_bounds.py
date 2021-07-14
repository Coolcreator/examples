from collections import deque

# приоритетная очередь
class SimpleQueue():
    # инициализация очереди
    def __init__(self):
        self.buffer = deque()
    # вставка элемента в очередь слева
    def push(self, value):
        self.buffer.appendleft(value)
    # удаление элемента из очереди справа
    def pop(self):
        return self.buffer.pop()
    # длина очереди
    def __len__(self):
        return len(self.buffer)

# узел дерева
class Node():
    # инициализация узла дерева
    def __init__(self, level, selected_items, cost, weight, bound):
        self.level = level
        self.selected_items = selected_items
        self.cost = cost
        self.weight = weight
        self.bound = bound


def branch_and_bounds(number, capacity, weight_cost):
    """
    Метод ветвей и границ для решения задачи о рюкзаке
    Источник дополнительной литературы и псевдокод алгоритма
    http://faculty.cns.uni.edu/~east/teaching/153/branch_bound/knapsack/overview_algorithm.html

    параметр number: общее количество предметов
    параметр capacity: вместимость рюкзака
    параметр weight_cost: массив кортежей вида [(вес, стоимость), ...]
    возвращаемое значение метода: кортеж вида (лучшая цена, [список отобранных
    предметов (1 - предмет находится в рюкзаке, 0 - иначе)])
    """
    priority_queue = SimpleQueue()

    # сортировка элементов в невозрастающем порядке значений (цена / вес)
    ratios = [(index, item[1] / float(item[0])) for index, item in enumerate(weight_cost)]
    ratios = sorted(ratios, key=lambda x: x[1], reverse=True)

    # лучший узел
    best_so_far = Node(0, [], 0.0, 0.0, 0.0)
    # инициализация узла и добавление ее в очередь
    a_node = Node(
        0,
        [],
        0.0,
        0.0,
        calculate_bound(best_so_far, number, capacity, weight_cost, ratios))
    priority_queue.push(a_node)

    # пока приоритетная очередь не пуста
    while len(priority_queue) > 0:
        # удаление из очереди крайнего правого элемента
        curr_node = priority_queue.pop()
        # если мы еще можем достичь лучшей стоиомости
        if curr_node.bound > best_so_far.cost:
            # создаем узел next_added = сurr_node + очередной узел
            # индекс нового элемента берется из отсортированного массива ratios
            curr_node_index = ratios[curr_node.level][0]
            next_item_cost = weight_cost[curr_node_index][1]
            next_item_weight = weight_cost[curr_node_index][0]
            next_added = Node(
                curr_node.level + 1,
                curr_node.selected_items + [curr_node_index],
                curr_node.cost + next_item_cost,
                curr_node.weight + next_item_weight,
                curr_node.bound
            )
            # если вес next_added меньше вместимости рюкзака
            if next_added.weight <= capacity:
                # если стоимость очередного узла > стоимости лучшего узла
                if next_added.cost > best_so_far.cost:
                    # переопределяем лучший узел
                    best_so_far = next_added
                # если еще можем достичь на nex_added лучшей стоимости
                if next_added.bound > best_so_far.cost:
                    # добавляем ее в очередь
                    priority_queue.push(next_added)
            # создаем узел nex_not_added = curr_node без добавления нового узла
            next_not_added = Node(
                curr_node.level + 1,
                curr_node.selected_items,
                curr_node.cost,
                curr_node.weight,
                curr_node.bound)
            # вычисляем границу на узле
            next_not_added.bound = calculate_bound(
                next_not_added,
                number, capacity,
                weight_cost,
                ratios)
            # если можем достичь на next_not_added лучшей стоимости
            if next_not_added.bound > best_so_far.cost:
                # добавляем ее в очередь
                priority_queue.push(next_not_added)

    # взвращаем лучшую стоимость и результирующий массив из нулей и едениц
    best_combination = [0] * number
    for w_c in best_so_far.selected_items:
        best_combination[w_c] = 1
    return int(best_so_far.cost), best_combination


# вычисление границы
def calculate_bound(node, number, capacity, weight_cost, ratios):
    # если уже превышена вместимость рюкзака
    if node.weight >= capacity:
        return 0
    # верхняя граница, вес и уровень лучшего узла
    upper_bound = node.cost
    total_weight = node.weight
    current_level = node.level
    # пока текущий уровень (очередной предмет) меньше общего числа предметов
    while current_level < number:
        # индекс очередного предмета из отсортированного списка
        current_index = ratios[current_level][0]
        # если общий вес + очередной дополнительный вес предмета
        # из отсортированного списка > вместимости рюкзака
        if total_weight + weight_cost[current_index][0] > capacity:
            # обновление цены, веса и верхней границы => выход из функции
            cost = weight_cost[current_index][1]
            weight = weight_cost[current_index][0]
            upper_bound += (capacity - total_weight) * cost/weight
            break
        # обновление верхней границы и общего веса предметов
        upper_bound += weight_cost[current_index][1]
        total_weight += weight_cost[current_index][0]
        # следующий уровень
        current_level += 1
    # верхняя граница
    return upper_bound
