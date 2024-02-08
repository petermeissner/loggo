import * as d3 from "https://cdn.jsdelivr.net/npm/d3@7/+esm";


// variables / selectors for general UI elements

d33 = d3;

// content
var content = d3.select("#content");

content.append("h1").text("Content");


// routes
d3.select('#routes')
  .selectAll('p')
  .data(['[GET] /', '[GET] /log_files', '[GET] /ws/echo', '[GET] /ws/test_stream'])
  .text(function(d){ return d;})
  .join('p');
