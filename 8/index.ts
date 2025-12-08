import * as fs from 'fs';

type Box = {
    x: number;
    y: number;
    z: number;
    parent: Box;
    circuitSize: number;
};

const input = fs.readFileSync('input.txt', 'utf-8').trim().split('\n').map(line => {
    const [x, y, z] = line.split(',').map(Number);
    const box: Box = { x, y, z, circuitSize: 1 } as Box;
    box.parent = box;
    return box;
});

function resetBoxes(boxes: Box[]): void {
    boxes.forEach(b => { b.parent = b; b.circuitSize = 1 });
}

function distance(box1: Box, box2: Box): number {
    const dx = Math.abs(box1.x - box2.x);
    const dy = Math.abs(box1.y - box2.y);
    const dz = Math.abs(box1.z - box2.z);
    return Math.sqrt(dx * dx + dy * dy + dz * dz);
}

function findRoot(box: Box): Box {
    if (box.parent === box) return box; 
    return findRoot(box.parent);
}

function join(box1: Box, box2: Box): void {
    const root1 = findRoot(box1);
    const root2 = findRoot(box2);
    if (root1 === root2) return;

    if (root1.circuitSize < root2.circuitSize) {
        root1.parent = root2;
        root2.circuitSize += root1.circuitSize;
    } else {
        root2.parent = root1;
        root1.circuitSize += root2.circuitSize;
    }
}

function buildEdges(boxes: Box[]) {
    const edges: { i: number; j: number; dist: number }[] = [];
    for (let i = 0; i < boxes.length; i++) {
        for (let j = i + 1; j < boxes.length; j++) {
            edges.push({ i, j, dist: distance(boxes[i]!, boxes[j]!) });
        }
    }
    return edges.sort((a, b) => a.dist - b.dist);
}

const uniqueCircuits = (boxes: Box[]) => new Set(boxes.map(findRoot)).size;

function partOne(boxes: Box[]): number {
    const edges = buildEdges(boxes);

    for (const { i, j } of edges.slice(0, 1000)) {
        join(boxes[i]!, boxes[j]!);
    }

    const networks = new Map<Box, number>();
    for (const box of boxes) {
        const root = findRoot(box);
        networks.set(root, root.circuitSize);
    }

    let ordered = Array.from(networks.entries()).sort((a, b) => b[1] - a[1]);
    return ordered[0]![1]! * ordered[1]![1]! * ordered[2]![1]!;
}

function partTwo(boxes: Box[]): number {
    const edges = buildEdges(boxes);
    let idx = boxes.length - 1;

    for (const { i, j } of edges.slice(0, idx)) {
        join(boxes[i]!, boxes[j]!);
    }

    let unique = uniqueCircuits(boxes);

    while (unique > 1) {
        for (let i = 0; i < unique - 1; i++) {
            for (const { i, j } of edges.slice(idx, idx + unique - 1)) {
                join(boxes[i]!, boxes[j]!);
            }
        }
        unique = uniqueCircuits(boxes);
        idx += unique - 1;
    }

    let final_edge = edges[idx];
    return boxes[final_edge!.i!]!.x! * boxes[final_edge!.j!]!.x!;
}

console.log("Part 1:", partOne(input));
resetBoxes(input);
console.log("Part 2:", partTwo(input));
