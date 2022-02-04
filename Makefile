install:
	sudo rm -rf ~/.gcsv
	go build .
	sudo mv ./gcsv /usr/bin/gcsv
	mkdir ~/.gcsv
	cp secret.json ~/.gcsv/secret.json