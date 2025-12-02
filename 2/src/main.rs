use std::fs;

fn part_one(ranges: &Vec<(u64, u64)>) -> u64 {
    let mut sum = 0;

    for range in ranges {
        for i in range.0..=(range.1 + 1) {
            let s = i.to_string();

            if s.len() % 2 == 0 && s[..s.len() / 2] == s[s.len() / 2..] {
                sum += i;
            }
        }
    }

    sum
}

fn part_two(ranges: &Vec<(u64, u64)>) -> u64 {
    let mut sum = 0;

    for range in ranges {
        for i in range.0..=(range.1 + 1) {
            let s = i.to_string();

            if (1..=s.len() / 2).any(|chunk_size| {
                s.len() % chunk_size == 0
                    && s.as_bytes()
                        .chunks(chunk_size)
                        .all(|chunk| chunk == &s.as_bytes()[..chunk_size])
            }) {
                sum += i;
            }
        }
    }

    sum
}

fn main() {
    let input = fs::read_to_string("input.txt").unwrap().trim().to_string();

    let ranges = input
        .split(",")
        .map(|range| {
            let mut parts = range.split("-").map(|num| num.parse::<u64>().unwrap());
            (parts.next().unwrap(), parts.next().unwrap())
        })
        .collect::<Vec<(u64, u64)>>();

    let sum = part_one(&ranges);
    println!("Solution for part 1: {}", sum);

    let sum = part_two(&ranges);
    println!("Solution for part 2: {}", sum);
}
