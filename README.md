# Multithreaded City Traffic Simulator (single-node)

Implement a multithreaded city traffic simulator where every car and semaphore will be independent entities. Street directions and semaphores must be respected by the cars.

![City Traffic](city-traffic.png)

## Technical Requirements

- The city's map can be static or automatically generated.
- Cars and semaphore number can be configured on game's start.
- For every car, define a random start and destination point.
- Define a random speed for each car.
- If a car detect another car on his route and it's slower, it must slow down its speed.
- Each car and semaphore behaviour will be implemented as a separated thread.
- Cars and Semaphores threads must use the same map or city layout data structure resource.
- Display finished cars' routes.
- Display each car's speed.

## General Requirements

- Make sure that you complete the below defined deliverables.
- Source code dependencies must be clearly documented.

## Deliverables

- Source code can be in a single student's account and the other team members can contribute to the same repository.
- Architecture Document - [ARCHITECTURE.md](ARCHITECTURE.md)
  - Details on how you designed and implemented your solution
  - Project architecture description, diagrams, charts and everything related to the way you think/design/build your program
- Build/Run automation (`Makefile` and documentation - [CITY-TRAFFIC.md](CITY-TRAFFIC.md)
  - Instructions on how to build and run your program
- Project's presentation (5-10 minutes)
  - A video presentation that will be delivered to the professor, in youtube preferable

## Permitted programming languages

- Multithreaded core backend
  - C
  - Go
- User Interface (optional)
  - Any
  - If it's terminal, output must be human-readable

## Grading Policy

| Concept                      | Grade |
| ---------------------------- | ----- |
| Architecture Document        | 20%   |
| Multithreaded implementation | 30%   |
| Build Automation             | 20%   |
| Coding best practices        | 10%   |
| Presentation                 | 20%   |
| TOTAL                        | 100%  |

- **Free Lab Bonus**

You can get an extra bonus if you implement an Artificial Intelleigence algorithm in semaphores synchronization. This must be documented in the `ARCHITECTURE.md` file.
