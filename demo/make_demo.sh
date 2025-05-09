#!/bin/bash
# Run demo_script and record
asciinema rec kontext.cast -c "./demo_script.sh" --overwrite

# Convert to GIF
asciicast2gif kontext.cast kontext.gif

# Optimize GIF
gifsicle -O3 --colors 256 kontext.gif -o kontext-demo.gif

echo "✅ Done! Your demo GIF is ready as kontext-demo.gif"
