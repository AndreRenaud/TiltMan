- What's a good name for a video game that combines the old marble labyrinth game (played on a wooden board where you have to tilt the board to roll the marble around and avoid falling into a hole), with pacman, where you have to 'tilt' the game to make pacman move
- Create a minimal ebitengine main.go file that opens a 1280x720 window. Title the window "TiltMan"
- Create a file marble.go, which draws a circle. It should have a type Marble which remembers its position, speed, direction etc... and each update it should update its position based on its velocity.
- Add a file 'map.go', which defines the map. The map is described using an ASCII drawing, such as:
  Where difference characters define what will be drawn in that square. The following characters are currently defined:
  - `#` a solid wall
  - `.` an empty bit of floor
  - `<` a bit of floor that slows the marble down
  - `>` a bit of floor that speeds the marble up
- Change the Marbles update function so that instead of actually changing the x/y coord, it calculates the new ones and returns them. We can then verify that we're not inside a wall before applying them.
- Add '(' & ')' character support to the map. These should slow down/speed up, but not as much as the < & > ones
- Add a new module spritesheet.go, which loads an image containing a bunch of tile images, and is capable of drawing any one of them on the screen arbitrarily as requested. Don't use it anywhere yet.
- Have the map.Draw function take a callback which should return a 32x32 ebiten.Image for the given x/y coord. The callback will be supplied with the m.Tiles map, and the current coordinate.
- Add a map from tile.Type to ebiten.Image to the game structure so we can cache the images
- Use the spritesheet module to load the assets/grass.png and assets/stone.png images. Then use grasses 0,0 image as the empty square, and stones 1,1 image as the wall square
- Need to make the wall drawing a bit smarter. If the wall has non-wall to the right of it, it should be using the sprite from 2,1 instead.
- Use the embed module to embed the assets/ directory, and update spritesheet to work with fs.FS
- Add some javascript to index.html so that deviceorientation events are captured & logged to the javascript console
- Add a channel of orientation events, which main.go will pull from (if it has some), and which events_wasm.go exposes a function to push into. Call this function from index.html when new orientation events happen
- Update index.html to make this a proper web app when saved to the home page on mobile devices
- (On ChatGPT) - Create a 192x192 icon for a game that is a marble run game where tilt the phone to play
- Write mapgenerator.go which creates a maze of arbitrary size (specified during construction). The maze should be constructed out of ascii characters. Make sure the border is all '#' symbols.
  - The output should be a string array, with one line per horizontal row of the maze.