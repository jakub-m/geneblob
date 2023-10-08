images:
	go run doc/generate_images.go
run: clean
	go run main.go
clean:
	rm -rvf tmp_*.png doc/_.gif doc/_*.png

