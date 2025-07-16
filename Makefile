GOROOT=$(shell go env GOROOT)

default: TiltMan web

TiltMan:
	go build -o $@ .

webserver: web
	    python3 -m http.server --directory ./web

web:
	rm -rf web/
	mkdir -p web/
	env GOOS=js GOARCH=wasm go build -o web/game.wasm .
	cp $(GOROOT)/lib/wasm/wasm_exec.js web/
	magick icons/icon.png -resize 192x192 web/icon-192x192.png
	magick icons/icon.png -resize 32x32 web/icon-32x32.png
	magick icons/icon.png -resize 512x512 web/icon-512x512.png
	optipng -quiet web/icon-192x192.png web/icon-32x32.png web/icon-512x512.png
	cp index.html manifest.json sw.js web/

publish: web
	rsync -a web/ kaka:/var/www/fortyfootgames.duckdns.org/TiltMan/

clean:
	rm -rf web
	rm -f TiltMan

.PHONY: web clean webserver default publish
