const std = @import("std");

fn partOne(allocator: std.mem.Allocator, contents: []u8) !usize {
    const width = std.mem.indexOf(u8, contents, "\n") orelse contents.len;
    var rows = std.mem.splitSequence(u8, contents, "\n");

    var active_cols = try allocator.alloc(bool, width);
    defer allocator.free(active_cols);
    @memset(active_cols, false);

    var splits: usize = 0;

    while (rows.next()) |row| {
        for (row, 0..) |col, i| {
            if (col == 'S') {
                active_cols[i] = true;
            } else if (col == '^' and active_cols[i]) {
                active_cols[i] = false;

                splits += 1;
                if (i > 0) {
                    active_cols[i - 1] = true;
                }
                if (i < width - 1) {
                    active_cols[i + 1] = true;
                }
            }
        }
    }

    return splits;
}

fn partTwo(allocator: std.mem.Allocator, contents: []u8) !usize {
    const width = std.mem.indexOf(u8, contents, "\n") orelse contents.len;
    var rows = std.mem.splitSequence(u8, contents, "\n");

    var timelines = try allocator.alloc(usize, width);
    defer allocator.free(timelines);
    @memset(timelines, 0);

    var timeline_sum: usize = 0;

    while (rows.next()) |row| {
        for (row, 0..) |col, i| {
            if (col == 'S') {
                timelines[i] = 1;
            } else if (col == '^' and timelines[i] > 0) {
                if (i > 0) {
                    timelines[i - 1] += timelines[i];
                }
                if (i < width - 1) {
                    timelines[i + 1] += timelines[i];
                }
                timelines[i] = 0;
            }
        }
    }

    for (timelines) |timeline| {
        timeline_sum += timeline;
    }

    return timeline_sum;
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

    const splits = try partOne(allocator, contents);
    std.debug.print("Solution for part 1: {}\n", .{splits});

    const timeline_sum = try partTwo(allocator, contents);
    std.debug.print("Solution for part 2: {}\n", .{timeline_sum});
}
