
all:app
app: dist/index.html 
	go build -o dist/middag
dist/index.html:
	node_modules/.bin/parcel build frontend/src/index.html
clean:
	rm -f dist/*
