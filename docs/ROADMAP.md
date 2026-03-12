The goal is to A ridership overlay is a way to combine two different datasets on the same map:
1️⃣ Transit performance (speed, delay, reliability)
2️⃣ Ridership demand (how many people actually use each route)

## Planned Enhancements
1. Ridership Map: ridership overlay of cta bus routes
2. Speed Ridership Overlay
- by segment (identify slow corridors)
- Thickness of map → number of riders
- Segment Color → average speed
3. Passenger delay → ridership x delay 
- e.g. "where would improvements help the most people?"
- riders, speed, impact
- 25k, slow, huge
= 10k, slow, moderate

## Possible Future Enhancements
Scheduled vs actual frequency
corridor congestion analysis
bus lane impact simulation
ridership vs speed visualizations

## Ridership Map

- [ ] Render 1 bus route
  - [ ] Import Ridership Data
  - [ ] Get route info from the cta bus tracker api
  - [ ] Merge into hardcoded geojson served on frontend
- [ ] Render all bus routes
