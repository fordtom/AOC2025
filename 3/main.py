def part_one(banks) -> int:
    joltages = []

    for bank in banks:
        s = bank[:-1]
        i = s.index(max(s))
        s = bank[i + 1 :]

        joltages.append(10 * int(bank[i]) + int(max(s)))

    return sum(joltages)


def part_two(banks) -> int:
    joltages = []

    for bank in banks:
        j = 0
        s = bank
        idx = -1

        for i in range(11, -1, -1):
            s = s[idx + 1 :]
            idx = s.index(max(s[:-i] if i != 0 else s))
            j += int(s[idx]) * 10**i

        joltages.append(j)

    return sum(joltages)


def main():
    with open("3/input.txt", "r") as file:
        banks = file.read().splitlines()

    print(f"Solution for part 1: {part_one(banks)}")
    print(f"Solution for part 2: {part_two(banks)}")


if __name__ == "__main__":
    main()
