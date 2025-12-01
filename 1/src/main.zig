const std = @import("std");

pub fn partOne(contents: []u8) !i32 {
    var steps = std.mem.splitSequence(u8, contents, "\n");
    var loc: i32 = 50;
    var incr: i32 = 0;

    while (steps.next()) |step| {
        if (step.len == 0) {
            continue;
        }
        const direction = step[0];
        const distance = try std.fmt.parseInt(i32, step[1..], 10);
        switch (direction) {
            'L' => {
                loc -= distance;
            },
            'R' => {
                loc += distance;
            },
            else => unreachable,
        }

        loc = @rem(loc, 100);

        if (loc == 0) {
            incr += 1;
        }
    }

    return incr;
}

pub fn partTwo(contents: []u8) !i32 {
    var steps = std.mem.splitSequence(u8, contents, "\n");
    var loc: i32 = 50;
    var incr: i32 = 0;

    while (steps.next()) |step| {
        if (step.len == 0) {
            continue;
        }

        const direction = step[0];
        const distance = try std.fmt.parseInt(i32, step[1..], 10);
        const rotations = @divTrunc(distance, 100);
        incr += rotations;
        const rem = @rem(distance, 100);

        switch (direction) {
            'L' => {
                if (loc > 0 and rem > 0 and loc <= rem) {
                    incr += 1;
                }

                loc -= rem;

                if (loc < 0) {
                    loc += 100;
                }
            },
            'R' => {
                loc += rem;

                if (loc > 99) {
                    loc -= 100;
                    incr += 1;
                }
            },
            else => unreachable,
        }
    }

    return incr;
}

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const file = try std.fs.cwd().openFile("input.txt", .{});
    defer file.close();

    const file_size = (try file.stat()).size;
    const contents = try file.readToEndAlloc(allocator, file_size);
    defer allocator.free(contents);

    const partOneSolution = try partOne(contents);
    std.debug.print("Solution for part 1: {}\n", .{partOneSolution});

    const partTwoSolution = try partTwo(contents);
    std.debug.print("Solution for part 2: {}\n", .{partTwoSolution});
}
