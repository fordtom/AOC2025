use std::fs;

struct RedTile {
    x: u64,
    y: u64,
}

enum Edge {
    Horizontal { x1: u64, x2: u64, y: u64 },
    Vertical { x: u64, y1: u64, y2: u64 },
}

impl RedTile {
    pub fn area_with(&self, other: &RedTile) -> u64 {
        (self.x.abs_diff(other.x) + 1) * (self.y.abs_diff(other.y) + 1)
    }

    pub fn get_other_vertices(&self, other: &RedTile) -> [RedTile; 2] {
        [
            RedTile {
                x: other.x,
                y: self.y,
            },
            RedTile {
                x: self.x,
                y: other.y,
            },
        ]
    }

    pub fn get_edge(&self, other: &RedTile) -> Edge {
        if self.x == other.x {
            Edge::Vertical {
                x: self.x,
                y1: self.y.min(other.y),
                y2: self.y.max(other.y),
            }
        } else {
            Edge::Horizontal {
                x1: self.x.min(other.x),
                x2: self.x.max(other.x),
                y: self.y,
            }
        }
    }

    pub fn get_rectangle_edges(&self, other: &RedTile) -> Vec<Edge> {
        let [a, b] = self.get_other_vertices(other);
        vec![
            self.get_edge(&a),
            a.get_edge(other),
            other.get_edge(&b),
            b.get_edge(self),
        ]
    }

    pub fn ray_crosses_edge(&self, edge: &Edge) -> bool {
        let ray = Edge::Horizontal {
            x1: self.x,
            x2: std::u64::MAX,
            y: self.y,
        };
        ray.does_intersect_ray_logic(edge)
    }
}

impl Edge {
    pub fn does_intersect_ray_logic(&self, other: &Edge) -> bool {
        match (self, other) {
            (Edge::Vertical { x, y1, y2 }, Edge::Horizontal { x1, x2, y }) => {
                *x1 <= *x && *x <= *x2 && *y1 <= *y && *y < *y2
            }
            (Edge::Horizontal { x1, x2, y }, Edge::Vertical { x, y1, y2 }) => {
                *x1 < *x && *x <= *x2 && *y1 <= *y && *y < *y2
            }
            _ => false,
        }
    }

    pub fn does_intersect_edge_logic(&self, other: &Edge) -> bool {
        match (self, other) {
            (Edge::Vertical { x, y1, y2 }, Edge::Horizontal { x1, x2, y }) => {
                x1 < x && x < x2 && y1 < y && y < y2
            }
            (Edge::Horizontal { x1, x2, y }, Edge::Vertical { x, y1, y2 }) => {
                x1 < x && x < x2 && y1 < y && y < y2
            }
            _ => false,
        }
    }
}

fn part_one(tiles: &Vec<RedTile>) -> u64 {
    let mut largest_area = 0;
    for tile in tiles {
        for other in tiles {
            let area = tile.area_with(other);
            if area > largest_area {
                largest_area = area;
            }
        }
    }

    largest_area
}

fn part_two(tiles: &Vec<RedTile>) -> u64 {
    let mut edges = vec![];
    for i in 0..tiles.len() {
        let second = if i + 1 < tiles.len() { i + 1 } else { 0 };
        edges.push(tiles[i].get_edge(&tiles[second]));
    }

    let mut largest_area = 0;
    for tile in tiles {
        'next_tile: for other in tiles {
            let new_vertices = tile.get_other_vertices(other);

            for v in &new_vertices {
                let crossings = edges.iter().filter(|edge| v.ray_crosses_edge(edge)).count();
                if crossings % 2 == 0 {
                    continue 'next_tile;
                }
            }

            let rect_edges = tile.get_rectangle_edges(other);
            for e in rect_edges {
                if edges.iter().any(|edge| edge.does_intersect_edge_logic(&e)) {
                    continue 'next_tile;
                }
            }

            let area = tile.area_with(other);
            if area > largest_area {
                largest_area = area;
            }
        }
    }

    largest_area
}

fn main() {
    let input = fs::read_to_string("input.txt").unwrap().trim().to_string();

    let tiles = input
        .split("\n")
        .map(|line| {
            let coords = line
                .split(",")
                .map(|num| num.parse::<u64>().unwrap())
                .collect::<Vec<u64>>();
            RedTile {
                x: coords[0],
                y: coords[1],
            }
        })
        .collect::<Vec<RedTile>>();

    let largest_area = part_one(&tiles);
    println!("Solution for part 1: {}", largest_area);

    let largest_area = part_two(&tiles);
    println!("Solution for part 2: {}", largest_area);
}
