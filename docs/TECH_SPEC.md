## Ridership Map

### Data Inputs

Ridership Data:  ETL to import ridership spreadsheet into DB (can automate later)
Route Segments: 
ETL to import routes into the DB from the bus tracker api 
Match the Bus Tracker routes to the ridership routes

### Visualization

Navigatable map (like google maps) of chicago. Something where I can define lines on the map using coordinates from the bus tracker api route segments

Filters: 

- month and year (to select what ridership records to show)
- metric: what data point to display via the thickness of the lines (avg weekday ridership, avg weekend ridership...etc.)
- display route info on hover