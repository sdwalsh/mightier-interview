  ______________________
((                      ))
 )) Mightier Interview ((
((                      ))
  ----------------------

Discussion Topics:
  ● Any interesting trade-offs or edge cases you ran into.
  ● Any steps necessary to run your program on a command line, including
  compilation steps (if required) and the runtime versions you used (e.g., python
  version 3.7.2) - please assume we are not familiar with the language
  development environment when writing your steps.

////////////
Discussion:
////////////

It's been a while since I worked with Go, so I figured this was a good opportunity to practice.
Plus, the runtime is included in the binary so once it's been compiled for a platform, it just works.

Since this is a game, it's a state machine-- transitioning between blue and red states. Scoring in this case
is applied during the transitions between states. If the plan was to expand this, a state machine library
(https://github.com/looplab/fsm) could be used to define both the states and the transitions making the
creation of more complicated scoring algorithms easier. I didn't look into it much given the timebox,
but it's something that could be of interest.

Only one external library is used. I brought in (github.com/gocarina/gocsv) in order to use struct tags
to easily marshal the csv data into my structs

There was on small error in the example from the document: HIBACHO_HERO (HIBACHI_HERO)

I left everything in main.go since it's only 165 lines. Left a number of comments throughout the code in order
to give a better idea of my thought process while writing.

Adding tests would be a good next step. Go's built-in testing library is fantastic and is fairly straight
forward to set up.

/////////////////
Compile and run:
/////////////////
I've provided binaries for both Linux x64_86, and MacOS x64_86 (mightier-interview and mightier-interview-mac)

Built with: go1.17.3 linux/amd64 on Ubuntu 21.10

Note: I don't believe I'm using any features that are specific to this version,
      and it'll likely compile on older versions

`go build followed by running the binary ./mightier-interview`
or
`go run main.go`

OPTIONAL FLAGS:
-properties=properties.csv
-session=session.json

e.g. `go run main.go -properties=123.csv -session=123.json`

default value for csv: properties.csv
default value for json: session.json
