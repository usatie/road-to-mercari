[![Go Reference](https://pkg.go.dev/badge/github.com/usatie/road-to-mercari/ex00/convert.svg)](https://pkg.go.dev/github.com/usatie/road-to-mercari/ex00/convert)

# Road to mercari
Convert all JPG files in a directory (including the subdirectories) to PNG format.

```
$> go mod tidy
$> go build
$> ./convert
error: invalid argument
$> ./convert nosuchdirectory
error: nosuchdirectory: no such file or directory
$> ls -1 images
42tokyo_logo.jpg
profile_photo.jpg
$> ./convert images
$> ls -1 images
42tokyo_logo.jpg
42tokyo_logo.png
profile_photo.jpg
profile_photo.png
$> echo 'aaa' > images/test.txt
$> ls -1 images
42tokyo_logo.jpg
42tokyo_logo.png
profile_photo.jpg
profile_photo.png
test.txt
$> ./convert images
error: images/test.txt is not a valid file
```
