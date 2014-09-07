package main

func getPercentValue(value, percent int) int {
    return value * percent / 100
}

func abs(x int) int {
    if x > 0 {
        return x
    }
    return -x
}

func max(x, y int) int {
    if x > y {
        return x
    }
    return y
}

func min(x, y int) int {
    if x < y {
        return x
    }
    return y
}
