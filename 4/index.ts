import * as fs from 'fs';

const input = fs.readFileSync('input.txt', 'utf-8').trim().split('\n');

const surrounding: [number, number][] = [[-1, -1], [-1, 0], [-1, 1], [0, -1], [0, 1], [1, -1], [1, 0], [1, 1]];

function isRoll(grid: string[], row: number, col: number): boolean {
    return grid[row]?.[col] === '@';
}

function shouldRemove(grid: string[], row: number, col: number): boolean {
    if (!isRoll(grid, row, col)) return false;
    const neighbours = surrounding.filter(([dr, dc]) => isRoll(grid, row + dr, col + dc)).length;
    return neighbours < 4;
}

function partOne(grid: string[]): number {
    let count = 0;
    for (let row = 0; row < grid.length; row++) {
        for (let col = 0; col < grid[row]!.length; col++) {
            if (shouldRemove(grid, row, col)) count++;
        }
    }
    return count;
}

function removalPass(grid: string[]): [string[], number] {
    let newGrid = grid.map(row => row.split(''));
    let count = 0;

    for (let row = 0; row < grid.length; row++) {
        for (let col = 0; col < grid[row]!.length; col++) {
            if (shouldRemove(grid, row, col)) {
                newGrid[row]![col] = '.';
                count++;
            }
        }
    }

    return [newGrid.map(row => row.join('')), count];
}

function partTwo(grid: string[]): number {
    let count = 0;
    let removed: number;

    do {
        [grid, removed] = removalPass(grid);
        count += removed;
    } while (removed > 0);

    return count;
}

console.log(partOne(input));
console.log(partTwo(input));
