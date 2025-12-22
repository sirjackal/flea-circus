# Project Euler: Problem 213 "Flea Circus"
A 30×30 grid of squares contains 900 fleas, initially one flea per square. When a bell is rung, each flea jumps to an adjacent square at random (usually 4 possibilities, except for fleas on the edge of the grid or at the corners).

What is the expected number of unoccupied squares after 50 rings of the bell? Give your answer rounded to six decimal places.

## Implementation plan in Golang:
* implement a single simulation first
* run multiple simulations and calculate the average
* run simulations in parallel (e.g. worker pool using channels and other sync primitives as required)
