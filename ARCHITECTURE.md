# Architecture

## Libraries

- [Ebiten library](https://ebiten.org/)
  - open source game library for the Go programming language. Ebitengine's simple API allows you to quickly and easily develop 2D games that can be deployed across multiple platforms

![Ebiten Architecture](https://camo.githubusercontent.com/674c4bb22fc289ea4584cb355090516c915ad544ade00c1bf86414e655b49147/68747470733a2f2f65626974656e2e6f72672f696d616765732f6f76657276696577322e312e706e67)

## Resources

[Ebiten Examples](https://ebiten.org/examples/flappy.html)
[Mastering Concurrency in Go](https://www.educative.io/courses/mastering-concurrency-in-go)

## ToDo's

- [x] The city's map can be static or automatically generated.
- [x] Cars and semaphore number can be configured on game's start.
- [x] For every car, define a random start and destination point.
- [x] Define a random speed for each car.
- [x] If a car detect another car on his route and it's slower, it must slow down its speed.
- [x] Each car and semaphore behaviour will be implemented as a separated thread.
- [x] Cars and Semaphores threads must use the same map or city layout data structure resource.
- [x] Display finished cars' routes.
- [x] Display each car's speed.
