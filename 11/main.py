from functools import cache


# We can treat it as a DAG if false
def has_cycle(devices: dict[str, list[str]]) -> bool:
    UNSEEN, WIP, DONE = 0, 1, 2
    status = {name: UNSEEN for name in devices}

    def check_cycle(device: str) -> bool:
        if device not in devices or device == "out":
            return False
        if status[device] == WIP:
            return True  # back edge = cycle
        if status[device] == DONE:
            return False
        status[device] = WIP
        for out in devices[device]:
            if check_cycle(out):
                return True
        status[device] = DONE
        return False

    return any(check_cycle(d) for d in devices)


def part_one(devices: dict[str, list[str]]) -> int:
    def walk(outputs: list[str]) -> int:
        total = 0
        for output in outputs:
            if output == "out":
                total += 1
            else:
                total += walk(devices[output])
        return total

    return walk(devices["you"])


def part_two(devices: dict[str, list[str]]) -> int:
    @cache
    def walk(device: str, end: str) -> int:
        if device == end:
            return 1
        if device == "out" or device not in devices:
            return 0
        return sum(walk(next, end) for next in devices[device])

    route1 = walk("svr", "dac") * walk("dac", "fft") * walk("fft", "out")
    route2 = walk("svr", "fft") * walk("fft", "dac") * walk("dac", "out")

    return route1 + route2


def main():
    with open("11/input.txt", "r") as file:
        lines = file.read().splitlines()
        devices = {
            name: [output.strip() for output in outputs.split(" ") if output.strip()]
            for name, outputs in [line.split(":") for line in lines if line.strip()]
        }

        if not has_cycle(devices):
            print("solution for part 1:", part_one(devices))
            print("solution for part 2:", part_two(devices))


if __name__ == "__main__":
    main()
