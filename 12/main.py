import numpy as np


def parse_shape(shape: str) -> np.ndarray:
    lines = shape.strip().split("\n")[1:]
    return np.array([[c == "#" for c in line] for line in lines])


def parse_map(map: str) -> tuple[np.ndarray, list[int]]:
    size, requirements = map.strip().split(":")
    x, y = size.split("x")
    x = int(x)
    y = int(y)
    requirements = [int(r) for r in requirements.strip().split(" ")]
    return np.zeros((x, y), dtype=bool), requirements


def rotations(shape: np.ndarray):
    for k in range(4):
        yield np.rot90(shape, k)


def overlaps(map: np.ndarray, shape: np.ndarray, pos: tuple[int, int]) -> bool:
    r, c = pos
    h, w = shape.shape
    if r + h > map.shape[0] or c + w > map.shape[1]:
        return True
    region = map[r : r + h, c : c + w]
    return np.any(region & shape)


def add(map: np.ndarray, shape: np.ndarray, pos: tuple[int, int]):
    r, c = pos
    h, w = shape.shape
    map[r : r + h, c : c + w] |= shape


def remove(map: np.ndarray, shape: np.ndarray, pos: tuple[int, int]):
    r, c = pos
    h, w = shape.shape
    map[r : r + h, c : c + w] &= np.invert(shape)


def solve(map: np.ndarray, shapes_to_place: list[np.ndarray]) -> bool:
    if not shapes_to_place:
        return True

    shape = shapes_to_place[0]
    remaining_shapes = shapes_to_place[1:]

    for rot in rotations(shape):
        for r in range(map.shape[0]):
            for c in range(map.shape[1]):
                if not overlaps(map, rot, (r, c)):
                    add(map, rot, (r, c))
                    if solve(map, remaining_shapes):
                        return True
                    remove(map, rot, (r, c))
    return False


def part_one(maps: list[tuple[np.ndarray, list[int]]], shapes: list[np.ndarray]) -> int:
    cells_per_shape = [int(shape.sum()) for shape in shapes]

    def can_fit_simple_by_size(map: np.ndarray, requirements: list[int]) -> bool:
        needed = sum(r * c for r, c in zip(requirements, cells_per_shape))
        return needed <= map.size

    solutions = 0
    for map, requirements in maps:
        if not can_fit_simple_by_size(map, requirements):
            continue

        shapes_to_place = []
        for shape_idx, count in enumerate(requirements):
            shapes_to_place.extend([shapes[shape_idx]] * count)

        if solve(map.copy(), shapes_to_place):
            solutions += 1

    return solutions


def main():
    with open("12/input.txt", "r") as file:
        sections = file.read().split("\n\n")

    maps = [parse_map(map) for map in sections[-1].strip().split("\n")]
    shapes = [parse_shape(shape) for shape in sections[:-1]]

    print("solution for part 1:", part_one(maps, shapes))
    print("solution for part 2: Merry Christmas!")


if __name__ == "__main__":
    main()
