package thanks

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/beevik/guid"
	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
	"github.com/noltio/diky2020-backend/pkg/server"
)

// ImageUploadHandler handles requests for thanks image upload
func ImageUploadHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if file == nil {
		c.Locals("uploadedFile", false)
		return c.Next()
	}
	if err != nil {
		fmt.Println(err.Error())
		return server.APIInternalServerError(c)
	}

	// filesize condition
	if file.Size > 8388608 {
		return server.APIError(c, "Maximální velikost nahrávaného obrázku je 8 MB.", 400)
	}

	// mime type condition
	ct := file.Header["Content-Type"][0]
	mimes := map[string]string{
		"image/jpeg": "jpg",
		"image/png":  "png",
		"image/gif":  "gif",
	}
	if _, ok := mimes[ct]; !ok {
		return server.APIError(c, "Povolené formáty obrázků jsou pouze .jpg, .jpeg, .png a .gif.", 400)
	}

	// filename (for both temp and resized file)
	filename := fmt.Sprintf("%s.%s", guid.New().String(), mimes[ct])
	filepath := "/upload/" + filename

	// save the original file
	c.SaveFile(file, filepath)

	// open the uploaded file
	nrf, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err.Error())
		return server.APIInternalServerError(c)
	}
	defer nrf.Close()

	// decode jpeg into image.Image
	var img image.Image
	switch ct {
	case "image/jpeg":
		img, err = jpeg.Decode(nrf)
	case "image/png":
		img, err = png.Decode(nrf)
	case "image/gif":
		img, err = gif.Decode(nrf)
	}
	if err != nil {
		fmt.Println(err.Error())
		return server.APIInternalServerError(c)
	}

	// resize to width 600 using Lanczos resampling
	// and preserve the original aspect ratio
	m := resize.Resize(600, 0, img, resize.Lanczos3)

	// save new file
	out, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err.Error())
		return server.APIInternalServerError(c)
	}
	defer out.Close()

	// write new image to file
	switch ct {
	case "image/jpeg":
		err = jpeg.Encode(out, m, nil)
	case "image/png":
		err = png.Encode(out, m)
	case "image/gif":
		err = gif.Encode(out, m, nil)
	}
	if err != nil {
		fmt.Println(err.Error())
		return server.APIInternalServerError(c)
	}

	// store the data in local variables and move to the next handler
	c.Locals("uploadedFile", true)
	c.Locals("uploadedFileName", filename)
	return c.Next()
}
