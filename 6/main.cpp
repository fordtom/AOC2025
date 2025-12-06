#include <array>
#include <fstream>
#include <iostream>
#include <string>
#include <vector>

constexpr int LINES = 5;
constexpr int DATA_ROWS = LINES - 1;

void apply_operation(long long& accumulator, char op, long long value) {
    if (op == '+') {
        accumulator += value;
    } else {
        accumulator *= value;
    }
}

long long get_initial_value(char op) {
    return op == '+' ? 0 : 1;
}

std::vector<size_t> find_operator_positions(const std::string& line) {
    std::vector<size_t> positions;
    for (size_t i = 0; i < line.size(); i++) {
        if (line[i] == '+' || line[i] == '*') {
            positions.push_back(i);
        }
    }
    return positions;
}

int main() {
    std::ifstream file("input.txt");
    std::array<std::string, LINES> lines;

    for (int i = 0; i < LINES; i++) {
        std::getline(file, lines[i]);
    }

    auto columns = find_operator_positions(lines[DATA_ROWS]);

    long long part1 = 0;
    long long part2 = 0;
    for (size_t i = 0; i < columns.size(); i++) {
        size_t start = columns[i];
        size_t len = (i + 1 < columns.size()) ? columns[i + 1] - columns[i] : lines[0].size() - start;

        long long out = get_initial_value(lines[DATA_ROWS][columns[i]]);
        for (int j = 0; j < (LINES - 1); j++) {
            long long val = std::stoll(lines[j].substr(start, len));

            apply_operation(out, lines[DATA_ROWS][columns[i]], val);
        }
        part1 += out;

        out = get_initial_value(lines[DATA_ROWS][columns[i]]);
        for (size_t j = 0; j < len; j++) {
            std::string digit_string;
            for (int row = 0; row < LINES - 1; row++) {
                char c = lines[row][start + j];
                if (c != ' ') {
                    digit_string += c;
                }
            }
            if (digit_string.empty()) {
                continue;
            }
            long long val = std::stoll(digit_string);

            apply_operation(out, lines[DATA_ROWS][columns[i]], val);
        }
        part2 += out;
    }

    std::cout << "Solution for part 1: " << part1 << '\n';
    std::cout << "Solution for part 2: " << part2 << '\n';

    return 0;
}
